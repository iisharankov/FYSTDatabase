package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

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

func (data *TransferData) Clock() {

	/*
		Clock() is the heart of the ACU Simulator. It keeps a steady pace of 20Hz using a ticker. Clock() uses two channels to communicate externally
		with other parts of the ACU Simulator.

		The first of which is clockSplineChan, a channel tied with the SplineGenorator() function. This channel feeds the processed spline model
		created by spineInterpolation to execute. The contents of this channel is an array of structs, specifically ProcessedTPT structs,
		which are simply embedded TimePositionTransfer structs that also contain a UnixTimestamp field for better time management. Once this
		channel returns new spline data, it is added to splineData, an array built specifically for storing the future positions of the
		Simulator. From here, upon every return of the ticker, if splineData is not empty, the next struct is removed from it and processed.
	*/

	// splineData is array of structs containing data for when to move to next Az/El nextSplinePoint is a single struct of that array
	var splineData *int
	var rcvdData int
	var i int

	twoHzClock := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case rcvdData = <-data.S3TransferChan:
				splineData = &rcvdData // Delete Spline by replacing with new spline data
				log.Println("Received New Spline!", &rcvdData, splineData)

			case <-twoHzClock.C:
				data.moveOffTelscope()
				i++
				log.Println(i)

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

	//--------------------- Finds all the rules in the database
	for _, val := range listOfFilesThatNeedToBeBackedUp {

		// Find the Name of the File with the given FileID. Only exists on entry with RuleID=1 (fyst bucket)
		query := fmt.Sprintf(`select l.URL from Log l join Rule r on r.RuleId=l.RuleID 
		join BackupLocation BL on BL.LocationID=r.LocationID 
		where l.FileID=%v and BL.S3Bucket="%v";`, val.FileID, "fyst")
		queryReturn, err := dbCon.QueryRead(query, &struct{ Name string }{})
		if err != nil {
			log.Println(err)
		}
		objName := queryReturn.Interface().([]struct{ Name string })[0].Name // type assert from reflect.Value to struct

		if err := tData.copyFile(val, objName); err != nil {
			log.Println("Error in Copying file \n", err)
		} else if err = addRowToLog(val.FileID, val.RuleID, objName); err != nil {
			log.Println("Error adding row to log\n", err)
		} else {
			log.Printf("Adding file=%v with rule=%v Log table\n", val.FileID, val.RuleID)
		}
	}
}

func (tData *TransferData) copyFile(src FilesThatNeedToBeBackedUp, objName string) error {
	log.Printf(" - - - - - File is %v, Rule is %v, bucket is %v\n", src.FileID, src.RuleID, src.BucketName)

	// Is ths needed? Buckets should exist.
	location := "us-east-1"
	tempFile := "temp.zip"
	bucket := strings.ToLower(strings.Replace(src.BucketName, " ", "", -1))

	// See if bucket exists in set, if not, attempt to create
	if !tData.dstS3.Buckets[bucket] {
		log.Printf("Bucket '%v' doesn't exist, creating. \n", bucket)
		tData.dstS3.makeBucket(bucket, location)
		tData.dstS3.Buckets[bucket] = true
	}

	// var ttt string
	// log.Println("~~~~ OBJECTS ARE")

	// objectCh := tData.srcS3.minioClient.ListObjects(tData.srcS3.ctx, "fyst", minio.ListObjectsOptions{
	// 	Recursive: true,
	// })
	// for object := range objectCh {
	// 	if object.Err != nil {
	// 		fmt.Println(object.Err)
	// 	}
	// 	fmt.Println(object.Key)
	// 	ttt = object.Key
	// }

	// ttt := strings.Split(src.URL[strings.LastIndex(src.URL, "/")+1:], ".")[0]
	log.Println("------Downloading file - ", objName)
	err := tData.srcS3.minioClient.FGetObject(tData.srcS3.ctx, "fyst", objName, tempFile, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	// Upload the zip file
	last := src.URL[strings.LastIndex(src.URL, "/")+1:]
	log.Println("------Uploading file called", last)
	n, err := tData.dstS3.minioClient.FPutObject(tData.dstS3.ctx, bucket, objName, tempFile, minio.PutObjectOptions{ContentType: "application/zip"})
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Successfully uploaded %s of size %d to bucket %s \n", src.BucketName, n.Size, src.URL)

	if err := os.Remove(tempFile); err != nil {
		log.Println(err)
	}

	return nil
}
