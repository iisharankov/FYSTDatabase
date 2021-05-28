package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iisharankov/FYSTDatabase/datasets"
)

func AddFileToDatabaseEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("AddFileToDatabaseEndpoint -", r.URL)

	var newFile datasets.File
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&newFile)
	log.Println("new file: ", newFile)

	// v // Check that  both FileName and MD5 are provided. Database checks FileName uniqueness
	if newFile.Name == "" || newFile.MD5Sum == "" {
		jsonResponse(w, fmt.Errorf("Name and MD5sum cannot be empty"), http.StatusConflict)
		return
	}
	// ^ //

	// v // Add file to Files table. PrepareQuery() method safely converts the field types.
	log.Println("Number of fields", reflect.ValueOf(datasets.FilesTable{}).NumField())
	if stmt, err := dbCon.PrepareQuery("insert into Files values(?, ?, ?, ?, ?, ?);"); err != nil {
		errMsg := fmt.Sprintf("Error while preparing File row insert: %v", err)
		jsonResponse(w, fmt.Errorf(errMsg), http.StatusBadRequest)
		return

	} else if _, err = stmt.Exec(nil, newFile.Name, newFile.DateCreated, newFile.Instrument,
		newFile.Size, newFile.MD5Sum); err != nil { // Wish you could easily unpack struct.
		errMsg := fmt.Sprintf("Error while executing File insert: %v", err)
		jsonResponse(w, fmt.Errorf(errMsg), http.StatusBadRequest)
		return
	}
	// ^ //

	// v // Find location row in LocationsTable to return metadata in response with
	var locationsTableRow []datasets.LocationsTable
	query := `select * from Locations b where b.LocationID=1` // LocationID=1 stands for Local object storage
	queryBackupLocation, err := dbCon.QueryRead(query, &datasets.LocationsTable{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}
	// ^ //

	// v // Convert reflect.Value() to ObjectFileTable and check for various server errors
	locationsTableRow, ok := queryBackupLocation.Interface().([]datasets.LocationsTable)
	log.Println(locationsTableRow)
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response."), http.StatusInternalServerError)
		return
	} else if len(locationsTableRow) == 0 {
		jsonResponse(w, fmt.Errorf("No location in database for local object storage."), http.StatusInternalServerError)
		return
	} else if locationsTableRow[0].S3Bucket == "" {
		jsonResponse(w, fmt.Errorf("Database failed to give bucket for upload."), http.StatusInternalServerError)
		return
	}
	// ^ //

	// Connect to local object storage to access metadata to return to client
	temp, err := connectToObjectStorage("Observatory")

	// v // If successful, respond to client with upload location details
	response := datasets.ClientUploadReply{
		S3Metadata: datasets.S3Metadata{
			Endpoint:        temp.address,
			AccessKeyID:     temp.accessID,
			SecretAccessKey: temp.secretID,
			UseSSL:          strconv.FormatBool(temp.useSSL),
		},
		FileName:       newFile.Name,
		LocationID:     locationsTableRow[0].LocationID,
		UploadLocation: locationsTableRow[0].S3Bucket,
	}

	if replyData, err := json.Marshal(response); err != nil {
		jsonResponse(w, fmt.Errorf("Server error packaging return response"), http.StatusInternalServerError)
		return
	} else {
		w.Write(append(replyData, "\n"...)) // Return JSON back to client
	}
	// ^ //
}

func GetFilesFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFilesFromDBEndpoint -", r.URL)

	outputRows, err := dbCon.QueryRead("select * from Files ORDER BY FileID DESC LIMIT 50;", &datasets.FilesTable{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	// v // Convert reflect.Value() to ObjectFileTable and check if empty or server error
	outputData, ok := outputRows.Interface().([]datasets.FilesTable)
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response"), http.StatusInternalServerError)
		return
	} else if len(outputData) == 0 {
		errMsg := fmt.Errorf("No files in database to return")
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return
	}
	// ^ //

	// v // Print out values within outputData
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(append(data, "\n"...))
	}
	// ^ //
}

func GetFileFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFileFromDBEndpoint -", r.URL)
	params := mux.Vars(r)

	// v // Request database to return file with the given Filename
	SQLQuery := fmt.Sprintf(`select * from Files where FileName="%v"`, params["filename"])
	outputRows, err := dbCon.QueryRead(SQLQuery, &datasets.FilesTable{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}
	// ^ //

	// v // Convert reflect.Value() to FilesTable and check if empty or server error
	outputData, ok := outputRows.Interface().([]datasets.FilesTable)
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response"), http.StatusInternalServerError)
		return
	} else if len(outputData) == 0 {
		errMsg := fmt.Errorf("File '%v' does not exist in database", params["filename"])
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return
	}
	// ^ //

	// v // Print out values within outputData (should be len 1 but for safety loop)
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(append(data, "\n"...))
	}
	// ^ //
}

func AddRowToCopiesTableEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("AddRowToCopiesTableEndpoint -", r.URL)
	params := mux.Vars(r)

	// v // Extract the body form the request, which should contain an int, which is the ruleID
	var ruleID int
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&ruleID)
	// ^ //

	// Find the FileID for the given FileName
	fileID, err := getFileIDFromFileName(params["filename"])
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	// v // We do not want to add a row to the Copies table when this endpoint is triggered if the
	// referenced file does not exist in the Records table. So the next query asks the database if
	// a row exists with the given FileID and RuleID.
	var recordsRow datasets.RecordsTable
	query := fmt.Sprintf(`select * FROM Records  WHERE FileID=%v and RuleID=%v;`, fileID, ruleID)
	if queryReturn, err := dbCon.QueryRead(query, &datasets.RecordsTable{}); err != nil {
		log.Println(err)
	} else if temp := queryReturn.Interface().([]datasets.RecordsTable); len(temp) == 1 {
		recordsRow = temp[0] // Above line makes sure the return query is of length one
	} else {
		log.Println(err)
	}
	// ^ //

	// v // Now check if struct is empty (means no row returned, I.E. no such row in Records)
	if (datasets.RecordsTable{}) == recordsRow {
		err := fmt.Errorf("No row in Records for %v. Cannot add row to Copies", params["filename"])
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}
	// ^ //

	// v // If this section is reached, we know a row exists in Records for this file, we can add it to Copies
	stmt, err := dbCon.PrepareQuery(`INSERT INTO Copies Values(?, ?, ?);`) // Try to add a row to Copies
	if err != nil {
		newErr := concatErrors(err, "Error in db.Perpare() while adding row to Copies")
		jsonResponse(w, newErr, http.StatusBadRequest)
		return
	} else if _, err = stmt.Exec(fileID, ruleID, params["filename"]); err != nil {
		newErr := concatErrors(err, "Error in query execution while adding row to Copies")
		jsonResponse(w, newErr, http.StatusBadRequest)
		return
	}
	// ^ //
	jsonResponse(w, nil, http.StatusOK)
}
