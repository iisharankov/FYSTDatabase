package main

import (
	"os"
)

var dbName = os.Getenv("DATABASE_NAME")
var dbAddress = os.Getenv("MYSQL_IP")
var dbUsername = os.Getenv("MYSQL_USER")
var dbPassword = os.Getenv("MYSQL_PASSWORD")

var minioEndpoint = os.Getenv("MINIO_ENDPOINT")
var minioAccessKeyID = os.Getenv("MINIO_ACCESS_ID")
var minioSecretAccessKey = os.Getenv("MINIO_SECRET_KEY")
var minioUseSSL = os.Getenv("MINIO_SSL")

const (
	sqlTimeLayout string = "2006-01-2 15:04:05"
)

// GlobalPTStackArray is a struct containing an array of structs
var dbCon DatabaseConnection

func main() {
	// go TransferClock()
	startAPIServer()

}
