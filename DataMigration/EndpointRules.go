package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iisharankov/FYSTDatabase/datasets"
)

func GetRulesFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "-", r.URL)
	params := mux.Vars(r)

	// v // Set up SQL Query depending if request asks for one file or multiple
	var SQLQuery string
	if strings.Contains(params["id"], "-") { // Hyphen suggests range of values
		splitRange := strings.Split(params["id"], "-")
		a, _ := strconv.Atoi(splitRange[0])
		b, _ := strconv.Atoi(splitRange[1])
		SQLQuery = "select * from Rules ORDER BY RuleID LIMIT " + strconv.Itoa(b-a+1) + " OFFSET " + strconv.Itoa(a-1)

	} else { // Single row requested
		idAsInt, _ := strconv.Atoi(params["id"]) // to avoid SQL injection
		SQLQuery = "select * from Rules where Rule.RuleID=" + string(idAsInt)
	}
	// ^ //

	// v // Send query to database
	outputRows, err := dbCon.QueryRead(SQLQuery, &datasets.RulesTable{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}
	// ^ //

	// v // Convert reflect.Value() to RulesTable and check if empty or server error
	outputData, ok := outputRows.Interface().([]datasets.RulesTable)
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response"), http.StatusInternalServerError)
		return
	} else if len(outputData) == 0 { // Will only trigger when single row rested
		jsonResponse(w, fmt.Errorf("Rule '%v' does not exist in database", params["id"]), http.StatusBadRequest)
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

func GetAllRulesFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "-", r.URL)

	// v // Send query to database requesting all the rules
	outputRows, err := dbCon.QueryRead("select * from Rules ORDER BY FileID DESC LIMIT 50;", &datasets.RulesTable{})
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}
	// ^ //

	// v // Convert reflect.Value() to RulesTable and check if empty or server error
	outputData, ok := outputRows.Interface().([]datasets.RulesTable)
	if !ok {
		jsonResponse(w, fmt.Errorf("Server error in casting database response"), http.StatusInternalServerError)
		return
	} else if len(outputData) == 0 { // Will only trigger when single row rested
		jsonResponse(w, fmt.Errorf("No rules in database to return"), http.StatusBadRequest)
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
