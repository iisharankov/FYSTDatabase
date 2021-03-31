package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iisharankov/FYSTDatabase/datasets"
)

func AddLogToDBEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Extract the body of the POST request
	var fileName string
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&fileName)

	// Make sure the ID in the header is actually an integer
	FileID, err := strconv.Atoi(params["id"])
	if err != nil {
		combinedErrors := concatErrors(err, "File ID could not be converted to an int!")
		jsonResponse(w, combinedErrors, http.StatusServiceUnavailable)
		return
	}

	// Tells you which RuleID corresponds to a FileID and Location string (locationName in BackupLocation)
	query := fmt.Sprintf(`select r.RuleID from Rule r 
		join ObjectFile o on o.InstrumentID=r.InstrumentID 
		join BackupLocation b on b.LocationID=r.LocationID
		where o.FileId = %v and b.S3Bucket = "%v"`, params["id"], params["location"])
	queryReturn, err := dbCon.QueryRead(query, &struct{ RuleID int }{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}
	returnQuery, _ := queryReturn.Interface().([]struct{ RuleID int }) // type assert from reflect.Value to struct

	if err = addRowToLog(FileID, returnQuery[0].RuleID, fileName); err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	jsonResponse(w, nil, http.StatusAccepted)
}

func GetCopiesOFLogFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Check to see if input is an actual integer
	_, err := strconv.Atoi(params["id"])
	if err != nil {
		errMsg := errors.New("Could not convert ID given to int, check value after 'files/'")
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return
	}

	// if int, use string version for simplicity
	SQLQuery := "select * from Log where Log.FileID=" + params["id"]
	var objectTable datasets.LogTable
	outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	outputData, _ := outputRows.Interface().([]datasets.LogTable)
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(data)
	}
	jsonResponse(w, err, http.StatusAccepted)
}

func GetLogsFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Set up SQL Query depending if request asks for one file or multiple
	var SQLQuery string
	if strings.Contains(params["id"], "-") {
		splitRange := strings.Split(params["id"], "-")

		a, _ := strconv.Atoi(splitRange[0])
		b, _ := strconv.Atoi(splitRange[1])
		SQLQuery = "select * from Log ORDER BY FileID LIMIT " + strconv.Itoa(b-a+1) + " OFFSET " + strconv.Itoa(a-1)
	} else {
		// if int, use string version for simplicity. No worry about SQL injection since above Atoi didn't fail
		SQLQuery = "select * from Log where Log.FileID=" + params["id"]
	}

	var logTable datasets.LogTable
	outputRows, err := dbCon.QueryRead(SQLQuery, &logTable)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	outputData, ok := outputRows.Interface().([]datasets.LogTable)
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

func GetAllLogsFromDBEndpoint(w http.ResponseWriter, r *http.Request) {

	// if int, use string version for simplicity
	SQLQuery := "select * from Log"
	var objectTable datasets.LogTable
	outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	outputData, _ := outputRows.Interface().([]datasets.LogTable)
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(data)
	}
	jsonResponse(w, err, http.StatusAccepted)
}
