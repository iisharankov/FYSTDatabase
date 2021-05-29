package main

import (
	"context"
	"log"
	"strconv"

	"github.com/iisharankov/FYSTDatabase/datasets"
	"github.com/minio/minio-go/v7"
)

func uploadData(reply datasets.ClientUploadReply, file datasets.File) error {
	minioUseSSL, _ := strconv.ParseBool(minioUseSSL) // Convert env var to bool

	minioEndpoint = "0.0.0.0:9001"
	var S3Instance = ObjectMetadata{
		ctx:      context.Background(),
		endpoint: minioEndpoint,
		id:       minioAccessKeyID,
		password: minioSecretAccessKey,
		useSSL:   minioUseSSL,
	}

	S3Instance.initMinio()
	_, err := copyFile(S3Instance, file, reply)
	if err != nil {
		log.Println("err in copyFile", err)
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

func copyFile(minioInstance ObjectMetadata, file datasets.File, serverReply datasets.ClientUploadReply) (int64, error) {
	reply := serverReply.UploadLocation
	bucketList, _ := minioInstance.ListBuckets()

	if !stringInSlice(reply, bucketList) {
		minioInstance.makeBucket(reply, "us-east-1")
	}

	// Upload the zip file
	log.Println(file.Name, file.URL)

	minioInstance.UploadObject(reply, file.Name, file.URL, "application/")
	return 0, nil
}
