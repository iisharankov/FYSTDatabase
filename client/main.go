package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type File struct {
	Name        string    `json:"name"`
	Instrument  int       `json:"instrument"`
	MD5Sum      string    `json:"md5sum"`
	DateCreated time.Time `json:"date_created"`
	Size        int       `json:"size"`
	URL         string    `json:"url"`
}

type ServerUploadReply struct {
	FileID         int    `json:"file_id"`
	UploadLocation string `json:"upload_location"`
}

func getenv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func main() {
	apiAddr := flag.String("api", getenv("FYST_DB_ADDR", "localhost:8700"), "DB API address")
	connection := NewDBAPI(*apiAddr)

	filescmd := flag.NewFlagSet("files", flag.ExitOnError)
	rulescmd := flag.NewFlagSet("rules", flag.ExitOnError)
	uploadcmd := flag.NewFlagSet("upload", flag.ExitOnError)

	switch os.Args[1] {
	case "files":
		id := filescmd.String("id", "", "Files to query")
		filescmd.Parse(os.Args[2:])

		ans, err := connection.getFiles(*id)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(ans))
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

		layout := "2006-01-02T15:04:05Z" // Golang Time is a mess...
		formatedDate, err := time.Parse(layout, *date)
		if err != nil {
			fmt.Println(err)
		}

		file := File{Name: *name, Instrument: *instrument, MD5Sum: *md5,
			DateCreated: formatedDate, Size: *size, URL: *url}

		// Tell server new file exists and unmarshal the reply for use
		var ans ServerUploadReply
		if reply, err := connection.requestToUploadFile(file); err != nil {
			log.Println(err)
		} else {
			if err = json.Unmarshal(reply, &ans); err != nil {
				panic(err)
			}
		}
		fmt.Printf("Reply was %v\n", ans)
		// upload file to bucket URL in JSON reply and then ask server to update log
		if err := uploadData(ans, file); err != nil {
			log.Println(err)
		} else {
			a, err := connection.requestToUpdateLog(ans)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("update Log response was:", string(a))
			}
		}

	default:
		fmt.Println("unexpected command")
		os.Exit(1)
	}
}
