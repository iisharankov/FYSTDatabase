package main

import (
	"context"

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
	Storage        string
	ByteHash       string
	LocationName   string
}

func uploadData(reply ServerUploadReply, file File) error {
	var S3Instance = ObjectMetadata{
		ctx:      context.Background(),
		endpoint: minioEndpoint,
		id:       minioAccessKeyID,
		password: minioSecretAccessKey,
		useSSL:   minioUseSSL}

	S3Instance.initMinio()

	_, err := copyFile(S3Instance, file, reply)
	if err != nil { // Inverse of normal!
		return err
	}
	return nil
}

func stringInSlice(a string, list []minio.BucketInfo) bool {
	for _, b := range list {
		if b.Name == a {
			return true
		}
	}
	return false
}

func copyFile(minioInstance ObjectMetadata, file File, serverReply ServerUploadReply) (int64, error) {
	reply := serverReply.UploadLocation
	bucketList, _ := minioInstance.ListBuckets()
	if !stringInSlice(reply, bucketList) {
		minioInstance.makeBucket(reply, "us-east-1")
	}

	// Upload the zip file
	minioInstance.UploadObject(reply, file.Name, file.URL, "application/zip")
	return 0, nil
}
