package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"sync"

<<<<<<< HEAD
	_ "github.com/go-sql-driver/mysql" // Used for SQL queries
)

=======
	_ "github.com/go-sql-driver/mysql" // Go wants a comment
)

var dbUsername string = "iisharankov"
var dbPassword string = "iisharankov"
var dbAddress string = ""
var dbName string = "mydb"

>>>>>>> e280b5333df62dc64de7d39bba93c5b883ccf2fc
// GlobalPTStackArray is a struct containing an array of structs
var dbCon DatabaseConnection

// DatabaseConnection is a struct
type DatabaseConnection struct {
	dbConnection *sql.DB
	sync.Mutex
}

// Connect connects to a database
<<<<<<< HEAD
func (dpCon *DatabaseConnection) connect(dbUsername, dbPassword, dbIP, dpName string) error {
=======
func (dpCon *DatabaseConnection) Connect(dbUsername, dbPassword, dbIP, dpName string) error {
>>>>>>> e280b5333df62dc64de7d39bba93c5b883ccf2fc
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

<<<<<<< HEAD
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

=======
// ExecuteQueryWTransaction takes a query and executes it with a transaction for safety
// TODO: make input a list of queries instead of just one!
func (dpCon *DatabaseConnection) ExecuteQueryWTransaction(insertStatement string) {
>>>>>>> e280b5333df62dc64de7d39bba93c5b883ccf2fc
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

<<<<<<< HEAD
// QueryRead takes a query returns a list of all the rows returned by the database
func (dpCon *DatabaseConnection) queryRead(SQLQuery string, p interface{}) ([]interface{}, error) {
	interfacetious := []interface{}{}
=======
// ExecuteQuery takes a query and executes it
func (dpCon *DatabaseConnection) ExecuteQuery(SQLQuery string, p interface{}) ([]interface{}, error) {
	interfacetious := []interface{}{}
	// interfacetious := make([]interface{}, 0)
>>>>>>> e280b5333df62dc64de7d39bba93c5b883ccf2fc

	rows, err := dbCon.dbConnection.Query(SQLQuery)
	defer rows.Close()
	if err != nil {
		return interfacetious, err
	}

	for rows.Next() {
<<<<<<< HEAD
		// reflect.TypeOf(p) gives the pointer to the address of p. We take Elem() to get the value,
		// and make a new copy (reflect.New()) which is also a pointer and must be extracted
		s := reflect.New(reflect.TypeOf(p).Elem()).Elem()

		// Uses reflect to create the correct columns and types from the struct (p) to scan() with.
=======
		// NewP := reflect.Zero(reflect.TypeOf(p).Elem())
		NewP := reflect.New(reflect.TypeOf(p).Elem()).Elem()
		// Elem() is * (kinda)

		fmt.Printf("p is %+T --- %#v ", p, p, p)
		fmt.Println("")
		fmt.Printf("NewP is %+T --- %#v", NewP, NewP, NewP)
		fmt.Println("")

		// Uses reflect to create the correct columns and types from the struct (p) to scan() with.
		s := reflect.ValueOf(p).Elem()
		fmt.Printf("S    is %+T --- %#v", s, s, s)
		fmt.Println("")
		s = NewP
>>>>>>> e280b5333df62dc64de7d39bba93c5b883ccf2fc
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			fmt.Printf("Field is: %T -- %#v\n", field, field)

			columns[i] = field.Addr().Interface()
		}
<<<<<<< HEAD

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
=======
		// Populates the interface p with
		err := rows.Scan(columns...)
		if err != nil {
			return interfacetious, err
		}

		// fmt.Println("p is:", p)
		interfacetious = append(interfacetious, NewP)
		// fmt.Println("RR", interfacetious[0])
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	// fmt.Println("END")
	return interfacetious, nil
}

// // Read takes a query output and a struct, and prints out the next line
// func (dpCon *DatabaseConnection) Read(rows *sql.Rows) (*sql.Rows, interface{}, error) {
// 	rows.Next() // Calls the next row in the sql.Rows type

// 	// Panic if you get any error from error from rows.Err()
// 	err = rows.Err()
// 	if err != nil {
// 		return rows, p, err
// 	}
// 	return rows, p, nil
// }

func main() {
	// Make connection to Database
	err := dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		fmt.Println(err)
	}

	temp := struct {
		RuleDescription string
		FileID          string
	}{}
	var SQLQuery string = "select b.InstrumentName, c.LocationName from Rule a join Instrument b on a.InstrumentID=b.InstrumentID join BackupLocation c on c.LocationID=a.LocationID where c.LocationName='Germany'; "
	outputRows, err := dbCon.ExecuteQuery(SQLQuery, &temp)
	if err != nil {
		fmt.Println("error on ExecuteQuery:", err)
		fmt.Println()
	}
	fmt.Println("YAY", outputRows)
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
	outputRowss, _ := dbCon.ExecuteQuery(FindFilesNotYetAddedQuery, &tempp)

	for i, val := range outputRowss {
		fmt.Println(i, val)
	}
	// hm := outputRows["outputRows"].([]outputRows)
	// fmt.Println(hm[0])
	// u, ok := outputRows["RuleDescription"].([]User)
	// if ok {
	// 	log.Printf("value = %+v\n", u)
	// }

	// if tErr.Error != nil {
	// rows, p, tErr = dbCon.Read(outputRows, &temp)
	// fmt.Println(rows, p)
	// }

	// var instrument InstrumentTable
	// SQLQuery = "select * from Instrument"
	// dbCon.Read(SQLQuery, &instrument)

	// var x string = `insert into Instrument values (3,'TEST','NotAnInstrument', 'Not real',1024,25,900,99)`
	// dbCon.ExecuteQuery(x)
	// dbCon.Read(SQLQuery, &instrument)

	// x = `DELETE FROM Instrument WHERE InstrumentID=3 LIMIT 1;`
	// dbCon.ExecuteQuery(x)
	// dbCon.Read(SQLQuery, &instrument)
>>>>>>> e280b5333df62dc64de7d39bba93c5b883ccf2fc
}
