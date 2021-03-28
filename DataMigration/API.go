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

// TODO: Generalized this by parsing db table name and number of parameters to upload? might be hard with reflect
func addRowToObjectFile(newFile datasets.File, w http.ResponseWriter) error {
	// Use SQL Prepare() method to safely convert the field types.
	stmt, err := dbCon.PrepareQuery("insert into ObjectFile values(?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Println("Error in db.Perpare()\n", err)
		return err
	}

	// Execute the command on the database (encoded already in stmt)
	_, err = stmt.Exec(nil, newFile.DateCreated, newFile.Instrument, newFile.Size, newFile.MD5Sum, newFile.URL)
	if err != nil {
		log.Println("Error in query execution\n", err)
		return err
	}

	log.Println("Added FileID row to ObjectFile Table")
	return nil
}

func concatErrors(err error, newError string) error {
	var errstrings []string
	errstrings = append(errstrings, err.Error())

	// Concatonate error with additional one.
	errstrings = append(errstrings, fmt.Errorf(newError).Error())

	// combine and return both errors for more useful debugging for user
	combinedErrors := fmt.Errorf(strings.Join(errstrings, " - "))
	return combinedErrors
}

// Shamelessly stolen
func jsonResponse(w http.ResponseWriter, err error, statusCode int) {
	var response struct {
		S string `json:"status"`
		M string `json:"message,omitempty"`
	}

	if err != nil {
		log.Println(err) // Catch all for errors for logging
		response.S = "error"
		response.M = err.Error()
	} else {
		response.S = "ok"
		statusCode = http.StatusOK
	}

	// TODO: Write vs writeHeader?
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print(err)
	}
}

func checkID(endpointSection string, w http.ResponseWriter) error {
	var err error
	// If endpointSection has a '-', it is a range and both elements should be tested recursively
	if strings.Contains(endpointSection, "-") {
		splitRange := strings.Split(endpointSection, "-")
		err = checkID(splitRange[0], w)
		if err != nil {
			return err
		}

		err = checkID(splitRange[1], w)
		if err != nil {
			return err
		}
		return nil
	}

	if endpointSection == "" {
		errMsg := errors.New("No file ID given. Either give file ID after '/'  or remove the '/'")
		return errMsg
	}

	// Check to see if input is an actual integer
	_, err = strconv.Atoi(endpointSection)
	if err != nil {
		errMsg := errors.New("Could not convert ID given to int, check value after 'files/'")
		return errMsg
	}

	return nil
}

func startAPIServer() {

	// Containers fail if db connection is first
	if err := dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName); err != nil {
		log.Println(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/files", AddFileToDatabaseEndpoint).Methods("POST")
	router.HandleFunc("/files", GetAllFilesFromDBEndpoint).Methods("GET")
	router.HandleFunc("/files/{id}", GetFilesFromDBEndpoint).Methods("GET")

	router.HandleFunc("/logs/{id}/{location}", AddLogToDBEndpoint).Methods("POST")
	router.HandleFunc("/logs", GetAllLogsFromDBEndpoint).Methods("GET")
	router.HandleFunc("/logs/{id}", GetLogsFromDBEndpoint).Methods("GET")
	router.HandleFunc("/logs/{id}/copies", GetCopiesOFLogFromDBEndpoint).Methods("GET")

	router.HandleFunc("/rules", GetAllRulesFromDBEndpoint).Methods("GET")
	router.HandleFunc("/rules/{id}", GetRulesFromDBEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":8700", router))
}
