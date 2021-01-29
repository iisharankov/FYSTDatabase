package overheadsql

// BackupLocationTable is an SQL table
type BackupLocationTable struct {
	LocationID      int    `json:"location_id"`
	LocationName    string `json:"location_name"`
	LocationAddress string `json:"location_address"`
	FTPIPAddress    string `json:"ftp_ip_address"`
}

// InstrumentTable is an SQL table
type InstrumentTable struct {
	InstrumentID   int    `json:"instrument_id"`
	InstrumentName string `json:"instrument_name"`
	FullName       string `json:"full_name"`
	Description    string `json:"description"`
	NumberOfPixels int    `json:"number_of_pixels"`
	FrequencyMin   int    `json:"frequency_min"`
	FrequencyMax   int    `json:"frequency_max"`
	TempRange      int    `json:"temp_range"`
}

// ObjectFileTable is an SQL table
type ObjectFileTable struct {
	FileID        int    `json:"file_id"`
	DateCreated   string `json:"date_created"` // dataetime?
	InstrumentID  int    `json:"instrument_id"`
	Size          int    `json:"size"`
	HashOfBytes   string `json:"hash_of_bytes"`
	ObjectStorage string `json:"object_storage"`
}

// RuleTable is an SQL table
type RuleTable struct {
	RuleID          int    `json:"rule_id"`
	RuleDescription string `json:"rule_description"`
	InstrumentID    int    `json:"instrument_id"`
	LocationID      int    `json:"location_id"`
	Active          int    `json:"active"` // tinyint or bool
}

// LogTable is an SQL table
type LogTable struct {
	FileID     int    `json:"file_id"`
	RuleID     int    `json:"rule_id"`
	BackupDate string `json:"backup_date"` // datetime?
	IsCopying  int    `json:"is_copying"`  // tiny int or bool
	URL        string `json:"url"`
}
