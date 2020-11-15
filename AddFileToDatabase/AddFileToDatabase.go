package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	OverheadSQL "github.com/iisharankov/FYSTDatabase/OverheadSQL"
)

const (
	dbUsername    string = "iisharankov"
	dbPassword    string = "iisharankov"
	dbAddress     string = ""
	dbName        string = "mydb"
	sqlTimeLayout string = "2006-01-2 15:04:05"
)

// GlobalPTStackArray is a struct containing an array of structs
var dbCon OverheadSQL.DatabaseConnection

func scanFiles() []string {
	filesList := make([]string, 0)

	cwd, _ := os.Getwd()
	rootDIR := path.Dir(cwd) + "/MockData/"
	matches, _ := filepath.Glob(rootDIR + "FakeFYST/*.txt")
	for _, match := range matches {
		filesList = append(filesList, match)
	}

	return filesList
}

// getFilesOnTelescope takes a string of file data and a field name from ObjectFileTable and outputs a list of those entries
func getFilesOnTelescope(filePaths []string, input string) []string {
	var FileMetadata []string

	for _, val := range filePaths {
		// fmt.Println(val[strings.LastIndex(val, "/")+1:]) // Prints the filename only (last occurrence after "/")
		fileLine := strings.Split(val, ",")
		var objectFileTable OverheadSQL.ObjectFileTable
		val := reflect.Indirect(reflect.ValueOf(objectFileTable))

		r := strings.NewReplacer(" ", "", "'", "") // Remove both whitespace and ' from all data.
		// append the appropriate
		switch input {
		case val.Type().Field(1).Name:
			FileMetadata = append(FileMetadata, strings.ReplaceAll(fileLine[0], "'", "")) // Don't strip whitespace since on date format
		case val.Type().Field(2).Name:
			FileMetadata = append(FileMetadata, r.Replace(fileLine[1]))
		case val.Type().Field(3).Name:
			FileMetadata = append(FileMetadata, r.Replace(fileLine[2]))
		case val.Type().Field(4).Name:
			FileMetadata = append(FileMetadata, r.Replace(fileLine[3]))
		case val.Type().Field(5).Name:
			FileMetadata = append(FileMetadata, r.Replace(fileLine[4]))
		default:
			fmt.Println("That field type not recognized")
			return FileMetadata
		}
	}
	return FileMetadata
}

func exctractData(filepaths []string) []string {
	var fileData []string

	for _, val := range filepaths {
		data, err := ioutil.ReadFile(val)
		if err != nil {
			fmt.Println("File reading error", err)
			var empty []string
			return empty
		}

		cleanedLine := strings.TrimSuffix(string(data), "\n")
		fileData = append(fileData, cleanedLine)
	}
	return fileData
}

func convertDate(layout, dateToConvert string) time.Time {
	t, err := time.Parse(layout, dateToConvert)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func main() {
	// Make connection to Database
	err := dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		fmt.Println(err)
	}

	var timeOfLastFileAdded time.Time
	var objectFileTable OverheadSQL.ObjectFileTable
	// outputRow, err := dbCon.QueryRead("select * from ObjectFile where FileID = (select max(FileID) from Log)", &objectFileTable)
	outputRow, err := dbCon.QueryRead("SELECT * FROM ObjectFile ORDER BY FileID DESC LIMIT 1", &objectFileTable)

	// Convert reflect.Value to struct and find timestamp of element
	valOfVal, _ := outputRow[0].(reflect.Value)
	concreteRow, _ := valOfVal.Interface().(OverheadSQL.ObjectFileTable)
	timeOfLastFileAdded = convertDate(sqlTimeLayout, concreteRow.DateCreated)

	filepaths := scanFiles() // Scan files in "telescope" directory and create list of them
	fileData := exctractData(filepaths)
	listOfFileDates := getFilesOnTelescope(fileData, "DateCreated") // Use scanned list of files and sort by DateCreated field only
	var elementsToAddToTable []int                                  // List of elements that are newer than last DB entry

	// Parse through listOfFileDates to find all elements that have a newer timestamp than timeOfLastFileAdded
	for i, val := range listOfFileDates {
		fileTime := convertDate(sqlTimeLayout, val)
		if fileTime.After(timeOfLastFileAdded) {
			elementsToAddToTable = append(elementsToAddToTable, i)
		}
	}

	fmt.Println("Adding", len(elementsToAddToTable), "files to database")

	//	We now have a list of all scanned files and a list of all elements in that first list that are new
	for _, element := range elementsToAddToTable {
		lineData := strings.Split(fileData[element], ",")

		// Convert each field of the line to the proper format
		date := lineData[0]
		instrumentID, _ := strconv.Atoi(strings.ReplaceAll(lineData[1], " ", ""))
		size, _ := strconv.Atoi(strings.ReplaceAll(lineData[2], " ", ""))
		hashOfBytes := strings.ReplaceAll(lineData[3], " ", "")
		locationOnDisk := "'" + filepaths[element] + "'" // strings.ReplaceAll(lineData[4], " ", "")
		objectStorage := strings.ReplaceAll(lineData[5], " ", "")

		// Find index for FileID column
		temp := struct{ FileID int }{}
		lastFileID, _ := dbCon.QueryRead("SELECT FileID FROM ObjectFile ORDER BY FileID DESC LIMIT 1", &temp)

		// Convert reflect.Value to struct to extract the int
		valOfVal, _ := lastFileID[0].(reflect.Value)
		temp, _ = valOfVal.Interface().(struct{ FileID int })

		// Create the query line to pass to the database
		addQueryLine := fmt.Sprintf("insert into ObjectFile values(%v, %v, %v, %v, %v, %v, %v);", temp.FileID+1, date, instrumentID, size, hashOfBytes, locationOnDisk, objectStorage)
		err := dbCon.QueryWrite(addQueryLine) // Add to Database
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Added FileID row")
		}

	}

}
