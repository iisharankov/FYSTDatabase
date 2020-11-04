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


// GlobalPTStackArray is a struct containing an array of structs
var dbCon DatabaseConnection

// DatabaseConnection is a struct
type DatabaseConnection struct {
	dbConnection *sql.DB
	sync.Mutex
}

// Connect connects to a database
func (dpCon *DatabaseConnection) connect(dbUsername, dbPassword, dbIP, dpName string) error {
	SQLConnectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUsername, dbPassword, dbIP, dpName)
	db, err := sql.Open("mysql", SQLConnectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	dbCon.dbConnection = db
	return nil
}


func (dpCon *DatabaseConnection) checkConnection() error {

	err := dbCon.dbConnection.Ping()
	if err != nil {
		return err
	}
	return nil
}

// QueryWriteWithTransaction takes a query and executes it with a transaction for safety
// TODO: make input a list of queries instead of just one!
func (dpCon *DatabaseConnection) queryWriteWithTransaction(insertStatement []string) {

	ctx := context.Background() // Create a new context, and begin a transaction
	tx, err := dbCon.dbConnection.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// tx is instance of *sql.Tx where queries can be executed
	for _, val := range insertStatement {
		_, err = tx.ExecContext(ctx, val)
		if err != nil {
			tx.Rollback() // rollback transaction if error returned
			return
		}
	}

	// tx.Commit() will fail if tx.Rollback() is called in above loop,
	// so this is safe to leave outside the loop.
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// QueryWrite takes a single query and executes it with no transactional safety
func (dpCon *DatabaseConnection) queryWrite(insertStatement string) {
	_, err := dbCon.dbConnection.Query(insertStatement)
	if err != nil {
		fmt.Println(err)
	}
}


// QueryRead takes a query returns a list of all the rows returned by the database
func (dpCon *DatabaseConnection) queryRead(SQLQuery string, p interface{}) ([]interface{}, error) {
	interfacetious := []interface{}{}
	rows, err := dbCon.dbConnection.Query(SQLQuery)
	defer rows.Close()
	if err != nil {
		return interfacetious, err
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
			fmt.Printf("Field is: %T -- %#v\n", field, field)

			columns[i] = field.Addr().Interface()
		}

		// Scans the next query row and populates columns, which has pointers to the memory addresses.
		if err := rows.Scan(columns...); err != nil {
			return []interface{}{}, err
		}
		interfacetious = append(interfacetious, s)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return interfacetious, nil
}
