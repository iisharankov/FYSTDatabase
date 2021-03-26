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

// S3Metadata holds the connection information for a given minio instance
type S3Metadata struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key"`
	SecretAccessKey string `json:"secret_key"`
	UseSSL          string `json:"use_ssl"`
}

// ClientUploadReply is the JSON sent to POST requests with metadata for where to upload given file
type ClientUploadReply struct {
	S3Metadata     S3Metadata
	FileID         int    `json:"file_id"`
	UploadLocation string `json:"upload_location"`
}
