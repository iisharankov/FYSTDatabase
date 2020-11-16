package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	OverheadSQL "github.com/iisharankov/FYSTDatabase/OverheadSQL"
)

type backupLocation struct {
	LocationID      int
	LocationName    string
	LocationAddress string
	FTPIPAddress    string
}

var dbUsername string = "iisharankov"
var dbPassword string = "iisharankov"
var dbAddress string = ""
var dbName string = "mydb"

// GlobalPTStackArray is a struct containing an array of structs
var dbCon OverheadSQL.DatabaseConnection

func main() {
	// Make connection to Database
	err := dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		fmt.Println(err)
	}

	var outputData []OverheadSQL.ObjectFileTable
	var objectFileTable OverheadSQL.ObjectFileTable

	var SQLQuery string = "select distinct f.* from ObjectFile f join Instrument i on i.InstrumentID=f.InstrumentID left join Log l on l.FileID=f.FileID join Rule r on r.InstrumentID=i.InstrumentID where l.FileID is null AND r.Active=1 order by f.DateCreated;"
	outputRows, err := dbCon.QueryRead(SQLQuery, &objectFileTable)
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range outputRows {
		// Convert reflect.Value to a OverheadSQL.ObjectFileTable struct and append to list of structs
		valOfVal, _ := val.(reflect.Value)
		concreteRow, _ := valOfVal.Interface().(OverheadSQL.ObjectFileTable)
		outputData = append(outputData, concreteRow)
	}

	//--------------------- Finds all the rules in the database
	var listOfRuleStructs []OverheadSQL.RuleTable
	var ruleTable OverheadSQL.RuleTable
	outputRows, err = dbCon.QueryRead("select * from Rule", &ruleTable)
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range outputRows {
		// Convert reflect.Value to a OverheadSQL.ObjectFileTable struct and append to list of structs
		valOfVal, _ := val.(reflect.Value)
		concreteRow, _ := valOfVal.Interface().(OverheadSQL.RuleTable)
		listOfRuleStructs = append(listOfRuleStructs, concreteRow)
	}

	//---------------------
	// type Temp struct {
	// 	InstrumentName string
	// 	LocationID     int
	// 	LocationName   string
	// }
	// var temp Temp
	// for _, val := range outputData {
	// 	fmt.Println("-")
	// 	SQLQuery = "select b.InstrumentName, a.LocationID, c.LocationName from Rule a join Instrument b on a.InstrumentID=b.InstrumentID join BackupLocation c on c.LocationID=a.LocationID where a.active=1 and b.InstrumentID=" + strconv.Itoa(val.InstrumentID)
	// 	outputRows, err = dbCon.QueryRead(SQLQuery, &temp)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	for _, val := range outputRows {
	// 		// Convert reflect.Value to a OverheadSQL.ObjectFileTable struct
	// 		valOfVal, _ := val.(reflect.Value)
	// 		concreteRow, _ := valOfVal.Interface().(Temp)
	// 		fmt.Println(concreteRow)
	// 	}
	// }

	for _, val := range outputData {
		locToCopyTo := getLocations(listOfRuleStructs, val.InstrumentID)
		for _, location := range locToCopyTo {

			var err error
			switch location {
			case 1:
				_, err = copyFile(val.LocationOnDisk, "Toronto")
			case 2:
				_, err = copyFile(val.LocationOnDisk, "Cornell")
			case 3:
				_, err = copyFile(val.LocationOnDisk, "Germany")
			default:
				fmt.Println("That field type not recognized")
				return
			}

			if err == nil {
				date := time.Now().Format("2006-01-02 15:04:05")
				addQueryLine := fmt.Sprintf("insert into Log values(%v, %v, '%v', %v);", val.FileID, location, date, 0)
				err = dbCon.QueryWrite(addQueryLine)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Added row in Log")
				}
			}
		}
	}

}

// --- USE QUERY TO ADD LOG TO LOG TABLE OTHERWISE COPY WIL HAPPEN AGAIN

// getLocations returns valid/active LocationIDs to copy data for a given InstrumentID
func getLocations(listOfRuleStructs []OverheadSQL.RuleTable, InstrumentID int) []int {
	var listOfValidLocationsToCopyTo []int // output list

	for _, val := range listOfRuleStructs {
		// If the InstrumentID matches and the rule is active, add to output list
		if val.InstrumentID == InstrumentID && val.Active == 1 {
			listOfValidLocationsToCopyTo = append(listOfValidLocationsToCopyTo, val.LocationID)
		}
	}
	return listOfValidLocationsToCopyTo
}

func copyFile(src, location string) (int64, error) {

	a, filename := filepath.Split(src)
	basePath := path.Dir(path.Dir(a))
	dst := basePath + "/" + location + "/" + filename

	// Opens source file and creates destination file
	source, _ := os.Open(src)
	defer source.Close()
	destination, _ := os.Create(dst)
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
