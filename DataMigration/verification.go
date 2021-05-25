package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

// func (tData *TransferData) verifyObject(filename string) error {
func verifyObject(filename string) error {
	var realMD5sum string
	var givenMD5sum string

	// location := "us-east-1"              // TODO: This should be called from the table
	tempFile := "verificationObject" // temp file for temp solution

	temp, err := connectToObjectStorage("Observatory")
	if err = temp.minioClient.FGetObject(temp.ctx, "fyst", filename,
		tempFile, minio.GetObjectOptions{}); err != nil {
		return err
	} else if realMD5sum, err = hash_file_md5(tempFile); err != nil {
		return err
	}

	// Remove temporary file now that we've calculated the MD5Sum of it
	if err := os.Remove(tempFile); err != nil {
		return err
	}

	query := fmt.Sprintf(`select MD5sum from Files where FileName="%v";`, filename)
	queryReturn, err := dbCon.QueryRead(query, &struct{ MD5sum string }{})
	if err != nil {
		return fmt.Errorf(`Query for file %v returned nothing. File may not have been uploaded`, filename)
	} else if temp := queryReturn.Interface().([]struct{ MD5sum string }); len(temp) == 1 {
		givenMD5sum = temp[0].MD5sum // Above line makes sure the return query is of length one
	} else {
		return fmt.Errorf("Query for file %v returned nothing", filename)
	}

	// Check that MD5Sum given in initial file request matches the MD5Sum from the uploaded object
	if givenMD5sum == realMD5sum {
		log.Printf("MD5Sum verification successful for %v\n", filename)
		return nil // MD5's match, all clear
	} else {
		return fmt.Errorf("verification failed. MD5sum hashes were different. Got %v but calculated %v", givenMD5sum, realMD5sum)
	}

}

func hash_file_md5(filePath string) (string, error) {
	// Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	// Open the passed argument and check for any error
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return returnMD5String, err
	}

	hash := md5.New() // Open a new hash interface to write to
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	// Get the 16 bytes hash and convert the bytes to a string
	returnMD5String = hex.EncodeToString(hash.Sum(nil)[:16]) //
	return returnMD5String, nil
}
