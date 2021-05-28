package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iisharankov/FYSTDatabase/datasets"
)

func concatErrors(err error, newError string) error {
	return fmt.Errorf(newError + " - " + err.Error())
}

// Shamelessly stolen
func jsonResponse(w http.ResponseWriter, err error, statusCode int) {
	var response datasets.ServerHTTPResponse

	// v // Encode message M and Status S depending on statusCode
	if 200 <= statusCode && statusCode < 300 { //
		response.S = "ok"
	} else {
		response.S = "error"
		response.M = err.Error()
	}
	// ^ //

	// v // Write header then populate body
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error while encoding jsonResponse body: ", err)
	}
	// ^ //
}

func startAPIServer() {

	// Containers fail if db connection is first
	if err := dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName); err != nil {
		log.Println(err) // TODO: Panic if database cannot connect
	}

	router := mux.NewRouter()
	router.HandleFunc("/files", AddFileToDatabaseEndpoint).Methods("POST")
	router.HandleFunc("/files", GetFilesFromDBEndpoint).Methods("GET")
	router.HandleFunc("/files/{filename}", GetFileFromDBEndpoint).Methods("GET")
	router.HandleFunc("/files/{filename}/copies", AddRowToCopiesTableEndpoint).Methods("POST")

	router.HandleFunc("/records/{filename}/{location}", AddRecordToDBEndpoint).Methods("POST")
	router.HandleFunc("/records", GetAllRecordsFromDBEndpoint).Methods("GET")
	router.HandleFunc("/records/{filename}", GetRecordFromDBEndpoint).Methods("GET")

	router.HandleFunc("/rules", GetAllRulesFromDBEndpoint).Methods("GET")
	router.HandleFunc("/rules/{id}", GetRulesFromDBEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe(":8700", router))
}
