package myminio

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
			log.Println("EWD")
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

func (minioInstance *ObjectMetadata) uploadObject(bucketName, objectName, filePath, contentType string) {
	// Upload the zip file with FPutObject
	n, err := minioInstance.minioClient.FPutObject(minioInstance.ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d \n", objectName, n.Size)
}

func (minioInstance *ObjectMetadata) listObjects(bucketName, prefix string) {
	// What to do with this, how do context.Contexts scale with other ctx's?
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

// func main() {

// 	endpoint := "0.0.0.0:9000"
// 	accessKeyID := "iisharankov"
// 	secretAccessKey := "iisharankov"
// 	useSSL := false
// 	var S3Instance = ObjectMetadata{ctx: context.Background(), endpoint: endpoint, id: accessKeyID, password: secretAccessKey, useSSL: useSSL}

// 	S3Instance.initMinio()

// 	// Make a new bucket called mymusic.
// 	bucketName := "testminiobucket"
// 	location := "us-east-1"
// 	S3Instance.makeBucket(bucketName, location)

// 	// Upload the zip file
// 	objectName := "TheBeatles_LetItBe.zip"
// 	filePath := "/home/iisharankov/Downloads/" + objectName
// 	contentType := "application/zip"
// 	S3Instance.uploadObject(bucketName, objectName, filePath, contentType)

// 	S3Instance.listObjects(bucketName, "myprefix")

// 	S3Instance.removeBucket(bucketName)

// 	S3Instance.removeObject(bucketName, objectName)

// 	S3Instance.listObjects(bucketName, "myprefix")

// 	S3Instance.removeBucket(bucketName)
// }

// More reading: https://dev.to/maelvls/why-is-go111module-everywhere-and-everything-about-go-modules-24k
