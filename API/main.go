package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type File struct {
	Name        string    `json:"name"`
	MD5Sum      string    `json:"md5sum"`
	DateCreated time.Time `json:"date_created"`
	Size        int       `json:"size"`
	Location    string    `json:"location"`
}

// type Command interface {
// 	Check() error
// }

// func (file File) Check() error {
// 	// XXX:TBD
// 	return nil
// }
func returnQueryParameter(r *http.Request, name string) (string, error) {
	parameterToReturn, ok := r.URL.Query()[name]
	if !ok {
		errMsg := fmt.Sprintf(`{"err": "r.URL.Query() returned false. Check parameter %v"}`, name)
		return "", fmt.Errorf(errMsg)
	} else if len(parameterToReturn) != 1 {
		errMsg := fmt.Sprintf(`{"err": "%v parameters found in URL query for '%v', only 1 accepted "}`, len(parameterToReturn), name)
		return "", fmt.Errorf(errMsg)
	}
	return parameterToReturn[0], nil
}

// filesEndpoint listens to the http request header for a curl command to update the IP/Port of the filesing, or to enable/disable
func filesEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	var statusCode int

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	endpoint := r.URL.Path
	endpointSections := strings.Split(endpoint, "/")

	if r.Method == "POST" {
		if endpointSections[1] == "files" {
			if len(endpointSections) == 2 {
				var x File
				err = dec.Decode(&x)
				statusCode = http.StatusAccepted
				diffFunc(x)
			}
		} else {
			statusCode = http.StatusNotFound
		}

		if err != nil {
			statusCode = http.StatusBadRequest
		}

		_, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			fmt.Println(err)
			return
		}

		// r.GetBody()
		// fmt.Println("Body is", r.Body)
		// r.ParseForm()
		// fmt.Println("form is", r.Form)

		// if r.Form["Module"][0] != "Services.PositionFiles"
		// fmt.Println(cmd)

	} else if r.Method == "GET" {
		if endpointSections[1] == "files" {
			if len(endpointSections) == 2 {
				SQLQuery := "select * from ObjectFile;"
				var objectTable OverheadSQL.ObjectFileTable
				outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
				if err != nil {
					fmt.Println(err)
				}
				outputData, _ := outputRows.Interface().([]OverheadSQL.ObjectFileTable)
				for _, val := range outputData {
					data, _ := json.Marshal(val)
					fmt.Fprintln(w, string(data))
				}
				statusCode = http.StatusAccepted

			} else if len(endpointSections) == 3 {
				if endpointSections[2] == "" {
					errMsg := errors.New("No file ID given. Either give file ID after '/'  or remove the '/'")
					jsonResponse(w, errMsg, http.StatusBadRequest)
					return
				}
				// Check to see if input is an actual integer
				_, err := strconv.Atoi(endpointSections[2])
				if err != nil {
					errMsg := errors.New("Could not convert ID given to int, check value after 'files/'")
					jsonResponse(w, errMsg, http.StatusBadRequest)
					return
				}

				// if int, use string version for simplicity. No worry about SQL injection since above Atoi didn't fail
				SQLQuery := "select * from ObjectFile where ObjectFile.FileID=" + endpointSections[2]
				var objectTable OverheadSQL.ObjectFileTable
				outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
				if err != nil {
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				outputData, ok := outputRows.Interface().([]OverheadSQL.ObjectFileTable)
				if !ok || len(outputData) == 0 { // Len will be 0 if index is out of range (nothing is returned)
					errMsg := errors.New("Error with ID given, may be out of range of last element")
					jsonResponse(w, errMsg, http.StatusBadRequest)
					return
				}

				data, _ := json.Marshal(outputData[0])
				fmt.Fprintln(w, string(data)) // Should only be a single output
				statusCode = http.StatusAccepted

			} else if len(endpointSections) == 4 && endpointSections[3] == "copies" {
				// Check to see if input is an actual integer
				_, err := strconv.Atoi(endpointSections[2])
				if err != nil {
					errMsg := errors.New("Could not convert ID given to int, check value after 'files/'")
					jsonResponse(w, errMsg, http.StatusBadRequest)
					return

				}

				// if int, use string version for simplicity
				SQLQuery := "select * from Log where Log.FileID=" + endpointSections[2]
				var objectTable OverheadSQL.LogTable
				outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
				if err != nil {
					fmt.Println(err)
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				outputData, _ := outputRows.Interface().([]OverheadSQL.LogTable)
				for _, val := range outputData {
					data, _ := json.Marshal(val)
					fmt.Fprintln(w, string(data))
				}
				statusCode = http.StatusAccepted
			}
		}
	} else {
		statusCode = http.StatusNotFound
	}
	jsonResponse(w, err, statusCode)
}

func diffFunc(cmd File) {
	// Find largest index for FileID column, so we can increment by one
	temp := struct{ FileID int }{} // Temp struct that has just integer (SQL query returns row of 1 int)
	queryReturn, _ := dbCon.QueryRead("SELECT FileID FROM ObjectFile ORDER BY FileID DESC LIMIT 1", &temp)
	lastFileID, _ := queryReturn.Interface().([]struct{ FileID int }) // type assert from reflect.Value to struct
	FileID := lastFileID[0].FileID + 1                                // Takes the first element out, and then takes the FileID field

	// Use SQL Prepare() method to safely convert the field types.
	stmt, err := dbCon.DBConnection.Prepare("insert into ObjectFile values(?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Println("Error in db.Perpare()\n", err)
	}

	// Execute the command on the database (encoded already in stmt)
	_, err = stmt.Exec(FileID, cmd.DateCreated, 1, cmd.Size, cmd.MD5Sum, cmd.Location, "??")
	if err != nil {
		fmt.Println("Error in query execution. Field types may differ and could not be cast\n", err)
	} else {
		fmt.Println("Added FileID row")
	}
}

// Shamelessly stolen
func jsonResponse(w http.ResponseWriter, err error, statusCode int) {
	var response struct {
		S string `json:"status"`
		M string `json:"message,omitempty"`
	}

	if err != nil {
		response.S = "error"
		response.M = err.Error()
	} else {
		response.S = "ok"
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print(err)
	}
}

func main() {

	err := dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		fmt.Println(err)
	}

	connectionTimeout := 1000 * time.Millisecond
	mux8700 := http.NewServeMux()
	mux8700.HandleFunc("/", filesEndpoint)
	server := &http.Server{
		Addr:         ":8700",
		Handler:      mux8700,
		ReadTimeout:  connectionTimeout,
		WriteTimeout: connectionTimeout,
	}

	log.Println("Started filesing port on 8700")
	server.ListenAndServe()
	log.Println("Finished server and closing")
}
