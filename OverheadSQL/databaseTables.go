package overheadsql

type backupLocationTable struct {
	LocationID      int
	LocationName    string
	LocationAddress string
	FTPIPAddress    string
}

type instrumentTable struct {
	InstrumentID   int
	InstrumentName string
	FullName       string
	Description    string
	NumberOfPixels int
	FrequencyMin   int
	FrequencyMax   int
	TempRange      int
}

type objectFileTable struct {
	FileID         int
	DateCreated    string // dataetime?
	InstrumentID   int
	Size           int
	HashOfBytes    string
	LocationOnDisk string
	ObjectStorage  string
}

type ruleTable struct {
	RuleID          int
	RuleDescription int
	InstrumentID    int
	LocationIDint   int
	Active          int // tinyint or bool
}

type logTable struct {
	FileID     int
	RuleID     int
	BackupDate string // datetime?
	IsCopying  int    // tiny int or bool
}
