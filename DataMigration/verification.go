package main

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
)

func verifyObject(filename string) error {
	ctx := context.Background()
	// Connect to minio instance known as "Observatory" and gets stats of object "filename"
	temp, err := connectToObjectStorage("Observatory")
	objInfo, err := temp.minioClient.StatObject(ctx, "fyst", filename, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	// The next block queires the Files table for what the MD5 is for the given FileName to
	// compare the recorded value with the calculated one above
	var givenMD5sum string
	query := fmt.Sprintf(`select MD5sum from Files where FileName="%v";`, filename)
	queryReturn, err := dbCon.QueryRead(query, &struct{ MD5sum string }{})
	if err != nil {
		return fmt.Errorf(`Query for file %v returned nothing. File may not have been uploaded`, filename)
	} else if temp := queryReturn.Interface().([]struct{ MD5sum string }); len(temp) == 1 {
		givenMD5sum = temp[0].MD5sum // Above line makes sure the return query is of length one
	} else {
		return fmt.Errorf("Query for file %v returned more than one entry", filename)
	}

	// Check that MD5Sum given in initial file request matches the MD5Sum from the uploaded object
	if givenMD5sum != objInfo.ETag { // objInfo.ETag is MD5 calculated by Minio
		return fmt.Errorf(`verification failed. MD5sum hashes were different. Got %v but calculated %v`, givenMD5sum, objInfo.ETag)
	}
	log.Printf("MD5Sum verification successful for %v\n", filename)
	return nil

}
