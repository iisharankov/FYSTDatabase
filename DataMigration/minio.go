package main

import (
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/* I initially created this file to hold all the methods I'd use for minio handeling, but
I found two issues pop up. 1) most of these methods are wrappers that don't really save
much (if any) code, like requestObject. It's literally wrapping another function for no
benifit. Also, I found I needed this library for both the server and client, so the
question arose if it should be it's own package. I don't think that's necessary, and this
file can probably be removed and integrated where the minio methods are used directly */

func (minioInstance *ObjectMetadata) initMinio() {
	// Initialize minio client object.
	// TODO: Return error for this!
	if client, err := minio.New(minioInstance.address, &minio.Options{
		Creds:  credentials.NewStaticV4(minioInstance.accessID, minioInstance.secretID, ""),
		Secure: minioInstance.useSSL,
	}); err != nil {
		log.Println(err)
	} else {
		minioInstance.minioClient = client
	}
}

func (minioInstance *ObjectMetadata) makeBucket(bucketName, location string) {
	opts := minio.MakeBucketOptions{Region: location}

	err := minioInstance.minioClient.MakeBucket(minioInstance.ctx, bucketName, opts)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioInstance.minioClient.BucketExists(minioInstance.ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket named %s already exists!\n", bucketName)
		} else {
			log.Println("Bucket create error")
			log.Println(err)
		}
	} else {
		log.Printf("Successfully created bucket called %s\n", bucketName)
	}
}

func (minioInstance *ObjectMetadata) requestObject(bucketName, objectName string) {
	// Upload the zip file with FPutObject
	err := minioInstance.minioClient.FGetObject(minioInstance.ctx,
		bucketName, objectName, "/tmp/myobject", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

// err := minioInstance.minioClient.(minioInstance.ctx, bucketName, objectName, minio.RemoveObjectOptions{GovernanceBypass: false})
// if err != nil {
// 	log.Fatalln(err)
// }

// func (minioInstance *ObjectMetadata) removeBucket(bucketName string) {

// 	err := minioInstance.minioClient.RemoveBucket(minioInstance.ctx, bucketName)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		log.Printf("Bucket %s removed successfully \n", bucketName)
// 	}
// }

// ListBuckets returns all buckets
// func (minioInstance *ObjectMetadata) ListBuckets() ([]minio.BucketInfo, error) {
// 	// What to do with this, how do context.Contexts scale with other ctx's?
// 	// ctx, cancel := context.WithCancel(context.Background())
// 	// defer cancel()

// 	buckets, err := minioInstance.minioClient.ListBuckets(minioInstance.ctx)
// 	return buckets, err
// }

// UploadObject uploads a given file from a given filepath to a cloud bucket
// func (minioInstance *ObjectMetadata) UploadObject(bucketName, objectName, filePath, contentType string) {
// 	// Upload the zip file with FPutObject
// 	n, err := minioInstance.minioClient.FPutObject(minioInstance.ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	log.Printf("Successfully uploaded %s of size %d to bucket %s \n", objectName, n.Size, bucketName)
// }

// func (minioInstance *ObjectMetadata) listObjects(bucketName, prefix string) {
// 	// What to do with this, how do context.Contexts scale with other ctx's?
// 	// ctx, cancel := context.WithCancel(context.Background())
// 	// defer cancel()

// 	objectCh := minioInstance.minioClient.ListObjects(minioInstance.ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

// 	for object := range objectCh {
// 		if object.Err != nil {
// 			log.Println(object.Err)
// 			return
// 		}
// 		log.Printf("Object with Key='%v' and size %d bytes found\n", object.Key, object.Size)
// 	}
// }

// func (minioInstance *ObjectMetadata) removeObject(bucketName, objectName string) {
// 	// Upload the zip file with FPutObject
// 	err := minioInstance.minioClient.RemoveObject(minioInstance.ctx, bucketName, objectName, minio.RemoveObjectOptions{GovernanceBypass: false})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	log.Printf("Successfully removed object %s \n", objectName)
// }
