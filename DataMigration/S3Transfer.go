package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

// TODO: This is in FYSTDatabase/datasets, but my go compiler will not
// TOO: recognize the struct even after I push the updated file to git.
// TODO: It should see it but can't? Is this just a local caching issue?
// FilesThatNeedToBeBackedUp lists all the data required to move file from FYST to external location
type FilesThatNeedToBeBackedUp struct {
	FileID         int
	RuleID         int
	InstrumentID   int
	Size           int
	InstrumentName string
	DateCreated    string
	URL            string
	ByteHash       string
	BucketName     string
}

/* uploadQueue is triggered at a set interval to check if there are new files
on the database that need to be copied externally. The query is quite involved
using all tables, and may not be best to be queried often. Another option is a
trigger through a channel from the server API when a POST log request is sent.
A third is a trigger somehow from the database when a file is addec to the Log
table. Not sure how to implement that but I know those triggers exist, might
be the cleanest. For now as a proof-of-concept it's checking on a loop */
func (data *TransferData) uploadQueue() {
	var possibleTrigger *int
	var rcvdTrigger int

	// TODO: Add this sampling to the main if we plan to keep a timing loop
	twoHzClock := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			// Trigger through channel is an example, is never called from
			// anywhere. Do we want this implementation?
			case rcvdTrigger = <-data.S3TransferChan:
				possibleTrigger = &rcvdTrigger // Delete Spline by replacing with new spline data
				log.Println("Received New Spline!", &rcvdTrigger, possibleTrigger)

			case <-twoHzClock.C:
				data.moveOffTelscope()
			}
		}
	}()
}

func (tData *TransferData) moveOffTelscope() {
	/* What files have not been copied to locations where active rules exist (pending uploads per
	all locations)? The following query scans which rules apply to every file in the database,
	and finds which files have not been copied to a location with an active rule. The result is
	every location a file should be copied to, but hasn't. */
	var SQLQuery string = `select f.FileID, r.RuleID, f.Size, i.InstrumentID,
	i.InstrumentName, f.DateCreated, f.ObjectStorage, f.HashOfBytes, BL.S3Bucket
	from ObjectFile f
	join Instrument i on i.InstrumentID=f.InstrumentID
	join Rule r on r.InstrumentID=i.InstrumentID
	join BackupLocation BL on BL.LocationID=r.LocationID
	left join Log l on l.FileID=f.FileID and r.RuleId=l.RuleID
	where l.FileID is null AND r.Active=1 order by f.DateCreated;`
	var structOfFilesThatNeedToBeBackedUp FilesThatNeedToBeBackedUp
	outputRows, err := dbCon.QueryRead(SQLQuery, &structOfFilesThatNeedToBeBackedUp)
	if err != nil {
		log.Println(err)
	}
	listOfFilesThatNeedToBeBackedUp, _ := outputRows.Interface().([]FilesThatNeedToBeBackedUp)

	// Try to upload each file returned by the above query to it's respective location.
	// This section could be optimized to spin up mulitple workers depending on locations/files, etc.
	for _, val := range listOfFilesThatNeedToBeBackedUp {

		/*
			Currently, each file gets uploaded to the fyst bucket on the local object storage. This means
			there exists a entry in log that will contain the Name of the object within the fyst bucket.
			This is the unique object storage name/ID, and is located in Logs since it may be different
			for each upload. The above query can't access this name, so we need to query for it separately
			so we can upload each 'val' object with the same name. If we choose a different naming scheme
			then this may not be needed */
		query := fmt.Sprintf(`select l.URL from Log l join Rule r on r.RuleId=l.RuleID 
		join BackupLocation BL on BL.LocationID=r.LocationID 
		where l.FileID=%v and BL.S3Bucket="%v";`, val.FileID, "fyst") // TODO: Hardcoded "fyst"
		queryReturn, err := dbCon.QueryRead(query, &struct{ Name string }{})
		if err != nil {
			log.Println(err)
		}

		objName := queryReturn.Interface().([]struct{ Name string })[0].Name // type assert from reflect.Value to struct
		if err := tData.copyFileToExternal(val, objName); err != nil {
			log.Println("Error in Copying file \n", err)
			return
		}

		// If the object was uploaded sucessfully, add a value in the log
		if err = addRowToLog(val.FileID, val.RuleID, objName); err != nil {
			// TODO: If object uploaded but addRowToLog fails, what should happen?
			// Should upload object be deleted, addRowToLog triggered again?
			// Need to maintain data integrity!
			log.Println("Error adding row to log\n", err)
		} else {
			log.Printf("Adding file=%v with rule=%v Log table\n", val.FileID, val.RuleID)
		}
	}
}

func (tData *TransferData) copyFileToExternal(src FilesThatNeedToBeBackedUp, objName string) error {
	log.Printf("Transfering file with FileID=%v, RuleID=%v to bucket named %v in %v\n",
		src.FileID, src.RuleID, src.BucketName, objName)

	// Is ths needed? Buckets should exist.
	location := "us-east-1" // TODO: This should be called from the table
	tempFile := "temp.zip"  // temp file for temp solution ;)

	// Bucket name must be lowercase and with no whitespace
	bucket := strings.ToLower(strings.Replace(src.BucketName, " ", "", -1))

	// See if bucket exists in set, if not, attempt to create
	if !tData.dstS3.Buckets[bucket] {
		log.Printf("Bucket '%v' doesn't exist, creating. \n", bucket)
		tData.dstS3.makeBucket(bucket, location)
		tData.dstS3.Buckets[bucket] = true // Store in set that is in TransferData.dstS3
	}

	// TODO: Temporary soln until streaming object download
	log.Println("Downloading file - ", objName)
	err := tData.srcS3.minioClient.FGetObject(tData.srcS3.ctx, "fyst", objName, tempFile, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	// Upload the zip file
	last := src.URL[strings.LastIndex(src.URL, "/")+1:]
	log.Println("Uploading file called", last)
	n, err := tData.dstS3.minioClient.FPutObject(tData.dstS3.ctx, bucket, objName, tempFile, minio.PutObjectOptions{ContentType: "application/zip"})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d to bucket %s \n", src.BucketName, n.Size, src.URL)
	if err := os.Remove(tempFile); err != nil {
		log.Println(err) // TODO: Temporary soln, need to delete local download
	}
	return nil
}
