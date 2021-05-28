package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// These are actions performed several times in multiple places. File title needs renaming

func addRowToRecord(FileID, RuleID int) error {
	date := time.Now().Format(sqlTimeLayout)

	stmt, err := dbCon.PrepareQuery("insert into Records values(?, ?, ?);")
	if err != nil {
		log.Println("Error in db.Perpare()\n", err)
		return err
	}

	// Execute the command on the database (encoded already in stmt)
	_, err = stmt.Exec(FileID, RuleID, date)
	if err != nil {
		log.Println("Error in query execution\n", err)
		return err
	}

	return nil
}

func getFileIDFromFileName(filename string) (int, error) {

	// v // Send query to database to find FileID corresponding to filename given
	query := fmt.Sprintf(`select FileID from Files where FileName="%v";`, filename)
	queryReturn, err := dbCon.QueryRead(query, &struct{ FileID int }{})
	if err != nil {
		return 0, fmt.Errorf("getFileIDFromFileName database query failed")
	}
	// ^ //

	// v // Convert reflect.Value() to proper struct and check if empty or server error
	temp, ok := queryReturn.Interface().([]struct{ FileID int })
	if !ok {
		return 0, fmt.Errorf("Server error in casting database response")
	} else if len(temp) == 0 { // Empty response assumes FileName does not exist
		return 0, fmt.Errorf("File '%v' does not exist in database", filename)
	}
	// ^ //

	return temp[0].FileID, nil
}
