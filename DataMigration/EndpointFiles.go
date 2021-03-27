package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iisharankov/FYSTDatabase/datasets"
)

func AddFileToDatabaseEndpoint(w http.ResponseWriter, r *http.Request) {
	var newFile datasets.File
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&newFile)
	log.Println("new file:", newFile)

	// Find largest index for FileID column
	query := "SELECT FileID FROM ObjectFile ORDER BY FileID DESC LIMIT 1"
	queryReturn, err := dbCon.QueryRead(query, &struct{ FileID int }{})
	if err != nil {
		log.Println(err)
	}

	var newFileID int
	if IDStruct, _ := queryReturn.Interface().([]struct{ FileID int }); len(IDStruct) == 0 {
		newFileID = 0 // Necessary due to empty table case (returns no row)
	} else {
		newFileID = IDStruct[0].FileID
	}

	// Try to add given file to ObjectFile table with next FileID
	err = addRowToObjectFile(newFile, w)
	if err != nil {
		var msg string = "addRowToObjectFile method failed, file not added to database"
		jsonResponse(w, concatErrors(err, msg), http.StatusBadRequest)
		return
	}

	// If successful, respond to client with upload location details
	var rule datasets.BackupLocationTable
	query = `select * from BackupLocation b where b.LocationID = 1`
	queryBackupLocation, err := dbCon.QueryRead(query, &rule)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	returnQuery, _ := queryBackupLocation.Interface().([]datasets.BackupLocationTable)
	if returnQuery[0].S3Bucket == "" {
		jsonResponse(w, err, http.StatusConflict)
		return
	}

	response := datasets.ClientUploadReply{
		S3Metadata: datasets.S3Metadata{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKeyID,
			SecretAccessKey: minioSecretKey,
			UseSSL:          minioUseSSL,
		},
		FileID:         newFileID,
		UploadLocation: returnQuery[0].S3Bucket,
	}

	// Create JSON with metadata necessary for the client
	replyData, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Client reply is: ", string(replyData))
	w.Write(replyData) // Upload JSON back to the client
}

func GetAllFilesFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	SQLQuery := "select * from ObjectFile;"
	var objectTable datasets.ObjectFileTable
	outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	// Convert reflect.Value() to ObjectFileTable, iterate over rows.
	outputData, _ := outputRows.Interface().([]datasets.ObjectFileTable)
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(data)
	}
	jsonResponse(w, err, http.StatusAccepted)
}

func GetFilesFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	pathComponents := strings.Split(r.URL.Path, "/")
	log.Println("pathComponents were", pathComponents)

	// Set up SQL Query depending if request asks for one file or multiple
	var SQLQuery string
	if strings.Contains(params["id"], "-") {
		splitRange := strings.Split(params["id"], "-")

		a, _ := strconv.Atoi(splitRange[0])
		b, _ := strconv.Atoi(splitRange[1])
		SQLQuery = "select * from ObjectFile ORDER BY FileID LIMIT " + strconv.Itoa(b-a+1) + " OFFSET " + strconv.Itoa(a-1)
	} else {
		// if int, use string version for simplicity. No worry about SQL injection since above Atoi didn't fail
		SQLQuery = "select * from ObjectFile where ObjectFile.FileID=" + params["id"]
	}

	var objectTable datasets.ObjectFileTable
	outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	outputData, ok := outputRows.Interface().([]datasets.ObjectFileTable)
	if !ok || len(outputData) == 0 { // Len will be 0 if index is out of range (nothing is returned)
		errMsg := errors.New("Error with ID given, may be out of range of last element")
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return
	}
	var i int
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(data)
		i++
	}
	log.Println(i)
	jsonResponse(w, err, http.StatusAccepted)
}
