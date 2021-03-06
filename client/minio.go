package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ObjectMetadata is a struct
type ObjectMetadata struct {
	ctx         context.Context
	minioClient *minio.Client
	endpoint    string
	id          string
	password    string
	useSSL      bool
}

func (minioInstance *ObjectMetadata) initMinio() {
	// Initialize minio client object.
	client, err := minio.New(minioInstance.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioInstance.id, minioInstance.password, ""),
		Secure: minioInstance.useSSL,
	})

	if err != nil {
		log.Println(err)
	} else {
		minioInstance.minioClient = client
	}
}

func (minioInstance *ObjectMetadata) makeBucket(bucketName, location string) {

	opts := minio.MakeBucketOptions{
		Region: location,
	}

	err := minioInstance.minioClient.MakeBucket(minioInstance.ctx, bucketName, opts)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioInstance.minioClient.BucketExists(minioInstance.ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket named %s already exists!\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket called %s\n", bucketName)
	}
}

func (minioInstance *ObjectMetadata) removeBucket(bucketName string) {

	err := minioInstance.minioClient.RemoveBucket(minioInstance.ctx, bucketName)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Bucket %s removed successfully \n", bucketName)
	}
}

// ListBuckets returns all buckets
func (minioInstance *ObjectMetadata) ListBuckets() ([]minio.BucketInfo, error) {
	// What to do with this, how do context.Contexts scale with other contexts?
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	buckets, err := minioInstance.minioClient.ListBuckets(minioInstance.ctx)
	return buckets, err
}

// UploadObject uploads a given file from a given filepath to a cloud bucket
func (minioInstance *ObjectMetadata) UploadObject(bucketName, objectName, filePath, contentType string) {
	// Upload the zip file with FPutObject
	// TODO: Setting filePath as objectName as well to help uniqueness
	n, err := minioInstance.minioClient.FPutObject(minioInstance.ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d to bucket %s \n", objectName, n.Size, bucketName)
}

func (minioInstance *ObjectMetadata) listObjects(bucketName, prefix string) {
	// What to do with this, how do context.Contexts scale with other contexts?
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	objectCh := minioInstance.minioClient.ListObjects(minioInstance.ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
			return
		}
		log.Printf("Object with Key='%v' and size %d bytes found\n", object.Key, object.Size)
	}
}

func (minioInstance *ObjectMetadata) removeObject(bucketName, objectName string) {
	// Upload the zip file with FPutObject
	err := minioInstance.minioClient.RemoveObject(minioInstance.ctx, bucketName, objectName, minio.RemoveObjectOptions{GovernanceBypass: false})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully removed object %s \n", objectName)
}
