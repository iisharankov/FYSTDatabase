package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iisharankov/FYSTDatabase/datasets"
)

func GetRulesFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Set up SQL Query depending if request asks for one file or multiple
	var SQLQuery string
	if strings.Contains(params["id"], "-") {
		splitRange := strings.Split(params["id"], "-")

		a, _ := strconv.Atoi(splitRange[0])
		b, _ := strconv.Atoi(splitRange[1])
		SQLQuery = "select * from Rule ORDER BY RuleID LIMIT " + strconv.Itoa(b-a+1) + " OFFSET " + strconv.Itoa(a-1)
	} else {
		// if int, use string version for simplicity. No worry about SQL injection since above Atoi didn't fail
		SQLQuery = "select * from Rule where Rule.RuleID=" + params["id"]
	}

	var ruleTable datasets.RuleTable
	outputRows, err := dbCon.QueryRead(SQLQuery, &ruleTable)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	outputData, ok := outputRows.Interface().([]datasets.RuleTable)
	if !ok || len(outputData) == 0 { // Len will be 0 if index is out of range (nothing is returned)
		errMsg := errors.New("Error with ID given, may be out of range of last element")
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return
	}
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(data)
	}

	jsonResponse(w, err, http.StatusAccepted)
}

func GetAllRulesFromDBEndpoint(w http.ResponseWriter, r *http.Request) {
	var ruleTable datasets.RuleTable
	outputRows, err := dbCon.QueryRead("select * from Rule;", &ruleTable)
	if err != nil {
		jsonResponse(w, err, http.StatusBadRequest)
		return
	}

	outputData, _ := outputRows.Interface().([]datasets.RuleTable)
	for _, val := range outputData {
		data, _ := json.Marshal(val)
		w.Write(data)
	}
	jsonResponse(w, err, http.StatusAccepted)
	return
}
