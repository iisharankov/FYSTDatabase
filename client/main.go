package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/iisharankov/FYSTDatabase/datasets"
)

var minioEndpoint = os.Getenv("MINIO_ENDPOINT")
var minioAccessKeyID = os.Getenv("MINIO_ACCESS_ID")
var minioSecretAccessKey = os.Getenv("MINIO_SECRET_KEY")
var minioUseSSL = os.Getenv("MINIO_SSL")

func getenv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	apiAddr := flag.String("api", getenv("FYST_DB_ADDR", "localhost:8700"), "DB API address")
	connection := NewDBAPI(*apiAddr)

	filescmd := flag.NewFlagSet("files", flag.ExitOnError)
	rulescmd := flag.NewFlagSet("rules", flag.ExitOnError)
	uploadcmd := flag.NewFlagSet("upload", flag.ExitOnError)

	switch os.Args[1] {
	case "files":
		filename := filescmd.String("filename", "", "Files to query")
		filescmd.Parse(os.Args[2:])

		reply, err := connection.getFiles(*filename)
		if err != nil {
			fmt.Println(err)
		}

		var ans datasets.FilesTable
		if err = json.Unmarshal(reply, &ans); err != nil {
			log.Println(err)
		} else {
			log.Println(ans)
		}

	case "rules":
		id := rulescmd.String("id", "", "Files to query")
		rulescmd.Parse(os.Args[2:])

		ans, err := connection.getRules(*id)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(ans))
		}

	case "upload":
		name := uploadcmd.String("name", "", "")
		instrument := uploadcmd.Int("instrument", 0, "")
		md5 := uploadcmd.String("md5", "", "")
		date := uploadcmd.String("date", "", "")
		size := uploadcmd.Int("size", 0, "")
		url := uploadcmd.String("url", "", "")
		uploadcmd.Parse(os.Args[2:])

		formatedDate, err := time.Parse("2006-01-02T15:04:05Z", *date)
		if err != nil {
			log.Println(err)
		}

		file := datasets.File{
			Name:        *name,
			Instrument:  *instrument,
			MD5Sum:      *md5,
			DateCreated: formatedDate,
			Size:        *size,
			URL:         *url,
		}

		// Tell server new file exists and unmarshal the reply for use
		var ans datasets.ClientUploadReply
		reply, err := connection.requestToUploadFile(file)
		if err != nil {
			log.Println(file.Name, err)
			return
		} else if err = json.Unmarshal(reply, &ans); err != nil {
			panic(err)
		}

		log.Printf("Reply was %v\n", ans)

		// a, err := connection.logGET()
		// log.Println("a was", string(a), "err was", err)
		// upload file to bucket URL in JSON reply and then ask server to update log
		err = uploadData(ans, file)
		if err != nil {
			log.Println(err)
			return
			// TODO: what to do if upload fails?
		}

		a, err := connection.requestToUpdateLog(file.Name, ans)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Log response:", string(a))

		b, err := connection.requestToUpdateCopies(ans.FileName, ans.LocationID)
		if err != nil {
			log.Println("-", err)
			return
		}
		log.Println("Copies response:", string(b))

	default:
		fmt.Println("unexpected command")
		os.Exit(1)
	}
}
