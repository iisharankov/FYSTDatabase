package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func initMinio(endpoint, ID, password string, useSSL bool) (*minio.Client, error) {
	// Initialize minio client object.
	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ID, password, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return s3Client, nil
}

func makeBucket(ctx context.Context, s3Client *minio.Client, bucketName, location string) {

	opts := minio.MakeBucketOptions{
		Region: "us-east-1",
	}

	err := s3Client.MakeBucket(ctx, bucketName, opts)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s3Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Println("EWD")
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func uploadObject(ctx context.Context, s3Client *minio.Client, bucketName, objectName, filePath, contentType string) {
	// Upload the zip file with FPutObject
	n, err := s3Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d \n", objectName, n)
}

func listObjects(ctx context.Context, s3Client *minio.Client, bucketName, prefix string) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := s3Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		// Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
			return
		}
		log.Println()
		log.Println(object)
	}
}
func main() {

	ctx := context.Background()
	// endpoint := "play.min.io"
	endpoint := "0.0.0.0:9000"
	accessKeyID := "iisharankov"
	secretAccessKey := "iisharankov"
	useSSL := false

	minioClient, err := initMinio(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Println(err)
	}

	// Make a new bucket called mymusic.
	bucketName := "mymusic"
	location := "us-east-1"
	makeBucket(ctx, minioClient, bucketName, location)

	// Upload the zip file
	// objectName := "TheBeatles_LetItBe.zip"
	// filePath := "/home/iisharankov/Downloads/" + objectName
	// contentType := "application/zip"
	// uploadObject(ctx, minioClient, bucketName, objectName, filePath, contentType)

	listObjects(ctx, minioClient, bucketName, "myprefix")
}

// Fix was:
// working directory is not part of a module
// go mod init  go/src/github.com/minio/minio-go

// More reading: https://dev.to/maelvls/why-is-go111module-everywhere-and-everything-about-go-modules-24k
