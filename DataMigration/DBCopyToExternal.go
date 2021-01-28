package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type FilesThatNeedToBeBackedUp struct {
	FileID         int
	RuleID         int
	InstrumentID   int
	Size           int
	InstrumentName string
	DateCreated    string
	Storage        string
	ByteHash       string
	LocationName   string
}

func moveOffTelscope() {
	var S3Instance = ObjectMetadata{
		ctx:      context.Background(),
		endpoint: minioEndpoint,
		id:       minioAccessKeyID,
		password: minioSecretAccessKey,
		useSSL:   minioUseSSL}

	S3Instance.initMinio()

	// Make connection to Database if it does not exist or has failed
	if err := dbCon.CheckConnection(); err != nil {
		err := dbCon.Pconnect(dbUsername, dbPassword, dbAddress, dbName)
		if err != nil {
			fmt.Println(err)
		}
	}

	var structOfFilesThatNeedToBeBackedUp FilesThatNeedToBeBackedUp

	// What files have not been copied to locations where active rules exist (pending uploads per all locations)?
	var SQLQuery string = `select f.FileID, r.RuleID, f.Size, i.InstrumentID, i.InstrumentName, f.DateCreated, f.ObjectStorage, f.HashOfBytes, BL.LocationName
	from ObjectFile f 
	join Instrument i on i.InstrumentID=f.InstrumentID
	join Rule r on r.InstrumentID=i.InstrumentID
	join BackupLocation BL on BL.LocationID=r.LocationID
	left join Log l on l.FileID=f.FileID and r.RuleId=l.RuleID
	where l.FileID is null AND r.Active=1 order by f.DateCreated;`
	outputRows, err := dbCon.QueryRead(SQLQuery, &structOfFilesThatNeedToBeBackedUp)
	if err != nil {
		fmt.Println(err)
	}
	listOfFilesThatNeedToBeBackedUp, _ := outputRows.Interface().([]FilesThatNeedToBeBackedUp)

	//--------------------- Finds all the rules in the database

	for _, val := range listOfFilesThatNeedToBeBackedUp {

		var err error
		switch val.RuleID {
		case 1:
			_, err = copyFile(S3Instance, val)
		case 2:
			_, err = copyFile(S3Instance, val)
		case 3:
			_, err = copyFile(S3Instance, val)
		default:
			fmt.Println("That field type not recognized")
			return
		}

		if err == nil {
			addRowToLog(val.FileID, val.RuleID)
		}
	}

}

func addRowToLog(FileID, location int) {
	// date := time.Now().Format("2006-01-02 15:04:05")

	// stmt, err := dbCon.PrepareQuery("insert into Log values(?, ?, ?, ?, ?);")
	// if err != nil {
	// 	log.Println("Error in db.Perpare()\n", err)
	// 	return
	// }

	// // Execute the command on the database (encoded already in stmt)
	// _, err = stmt.Exec(FileID, location, date, 0, "")
	// if err != nil {
	// 	log.Println("Error in query execution\n", err)
	// 	return
	// }

	log.Println("Added row to Log")
}

// getLocations returns valid/active LocationIDs to copy data for a given InstrumentID
func getLocations(listOfRuleStructs []RuleTable, InstrumentID int) []int {
	var listOfValidLocationsToCopyTo []int // output list

	for _, val := range listOfRuleStructs {
		// If the InstrumentID matches and the rule is active, add to output list
		if val.InstrumentID == InstrumentID && val.Active == 1 {
			listOfValidLocationsToCopyTo = append(listOfValidLocationsToCopyTo, val.LocationID)
		}
	}
	return listOfValidLocationsToCopyTo
}

func copyFile(minioInstance ObjectMetadata, src FilesThatNeedToBeBackedUp) (int64, error) {

	location := "us-east-1"
	minioInstance.makeBucket(src.LocationName, location)

	// Upload the zip file
	last := src.Storage[strings.LastIndex(src.Storage, "/")+1:]
	minioInstance.UploadObject(src.LocationName, last, src.Storage, "application/zip")

	return 0, nil
}
