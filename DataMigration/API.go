package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

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
