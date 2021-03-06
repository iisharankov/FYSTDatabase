package datasets

import "time"

// File is the struct type the client sends to the server. Both require it
type File struct {
	Name        string    `json:"name"`
	Instrument  int       `json:"instrument"`
	MD5Sum      string    `json:"md5sum"`
	DateCreated time.Time `json:"date_created"`
	Size        int       `json:"size"`
	URL         string    `json:"url"`
}

// FilesThatNeedToBeBackedUp lists all the data required to move file from FYST to external location
type FilesThatNeedToBeBackedUp struct {
	FileID         int    `json:"file_d"`
	FileName       string `json:"file_name"`
	FileSize       int    `json:"file_size"`
	DateCreated    string `json:"date"`
	MD5sum         string `json:"md5"`
	RuleID         int    `json:"rule_id"`
	BucketName     string `json:"bucket_name"`
	UploadLocation string `json:"upload_location"`
}

// S3Metadata holds the connection information for a given minio instance
type S3Metadata struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key"`
	SecretAccessKey string `json:"secret_key"`
	UseSSL          string `json:"use_ssl"`
} // TODO: Change UseSSL to bool finally?

// ClientUploadReply is the JSON sent to POST requests with metadata for where to upload given file
type ClientUploadReply struct {
	S3Metadata     S3Metadata
	FileName       string `json:"file_name"`
	LocationID     int    `json:"location_id"`
	UploadLocation string `json:"upload_location"`
}

// Response given from FYST Server to client connections through HTTP
type ServerHTTPResponse struct {
	S string `json:"status"`
	M string `json:"message,omitempty"`
}
