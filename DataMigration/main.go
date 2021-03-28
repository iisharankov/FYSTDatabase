package main

import (
	"context"
	"os"
	"strconv"

	"github.com/iisharankov/FYSTDatabase/datasets"
	"github.com/minio/minio-go/v7"
)

var dbName = os.Getenv("DATABASE_NAME")
var dbAddress = os.Getenv("MYSQL_IP")
var dbUsername = os.Getenv("MYSQL_USER")
var dbPassword = os.Getenv("MYSQL_PASSWORD")

var minioEndpoint = os.Getenv("MINIO_ENDPOINT")
var minioAccessKeyID = os.Getenv("MINIO_ACCESS_ID")
var minioSecretKey = os.Getenv("MINIO_SECRET_KEY")
var minioUseSSL = os.Getenv("MINIO_SSL")

var s3Endpoint = os.Getenv("S3_ENDPOINT")
var s3AccessKeyID = os.Getenv("S3_ACCESS_ID")
var s3SecretKey = os.Getenv("S3_SECRET_KEY")
var s3UseSSL = os.Getenv("S3_SSL")

const (
	sqlTimeLayout string = "2006-01-2 15:04:05"
)

// GlobalPTStackArray is a struct containing an array of structs
var dbCon DatabaseConnection

// SimulatorMetadata is a class to easily pass around internal channels
type ServerMetadata struct {
}

// ObjectMetadata is a struct
type ObjectMetadata struct {
	ctx         context.Context
	minioClient *minio.Client
	endpoint    string
	id          string
	password    string
	useSSL      bool
	Buckets     map[string]bool `json:"buckets"`
}

type TransferData struct {
	S3TransferChan chan int
	srcS3          *ObjectMetadata
	dstS3          *ObjectMetadata
	currentFile    datasets.File
}

func main() {

	s3UseSSL, _ := strconv.ParseBool(s3UseSSL)       // Convert env var to bool
	minioUseSSL, _ := strconv.ParseBool(minioUseSSL) // Convert env var to bool

	transferData := TransferData{
		S3TransferChan: make(chan int),
		srcS3: &ObjectMetadata{ // Destination S3 instance
			ctx:      context.Background(),
			endpoint: minioEndpoint,
			id:       minioAccessKeyID,
			password: minioSecretKey,
			useSSL:   minioUseSSL,
			Buckets:  make(map[string]bool)},
		dstS3: &ObjectMetadata{ // Source S3 instance (local)
			ctx:      context.Background(),
			endpoint: s3Endpoint,
			id:       s3AccessKeyID,
			password: s3SecretKey,
			useSSL:   s3UseSSL,
			Buckets:  make(map[string]bool)},
	}

	transferData.srcS3.initMinio()
	transferData.dstS3.initMinio()

	// p := makeSimulator()
	go transferData.Clock()

	startAPIServer()
}
