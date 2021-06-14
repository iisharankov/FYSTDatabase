package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/iisharankov/FYSTDatabase/datasets"
)

// DBAPI manages communication with the DBAPI.
type DBAPI struct {
	Host   string
	client *http.Client
}

// NewDBAPI returns a new connection to host.
func NewDBAPI(host string) *DBAPI {
	return &DBAPI{
		Host: host,
		client: &http.Client{
			Timeout: 10000 * time.Millisecond,
		},
	}
}

func (dbapi *DBAPI) do(req *http.Request) ([]byte, error) {
	// Send request req to server
	resp, err := dbapi.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the body of the http response into b
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	trimmedString := strings.Trim(string(b), "\n") // Trims extra newline

	// Send request req to server
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(trimmedString)
	}

	// // unpack b into the struct it was encoded as
	// var response datasets.ServerHTTPResponse
	// if err = json.Unmarshal(b, &response); err != nil {
	// 	log.Println("REFREFREFERFs")
	// 	log.Println(err)
	// }
	//
	// log.Println(response)
	// // Server responds with Status S of "ok" or "error". Respond accordingly
	// if response.S == "error" {
	// 	return nil, fmt.Errorf("Error from Server: " + response.M)
	// }
	// return []byte(response.M), nil
	return []byte(trimmedString), nil
}

func (dbapi *DBAPI) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	req.Host = dbapi.Host
	req.URL.Host = dbapi.Host
	req.URL.Scheme = "http"
	return req, err
}

func (dbapi *DBAPI) get(path string) ([]byte, error) {
	req, err := dbapi.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return dbapi.do(req)
}

func (dbapi *DBAPI) post(path, contentType string, body io.Reader) ([]byte, error) {
	req, err := dbapi.newRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	return dbapi.do(req)
}

// Patch command depreciated
// func (dbapi *DBAPI) patch(path, contentType string, body io.Reader) ([]byte, error) {
// 	req, err := dbapi.newRequest("PATCH", path, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", contentType)
// 	return dbapi.do(req)
// }

// * // * // client specific functions below // * // * //
func (dbapi *DBAPI) getFiles(filename string) ([]byte, error) {
	var endpoint string = "/files"
	if filename != "" { // If a specific ID was given, add ID to endpoint
		endpoint += "/" + filename
	}

	return dbapi.get(endpoint)
}

func (dbapi *DBAPI) getRules(id string) ([]byte, error) {
	var endpoint string = "/rules"
	if id != "" { // If a specific ID was given, add ID to endpoint
		endpoint += "/" + id
	}

	a, err := dbapi.get(endpoint)
	return a, err
}

func (dbapi *DBAPI) requestToUploadFile(file datasets.File) ([]byte, error) {
	var body bytes.Buffer
	dec := json.NewEncoder(&body)
	if err := dec.Encode(file); err != nil {
		return nil, err
	}

	switch {
	case file.Name == "":
		return nil, fmt.Errorf("Name cannot be empty. Canceled.")
	case file.MD5Sum == "":
		return nil, fmt.Errorf("MD5sum cannot be empty. Canceled.")
	default:
		return dbapi.post("/files", "application/x-www-form-urlencoded", &body)
	}
}

func (dbapi *DBAPI) requestToUpdateLog(name string, reply datasets.ClientUploadReply) ([]byte, error) {
	var body bytes.Buffer

	endpoint := fmt.Sprintf("/records/%v/%v", reply.FileName, reply.UploadLocation)
	a, err := dbapi.post(endpoint, "application/x-www-form-urlencoded", &body)
	return a, err
}

func (dbapi *DBAPI) requestToUpdateCopies(filename string, LocationID int) ([]byte, error) {
	var body bytes.Buffer
	dec := json.NewEncoder(&body)
	err := dec.Encode(LocationID)

	// TODO: Still pass locationID AND FileID to endpoint in body to populate Copies table, normalization!
	endpoint := fmt.Sprintf("/files/%v/copies", filename)
	a, err := dbapi.post(endpoint, "application/x-www-form-urlencoded", &body)
	return a, err
}

// func (dbapi *DBAPI) requestToUpdateObjectFile(file datasets.File, reply datasets.ClientUploadReply) ([]byte, error) {
// 	var body bytes.Buffer
// 	dec := json.NewEncoder(&body)
// 	err := dec.Encode(true)
// 	log.Println("1")
// 	endpoint := fmt.Sprintf("/files/%d", reply.FileID)
// 	a, err := dbapi.patch(endpoint, "application/x-www-form-urlencoded", &body)
// 	return a, err
// }

// func (dbapi *DBAPI) requestToUpdateLog(reply datasets.ClientUploadReply) ([]byte, error) {
// 	var body bytes.Buffer
// 	dec := json.NewEncoder(&body)
// 	err := dec.Encode(reply.UploadLocation)

// 	endpoint := fmt.Sprintf("/files/%d", reply.FileID)
// 	// endpoint = "/records"
// 	log.Println(endpoint)
// 	// a, err := dbapi.post(endpoint, "application/x-www-form-urlencoded", &body)
// 	a, err := dbapi.post(endpoint, "application/x-www-form-urlencoded", &body)
// 	return a, err
// }

func (dbapi *DBAPI) logGET() ([]byte, error) {
	var endpoint string = "/records"
	a, err := dbapi.get(endpoint)
	return a, err
}
