package main

import (
	"log"
	"os"
	"sync"

	"github.com/minio/minio-go/v7"
)

// Recieved from docker-compose or env vars if running locally
var dbName = os.Getenv("DATABASE_NAME")
var dbAddress = os.Getenv("MYSQL_IP")
var dbUsername = os.Getenv("MYSQL_USER")
var dbPassword = os.Getenv("MYSQL_PASSWORD")

const (
	sqlTimeLayout string = "2006-01-2 15:04:05"
)

// DatabaseConnection holds the connection to the database
// so methods in DBOverhead can share the connection
var dbCon DatabaseConnection

/* ObjectStorageConnection stores all the information for a given minio instance,
currently it also stores the buckets created during the instance, but this
may be possible to offload to the database. */
type ObjectStorageConnection struct { // TODO Rename
	minioClient *minio.Client
	address     string
	accessID    string
	secretID    string
	useSSL      bool
}

// This may need an overhaul since this is not generalized well.
type TransferData struct {
	S3TransferChan chan int
	sync.Mutex
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	transferData := TransferData{S3TransferChan: make(chan int)}
	go transferData.uploadQueue()

	// Start API that has all the endpoints ready to listen for traffic
	startAPIServer()
}
