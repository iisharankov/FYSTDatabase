package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/iisharankov/FYSTDatabase/datasets"
	"github.com/minio/minio-go/v7"
)

/* uploadQueue is triggered at a set interval to check if there are new files
on the database that need to be copied externally. The query is quite involved
using all tables, and may not be best to be queried often. Another option is a
trigger through a channel from the server API when a POST log request is sent.
A third is a trigger somehow from the database when a file is addec to the Record
table. Not sure how to implement that but I know those triggers exist, might
be the cleanest. For now as a proof-of-concept it's checking on a loop */
func (tData *TransferData) uploadQueue() {
	// var possibleTrigger *int
	// var rcvdTrigger int

	// TODO: Add this sampling to the main if we plan to keep a timing loop
	twoHzClock := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			// Trigger through channel is an example.
			// case rcvdTrigger = <-tData.S3TransferChan:
			// 	possibleTrigger = &rcvdTrigger //
			// 	log.Println("Received something!", &rcvdTrigger, possibleTrigger)

			case <-twoHzClock.C:
				tData.moveOffTelscope()
			}
		}
	}()
}

func (tData *TransferData) moveOffTelscope() {
	tData.Lock()
	defer tData.Unlock()

	// v // What files have not been copied to locations where active rules exist
	// (i.e. pending uploads per all locations)? The following query scans which rules
	// apply to every file in the database, and finds which files have not been
	// copied to a location with an active rule. The result is every location a
	// file should be copied to, but hasn't.
	var SQLQuery string = `select f.FileID, c.URL, f.Size, f.DateCreated, f.MD5sum, 
	r.RuleID, BL.S3Bucket, BL.LocationName
	from Files f join Copies c on c.FileID=f.FileID
	join Instruments i on i.InstrumentID=f.InstrumentID
	join Rules r on r.InstrumentID=i.InstrumentID
	join Locations BL on BL.LocationID=r.LocationID
	left join Records l on l.FileID=f.FileID and r.RuleId=l.RuleID
	where l.FileID is null AND r.Active=1 order by f.DateCreated;`
	var listOfFilesThatNeedToBeBackedUp []datasets.FilesThatNeedToBeBackedUp
	if outputRows, err := dbCon.QueryRead(SQLQuery, &datasets.FilesThatNeedToBeBackedUp{}); err != nil {
		log.Println(err)
	} else {
		listOfFilesThatNeedToBeBackedUp, _ = outputRows.Interface().([]datasets.FilesThatNeedToBeBackedUp)
	}
	log.Println("Files to upload: ", listOfFilesThatNeedToBeBackedUp)

	// v // Try to upload each file returned by the above query to it's respective location.
	// This section could be optimized to spin up multiple workers depending on locations/files, etc.
	for _, val := range listOfFilesThatNeedToBeBackedUp {
		// objURL may be an empty table, in which case nothing should happen
		// if objURL := queryReturn.Interface().([]struct{ Name string }); len(objURL) != 0 {
		if err := tData.copyFileToExternal(val); err != nil {
			log.Println("Error in Copying file \n", err)
			return
		}

		// If the object was uploaded sucessfully, add a value in the log
		if err := addRowToRecord(val.FileID, val.RuleID); err != nil {
			log.Println("Error adding row to log\n", err)
		} else {
			log.Printf("Added file=%v with rule=%v into Records table\n", val.FileID, val.RuleID)
		}
	}
}

// connectToObjectStorage takes a string and returns
func connectToObjectStorage(location string) (ObjectMetadata, error) {

	// Below block queries Locations table for Location=location row (should be unique entry)
	var locationsRow datasets.LocationsTable
	query := fmt.Sprintf(`select * from Locations where LocationName="%v";`, location)
	if queryReturn, err := dbCon.QueryRead(query, &datasets.LocationsTable{}); err != nil {
		return ObjectMetadata{}, err
	} else if temp := queryReturn.Interface().([]datasets.LocationsTable); len(temp) == 1 {
		// Above line makes sure the return query is of length one
		locationsRow = temp[0]
	} else {
		errMsg := fmt.Sprintf(`Query requesting Location from Locations table 
		returned %v rows when it should have returned only one`, len(temp))
		return ObjectMetadata{}, fmt.Errorf(errMsg)
	}

	// Use the data from the SQL query to populate the ObjectMetadata struct
	transferData := ObjectMetadata{
		ctx:      context.Background(),
		address:  locationsRow.Address,
		accessID: locationsRow.AccessID,
		secretID: locationsRow.SecretID,
		useSSL:   locationsRow.SSL}

	// Initialize the Minio connection and return it
	transferData.initMinio()

	return transferData, nil
}

func (tData *TransferData) copyFileToExternal(src datasets.FilesThatNeedToBeBackedUp) error {

	localObjSrg, err := connectToObjectStorage("Observatory")
	remoteObjSrg, err := connectToObjectStorage(src.UploadLocation)

	// v // See if bucket exists (must be lowercase with no whitespace), else create
	bucket := strings.ToLower(strings.Replace(src.BucketName, " ", "", -1))
	if ok, err := remoteObjSrg.minioClient.BucketExists(remoteObjSrg.ctx, bucket); err != nil {
		return err
	} else if ok == false {
		log.Printf("Bucket '%v' doesn't exist, creating. \n", bucket)
		remoteObjSrg.makeBucket(bucket, "us-east-1") // TODO: Location should be called from DB
	}
	// ^ //

	// TODO: Temporary soln as streaming fails at high throughput
	tempFile := src.FileName
	err = localObjSrg.minioClient.FGetObject(localObjSrg.ctx, "fyst", src.FileName, tempFile, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	// Upload the zip file
	_, err = remoteObjSrg.minioClient.FPutObject(remoteObjSrg.ctx, bucket, src.FileName, tempFile, minio.PutObjectOptions{ContentType: "application/zip"})
	if err != nil {
		return err
	}

	if err := os.Remove(tempFile); err != nil {
		log.Println(err) // TODO: Temporary soln, need to delete local download
	}

	// Streaming method - Broken - Fails at high throughput
	// v // Stream object from local object storage to off-site
	// object, err := localObjSrg.minioClient.GetObject(localObjSrg.ctx, "fyst", src.FileName, minio.GetObjectOptions{})
	// if err != nil {
	// 	return err
	// }

	// _, err = remoteObjSrg.minioClient.PutObject(remoteObjSrg.ctx, bucket, src.FileName, object, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	// if err != nil {
	// 	return err
	// }
	// ^ //

	return nil
}
