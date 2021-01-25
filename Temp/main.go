package main

import (
	"fmt"

	OverheadSQL "github.com/iisharankov/FYSTDatabase/OverheadSQL"
)

var dbCon OverheadSQL.DatabaseConnection

func main() {

	// Make connection to Database
	dbCon.Connect("iisharankov", "iisharankov", "", "mydb")

	temp := struct {
		RuleDescription string
		FileID          string
	}{}
	var SQLQuery string = "select b.InsturmentName, c.LocationName from Rule a join Instrument b on a.InstrumentID=b.InstrumentID join BackupLocation c on c.LocationID=a.LocationID where c.LocationName='Germany'; "
	dbCon.QueryRead(SQLQuery, &temp)

	var instrument OverheadSQL.InstrumentTable
	SQLQuery = "select * from Instrument"
	dbCon.QueryRead(SQLQuery, &instrument)

	var x string = `insert into Instrument values (4,'TEST','NotAnInstrument', 'Not real',1024,25,900,99)`
	dbCon.QueryWrite(x)
	dbCon.QueryRead(SQLQuery, &instrument)

	x = `DELETE FROM Instrument WHERE InstrumentID=3 LIMIT 1;`
	dbCon.QueryWrite(x)
	dbCon.QueryRead(SQLQuery, &instrument)
}

func mainn() {

	// Make connection to Database
	err := dbCon.Connect("iisharankov", "iisharankov", "", "mydb")
	if err != nil {
		fmt.Println(err)
	}

	temp := struct {
		RuleDescription string
		FileID          string
	}{}
	var SQLQuery string = "select b.InstrumentName, c.LocationName from Rule a join Instrument b on a.InstrumentID=b.InstrumentID join BackupLocation c on c.LocationID=a.LocationID where c.LocationName='Germany'; "
	outputRows, err := dbCon.QueryRead(SQLQuery, &temp)
	if err != nil {
		fmt.Println("error on ExecuteQuery:", err)
	}

	for i, val := range outputRows {
		fmt.Println(i, val)
	}

	tempp := struct {
		FileID       int
		Size         int
		InstrumentID string
		Date         string
	}{}
	var FindFilesNotYetAddedQuery string = "select distinct f.FileID, f.Size, i.InstrumentName, f.DateCreated from ObjectFile f join Instrument i on i.InstrumentID=f.InstrumentID left join Log l on l.FileID=f.FileID join Rule r on r.InstrumentID=i.InstrumentID where l.FileID is null AND r.Active=1 order by f.DateCreated;"
	outputRowss, _ := dbCon.QueryRead(FindFilesNotYetAddedQuery, &tempp)

	for i, val := range outputRowss {
		fmt.Println(i, val)
	}

	fmt.Println("------------")

	var x string = `insert into Instrument values (4,'TEST','NotAnInstrument', 'Not real',1024,25,900,99)`
	var Listt []string
	Listt = append(Listt, x)
	dbCon.QueryWriteWithTransaction(Listt)
	fmt.Println("------------")

	var instrument OverheadSQL.instrumentTable
	SQLQuery = "select * from Instrument"
	outputRowss, _ = dbCon.QueryRead(SQLQuery, &instrument)
	for i, val := range outputRowss {
		fmt.Println(i, val)
	}
	fmt.Println("------------")
	y := `DELETE FROM Instrument WHERE InstrumentID=4 LIMIT 1;`
	// xx := `DELETE FROM Instrument WHERE InstrumentID=3 LIMIT 1;`
	yy := []string{y}
	dbCon.QueryWriteWithTransaction(yy)

	fmt.Println("------------")
	outputRowss, _ = dbCon.QueryRead(SQLQuery, &instrument)
	for i, val := range outputRowss {
		fmt.Println(i, val)
	}
}
