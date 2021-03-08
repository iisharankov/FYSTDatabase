package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// File has
type File struct {
	Name        string    `json:"name"`
	Instrument  int       `json:"instrument"`
	MD5Sum      string    `json:"md5sum"`
	DateCreated time.Time `json:"date_created"`
	Size        int       `json:"size"`
	URL         string    `json:"url"`
}

// ClientUploadReply is the JSON sent to POST requests with metadata for where to upload given file
type ClientUploadReply struct {
	FileID         int    `json:"file_id"`
	UploadLocation string `json:"upload_location"`
}

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

// TODO: Generalized this by parsing db table name and number of parameters to upload? might be hard with reflect
func addRowToObjectFile(newFile File, newFileID int, w http.ResponseWriter) bool {
	// Use SQL Prepare() method to safely convert the field types.
	stmt, err := dbCon.PrepareQuery("insert into ObjectFile values(?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Println("Error in db.Perpare()\n", err)
		jsonResponse(w, err, http.StatusBadRequest)
		return false
	}

	// Execute the command on the database (encoded already in stmt)
	_, err = stmt.Exec(newFileID, newFile.DateCreated, newFile.Instrument, newFile.Size, newFile.MD5Sum, newFile.URL)
	if err != nil {
		log.Println("Error in query execution\n", err)
		jsonResponse(w, err, http.StatusBadRequest)
		return false
	}

	log.Println("Added FileID row to ObjectFile Table")
	return true
}

// filesEndpoint listens to the http request header for a curl command to update the IP/Port of the filesing, or to enable/disable
func filesEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	endpoint := r.URL.Path
	endpointSections := strings.Split(endpoint, "/")

	// Check if connection to DB is possible, if not, fail.
	if err = dbCon.CheckConnection(); err != nil {
		var errstrings []string
		errstrings = append(errstrings, err.Error())

		// Concatonate error with additional one.
		errstrings = append(errstrings, fmt.Errorf("connection to database failed and server aborted").Error())

		// combine and return both errors for more useful debugging for user
		jsonResponse(w, fmt.Errorf(strings.Join(errstrings, " - ")), http.StatusServiceUnavailable)
		return
	}

	// Process POST or GET request
	if r.Method == "POST" {
		if endpointSections[1] == "files" {
			if len(endpointSections) == 2 {
				// This endpoint is when a file has been given to be added to a table within the database

				var newFile File
				err = dec.Decode(&newFile)
				log.Println("new file:", newFile)

				// Find largest index for FileID column, so we can increment by one
				temp := struct{ FileID int }{} // Temp struct that has just integer (SQL query returns row of 1 int)
				queryReturn, err := dbCon.QueryRead("SELECT FileID FROM ObjectFile ORDER BY FileID DESC LIMIT 1", &temp)
				if err != nil {
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				// Convert SQL output to golang int
				lastFileID, _ := queryReturn.Interface().([]struct{ FileID int }) // type assert reflect.Value to struct
				var newFileID int
				if len(lastFileID) == 0 {
					newFileID = 1 // If ObjectFile table empty, lastFileID is empty, so can't take FileID field
				} else {
					// Takes the first (only) row out, and calls FileID field to increment
					newFileID = lastFileID[0].FileID + 1
				}

				// Try to add given file to ObjectFile table with next FileID
				ok := addRowToObjectFile(newFile, newFileID, w)
				if ok { // If sucessful, respond to client with upload location details

					// Query BackupLocation Table to find bucket to upload file
					var rule BackupLocationTable
					query := `select * from BackupLocation b where b.LocationID = 1`
					queryReturn, err := dbCon.QueryRead(query, &rule)
					if err != nil {
						jsonResponse(w, err, http.StatusBadRequest)
						return
					}

					// convert reflect.Value to BackupLocationTable table
					returnQuery, _ := queryReturn.Interface().([]BackupLocationTable)

					// Create JSON with metadata necessary for the client
					replyData, err := json.Marshal(ClientUploadReply{
						FileID:         newFileID,
						UploadLocation: returnQuery[0].S3Bucket,
					})
					if err != nil {
						log.Fatal(err)
					}

					w.Write(replyData) // Upload JSON back to the client
					jsonResponse(w, err, http.StatusAccepted)
					return
				} else {
					var msg string = "addRowToObjectFile method failed, file not added to database"
					log.Println(msg)
					err := errors.New(msg)
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

			} else if len(endpointSections) == 3 {
				if endpointSections[2] != "" {
					// POST reuqest to add file to Logs table of database

					var copyLocation string
					err = dec.Decode(&copyLocation)
					log.Println("new log to add:", copyLocation)
					FileID, err := strconv.Atoi(endpointSections[2])

					// Tells you which RuleID corresponds to a FileID and Location string (locationname in BackupLocation)
					query := fmt.Sprintf(`select r.RuleID from Rule r 
					join ObjectFile o on o.InstrumentID=r.InstrumentID 
					join BackupLocation b on b.LocationID=r.LocationID
					where o.FileId = %v and b.S3Bucket = "%v"`, FileID, copyLocation)

					temp := struct{ RuleID int }{} // Temp struct that has just integer (SQL query returns row of 1 int)
					queryReturn, err := dbCon.QueryRead(query, &temp)
					if err != nil {
						jsonResponse(w, err, http.StatusBadRequest)
						return
					}
					returnQuery, _ := queryReturn.Interface().([]struct{ RuleID int }) // type assert from reflect.Value to struct

					if err := addRowToLog(FileID, returnQuery[0].RuleID); err != nil {
						jsonResponse(w, err, http.StatusBadRequest)
						return
					}

					jsonResponse(w, nil, http.StatusAccepted)
					return

				} else {
					err = errors.New("value after /files/ was empty or not given")
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}
			}
		} else {
			jsonResponse(w, err, http.StatusNotFound)
			return
		}

		// _, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }

	} else if r.Method == "GET" {
		if endpointSections[1] == "files" {
			if len(endpointSections) == 2 {
				// Client is asking for all files from database

				SQLQuery := "select * from ObjectFile;"
				var objectTable ObjectFileTable
				outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
				if err != nil {
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				// Convert reflect.Value() to ObjectFileTable, iterate over rows.
				outputData, _ := outputRows.Interface().([]ObjectFileTable)
				for _, val := range outputData {
					data, _ := json.Marshal(val)
					w.Write(data)
				}
				jsonResponse(w, err, http.StatusAccepted)
				return

			} else if len(endpointSections) == 3 && checkID(endpointSections[2], w) == nil {

				var SQLQuery string
				if strings.Contains(endpointSections[2], "-") {
					splitRange := strings.Split(endpointSections[2], "-")

					a, _ := strconv.Atoi(splitRange[0])
					b, _ := strconv.Atoi(splitRange[1])
					SQLQuery = "select * from ObjectFile ORDER BY FileID LIMIT " + strconv.Itoa(b-a+1) + " OFFSET " + strconv.Itoa(a-1)
				} else {
					// if int, use string version for simplicity. No worry about SQL injection since above Atoi didn't fail
					SQLQuery = "select * from ObjectFile where ObjectFile.FileID=" + endpointSections[2]
				}

				var objectTable ObjectFileTable
				outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
				if err != nil {
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				outputData, ok := outputRows.Interface().([]ObjectFileTable)
				if !ok || len(outputData) == 0 { // Len will be 0 if index is out of range (nothing is returned)
					errMsg := errors.New("Error with ID given, may be out of range of last element")
					jsonResponse(w, errMsg, http.StatusBadRequest)
					return
				}

				for _, val := range outputData {
					data, _ := json.Marshal(val)
					w.Write(data)
				}

				jsonResponse(w, err, http.StatusAccepted)
				return

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
				var objectTable LogTable
				outputRows, err := dbCon.QueryRead(SQLQuery, &objectTable)
				if err != nil {
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				outputData, _ := outputRows.Interface().([]LogTable)
				for _, val := range outputData {
					data, _ := json.Marshal(val)
					w.Write(data)
				}
				jsonResponse(w, err, http.StatusAccepted)
				return
			}
		} else if endpointSections[1] == "rules" {
			if len(endpointSections) == 2 {

				var ruleTable RuleTable
				outputRows, err := dbCon.QueryRead("select * from Rule;", &ruleTable)
				if err != nil {
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				outputData, _ := outputRows.Interface().([]RuleTable)
				for _, val := range outputData {
					data, _ := json.Marshal(val)
					w.Write(data)
				}
				jsonResponse(w, err, http.StatusAccepted)
				return

			} else if len(endpointSections) == 3 && checkID(endpointSections[2], w) == nil {
				// if int, use string version for simplicity. No worry about SQL injection since above Atoi didn't fail
				SQLQuery := "select * from Rule where Rule.RuleID=" + endpointSections[2]
				var ruleTable RuleTable
				outputRows, err := dbCon.QueryRead(SQLQuery, &ruleTable)
				if err != nil {
					jsonResponse(w, err, http.StatusBadRequest)
					return
				}

				outputData, ok := outputRows.Interface().([]RuleTable)
				if !ok || len(outputData) == 0 { // Len will be 0 if index is out of range (nothing is returned)
					errMsg := errors.New("Error with ID given, may be out of range of last element")
					jsonResponse(w, errMsg, http.StatusBadRequest)
					return
				}

				data, _ := json.Marshal(outputData[0])
				w.Write(data) // Should only be a single output
				jsonResponse(w, err, http.StatusAccepted)
				return

			}
		}
	} else {
		jsonResponse(w, err, http.StatusNotFound)
		return
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

func checkID(endpointSection string, w http.ResponseWriter) error {
	var err error
	// If endpointSection has a '-', it is a range and both elements should be tested recursively
	if strings.Contains(endpointSection, "-") {
		splitRange := strings.Split(endpointSection, "-")
		err = checkID(splitRange[0], w)
		if err != nil {
			return err
		}

		err = checkID(splitRange[1], w)
		if err != nil {
			return err
		}
		return nil
	}

	if endpointSection == "" {
		errMsg := errors.New("No file ID given. Either give file ID after '/'  or remove the '/'")
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return errMsg
	}

	// Check to see if input is an actual integer
	_, err = strconv.Atoi(endpointSection)
	if err != nil {
		errMsg := errors.New("Could not convert ID given to int, check value after 'files/'")
		jsonResponse(w, errMsg, http.StatusBadRequest)
		return errMsg
	}

	return nil
}

func startAPIServer() {
	// Containers fail if db connection is first
	// if err := dbCon.Connect(dbUsername, dbPassword, dbAddress, dbName); err != nil {
	// 	log.Println(err)
	// }

	connectionTimeout := 1000 * time.Millisecond
	mux8700 := http.NewServeMux()
	mux8700.HandleFunc("/", filesEndpoint)
	server := &http.Server{
		Addr:         ":8700",
		Handler:      mux8700,
		ReadTimeout:  connectionTimeout,
		WriteTimeout: connectionTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
	// go server.ListenAndServe()
	log.Println("Started filesing port on 8700.")
}
