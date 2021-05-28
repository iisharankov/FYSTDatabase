package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iisharankov/FYSTDatabase/datasets"
)

func AddRecordToDBEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("AddRecordToDBEndpoint -", r.URL)
	params := mux.Vars(r)

	log.Println(r.URL, "1")

	// check that the MD5 of the object on Minio is the same as the given one in the database
	if err := verifyObject(params["filename"]); err != nil {
		errMsg := fmt.Errorf("Error verifying object at Files endpoint: " + err.Error())
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return
	}

	log.Println(r.URL, "2")

	// Finds which RuleID corresponds to FileName & Location string given
	query := fmt.Sprintf(`select r.RuleID from Rules r 
		join Files f on f.InstrumentID=r.InstrumentID join Locations l on l.LocationID=r.LocationID
		where f.FileName="%v" and l.S3Bucket="%v" and r.Active=1;`, params["filename"], params["location"])

	outputRows, err := dbCon.QueryRead(query, &struct{ RuleID int }{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	log.Println(r.URL, "3")

	// Convert reflect.Value() to RecordsTable and check if empty or server error
	ruleIDRow, ok := outputRows.Interface().([]struct{ RuleID int })
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response"), http.StatusInternalServerError)
		return
	} else if len(ruleIDRow) == 0 {
		// Will only trigger if file has no local object storage rule, or rule is not active
		jsonResponse(w, fmt.Errorf("No records in database to return"), http.StatusBadRequest)
		return
	}

	log.Println(r.URL, "4")

	// Find the FileID to add a entry to record, since we only were given FileName
	FileID, err := getFileIDFromFileName(params["filename"])
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	log.Println(r.URL, "5")

	// Attempt to add the log to the Records table
	if err = addRowToRecord(FileID, ruleIDRow[0].RuleID); err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	jsonResponse(w, nil, http.StatusOK)
}

func GetRecordFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("GetRecordFromDBEndpoint -", r.URL)
	params := mux.Vars(r)

	// Find the fileID for the given filename
	fileID, err := getFileIDFromFileName(params["filename"])
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	// Send query to database asking for records corresponding to fileID
	SQLQuery := fmt.Sprintf(`select * from Records where FileID=%v;`, fileID)
	outputRows, err := dbCon.QueryRead(SQLQuery, &datasets.RecordsTable{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	// Convert reflect.Value() to RecordsTable and check if empty or server error
	outputData, ok := outputRows.Interface().([]datasets.RecordsTable)
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response"), http.StatusInternalServerError)
		return
	} else if len(outputData) == 0 { // Will only trigger when single row rested
		jsonResponse(w, fmt.Errorf("Record '%v' does not exist in database", params["filename"]), http.StatusBadRequest)
		return
	}

	// Print out values within outputData
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(append(data, "\n"...))
	}

}

func GetAllRecordsFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllRecordsFromDBEndpoint -", r.URL)

	// Send query to database requesting all records
	SQLQuery := "select * from Records ORDER BY FileID DESC LIMIT 50;"
	outputRows, err := dbCon.QueryRead(SQLQuery, &datasets.RecordsTable{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	// Convert reflect.Value() to RecordsTable and check if empty or server error
	outputData, ok := outputRows.Interface().([]datasets.RecordsTable)
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response"), http.StatusInternalServerError)
		return
	} else if len(outputData) == 0 { // Will only trigger when single row rested
		jsonResponse(w, fmt.Errorf("No records in database to return"), http.StatusBadRequest)
		return
	}

	// Print out values within outputData
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(append(data, "\n"...))
	}

}
