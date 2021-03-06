package datasets

// LocationsTable is an SQL table
type LocationsTable struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	S3Bucket     string `json:"s3bucket"`
	Address      string `json:"address"`
	AccessID     string `json:"access_id"`
	SecretID     string `json:"secret_id"`
	SSL          bool   `json:"ssl"`
}

// InstrumentsTable is an SQL table
type InstrumentsTable struct {
	InstrumentID   int    `json:"instrument_id"`
	InstrumentName string `json:"instrument_name"`
	FullName       string `json:"full_name"` //
	Description    string `json:"description"`
	NumberOfPixels int    `json:"number_of_pixels"`
	FrequencyMin   int    `json:"frequency_min"`
	FrequencyMax   int    `json:"frequency_max"`
}

// FilesTable is an SQL table
type FilesTable struct {
	FileID       int    `json:"file_id"`
	FileName     string `json:"file_name"`
	DateCreated  string `json:"date_created"` // dataetime?
	InstrumentID int    `json:"instrument_id"`
	Size         int    `json:"size"`
	MD5sum       string `json:"md5_sum"`
}

// RulesTable is an SQL table
type RulesTable struct {
	RuleID          int    `json:"rule_id"`
	RuleDescription string `json:"rule_description"`
	InstrumentID    int    `json:"instrument_id"`
	LocationID      int    `json:"location_id"`
	Active          int    `json:"active"` // tinyint or boolgo get
}

// RecordsTable is an SQL table
type RecordsTable struct {
	FileID     int    `json:"file_id"`
	RuleID     int    `json:"rule_id"`
	BackupDate string `json:"backup_date"` // datetime?
}

// CopiesTable is an SQL table
type CopiesTable struct {
	FileID     int    `json:"file_id"`
	LocationID int    `json:"location_id"`
	URL        string `json:"url"`
}

var randomvar int
