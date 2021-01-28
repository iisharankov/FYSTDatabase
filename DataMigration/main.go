package main

import OverheadSQL "github.com/iisharankov/FYSTDatabase/OverheadSQL"

const (
	dbUsername    string = "iisharankov"
	dbPassword    string = "iisharankov"
	dbAddress     string = ""
	dbName        string = "mydb"
	sqlTimeLayout string = "2006-01-2 15:04:05"

	minioEndpoint        string = "0.0.0.0:9000"
	minioAccessKeyID     string = "iisharankov"
	minioSecretAccessKey string = "iisharankov"
	minioUseSSL          bool   = false
)

// GlobalPTStackArray is a struct containing an array of structs
var dbCon OverheadSQL.DatabaseConnection

func main() {
	startAPIServer()

}
