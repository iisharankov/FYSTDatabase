package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Used for SQL queries
)

// DatabaseConnection is a struct
type DatabaseConnection struct {
	DBConnection *sql.DB
	sync.Mutex
}

// Connect connects to a database
func (dpCon *DatabaseConnection) Connect(dbUsername, dbPassword, dbIP, dpName string) error {
	SQLConnectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUsername, dbPassword, dbIP, dpName)

	db, err := sql.Open("mysql", SQLConnectionString)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	dbCon.DBConnection = db
	return nil
}

// // Connect connects to a database
// func (dpCon *DatabaseConnection) Connect(dbUsername, dbPassword, dbIP, dpName string) (emptyDB *sql.DB, err error) {
// 	SQLConnectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUsername, dbPassword, dbIP, dpName)
// 	db, err := sql.Open("mysql", SQLConnectionString)
// 	if err != nil {
// 		return emptyDB, err
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		return emptyDB, err
// 	}
// 	dbCon.DBConnection = db
// 	return dbCon.DBConnection, nil
// }

// CheckConnection checks if the connection to the database is active, reconnect if necessary
func (dpCon *DatabaseConnection) CheckConnection() (err error) {

	// Connection may be nil if never established
	if dbCon.DBConnection == nil {
		log.Println("No connection to database established, attempting to connect")
		if dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName); err != nil {
			log.Println("Database reconnect failed")
			return err
		}
	}

	if dbCon.DBConnection.Ping(); err != nil {
		log.Println("Ping to Database returned error, reconnecting")
		dbCon.DBConnection = nil

		if dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName); err != nil {
			log.Println("Database reconnect failed")
			return err
		}
	}
	return nil
}

// PrepareQuery prepares a query with the given string
func (dpCon *DatabaseConnection) PrepareQuery(prepareString string) (stmt *sql.Stmt, err error) {
	if dbCon.CheckConnection(); err != nil {
		return nil, err
	}

	stmt, err = dbCon.DBConnection.Prepare(prepareString)
	return stmt, err
}

// QueryWriteWithTransaction takes a query and executes it with a transaction for safety
func (dpCon *DatabaseConnection) QueryWriteWithTransaction(insertStatement []string) (err error) {
	if dbCon.CheckConnection(); err != nil {
		return err
	}

	ctx := context.Background() // Create a new context, and begin a transaction
	tx, err := dbCon.DBConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// tx is instance of *sql.Tx where queries can be executed
	for _, val := range insertStatement {
		_, err = tx.ExecContext(ctx, val)
		if err != nil {
			tx.Rollback() // rollback transaction if error returned
			return err
		}
	}

	// tx.Commit() will fail if tx.Rollback() is called in above loop,
	// so this is safe to leave outside the loop.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// QueryWrite takes a single query and executes it with no transactional safety
func (dpCon *DatabaseConnection) QueryWrite(insertStatement string) (err error) {
	if dbCon.CheckConnection(); err != nil {
		return err
	}

	_, err = dbCon.DBConnection.Query(insertStatement)
	if err != nil {
		return err
	}
	return nil
}

// QueryRead takes a query returns a list of all the rows returned by the database
func (dpCon *DatabaseConnection) QueryRead(SQLQuery string, p interface{}) (emptyOutput reflect.Value, err error) {
	if dbCon.CheckConnection(); err != nil {
		return emptyOutput, err
	}

	// Here we are relying heavily on the reflect package to create a generic function for all types of SQL tables
	// The output is a slice of Values that can be then type asserted to the given structure p (outside of this function)
	outputData := reflect.New(reflect.SliceOf(reflect.TypeOf(p).Elem())).Elem()

	// Send query to database to receive rows
	rows, err := dbCon.DBConnection.Query(SQLQuery)
	defer rows.Close()
	if err != nil {
		return emptyOutput, err
	}

	for rows.Next() {
		// reflect.TypeOf(p) gives the pointer to the address of p. We take Elem() to get the value,
		// and make a new copy (reflect.New()) which is also a pointer and must be extracted
		s := reflect.New(reflect.TypeOf(p).Elem()).Elem()

		// Uses reflect to create the correct columns and types from the struct (p) to scan() with.
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}
		// Scans the next query row and populates columns, which has pointers to the memory addresses.
		if err := rows.Scan(columns...); err != nil {
			return emptyOutput, err
		}
		outputData = reflect.Append(outputData, s)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	return outputData, nil
}
