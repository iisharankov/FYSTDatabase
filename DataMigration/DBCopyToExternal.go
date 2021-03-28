package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func addRowToLog(FileID, RuleID int, name string) error {
	date := time.Now().Format("2006-01-02 15:04:05")
	log.Println("Adding row to Log table")

	stmt, err := dbCon.PrepareQuery("insert into Log values(?, ?, ?, ?, ?);")
	if err != nil {
		log.Println("Error in db.Perpare()\n", err)
		return err
	}

	// Execute the command on the database (encoded already in stmt)
	_, err = stmt.Exec(FileID, RuleID, date, 0, name)
	if err != nil {
		log.Println("Error in query execution\n", err)
		return err
	}

	return nil
}

// getLocations returns valid/active LocationIDs to copy data for a given InstrumentID
func getLocations(listOfRuleStructs []datasets.RuleTable, InstrumentID int) []int {
	var listOfValidLocationsToCopyTo []int // output list

	for _, val := range listOfRuleStructs {
		// If the InstrumentID matches and the rule is active, add to output list
		if val.InstrumentID == InstrumentID && val.Active == 1 {
			listOfValidLocationsToCopyTo = append(listOfValidLocationsToCopyTo, val.LocationID)
		}
	}
	return listOfValidLocationsToCopyTo
}
