package overheadsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"sync"

	_ "github.com/go-sql-driver/mysql" // maybe unneeded here
)

var dbUsername string = "iisharankov"
var dbPassword string = "iisharankov"
var dbAddress string = "" // "192.168.0.20"
var dbName string = "mydb"

// GlobalPTStackArray is a struct containing an array of structs
var dbCon DatabaseConnection

// InstrumentTable is a SQL Table
type InstrumentTable struct {
	InstrumentID   int
	InstrumentName string
	FullName       string
	Description    string
	NumberOfPixels int
	FrequencyMin   int
	FrequencyMax   int
	TempRange      int
}

// DatabaseConnection is a struct
type DatabaseConnection struct {
	dbConnection *sql.DB
	sync.Mutex
}

// Connect connects to a database
func (dpCon *DatabaseConnection) Connect(dbUsername, dbPassword, dbIP, dpName string) error {

	// "ivan:ivan@tcp(192.168.0.20:3306)/mydb"
	SQLConnectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUsername, dbPassword, dbIP, dpName)
	db, err := sql.Open("mysql", SQLConnectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
		// fmt.Println(err.Error())
	}
	dbCon.dbConnection = db
	return nil
}

// ExecuteQuery takes a query and executes it
func (dpCon *DatabaseConnection) ExecuteQuery(insertStatement string) {
	ctx := context.Background() // Create a new context, and begin a transaction
	tx, err := dbCon.dbConnection.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	} // tx is instance of *sql.Tx where queries can be executed

	_, err = tx.ExecContext(ctx, insertStatement)
	if err != nil {
		tx.Rollback() // rollback transaction if error returned
		return
	}
	err = tx.Commit() // Else commit transaction
	if err != nil {
		log.Fatal(err)
	}
}

// Read sends the query to the Database for evaluation
func (dpCon *DatabaseConnection) Read(SQLQuery string, p interface{}) {
	rows, err := dbCon.dbConnection.Query(SQLQuery)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	// Loop through the rows output
	for rows.Next() {

		// This code uses reflect to create the correct columns and types from the struct (p) to rows.Scan() with
		s := reflect.ValueOf(p).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err = rows.Scan(columns...)
		if err != nil {
			panic(err)
		}
		fmt.Println(p)
	}

	// Panic if you get any error from error from rows.Err()
	err = rows.Err()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("- - - -")
	}
}
