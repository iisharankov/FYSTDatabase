package overheadsql

type BackupLocationTable struct {
	LocationID      int
	LocationName    string
	LocationAddress string
	FTPIPAddress    string
}

type InstrumentTable struct {
	InstrumentID   int
	InstrumentName string
	FullName       string
	Description    string
	NumberOfPixels int
	FrequencyMin   int
	FrequencyMax   int
	TempRange      int
}

type ObjectFileTable struct {
	FileID         int
	DateCreated    string // dataetime?
	InstrumentID   int
	Size           int
	HashOfBytes    string
	LocationOnDisk string
	ObjectStorage  string
}

type RuleTable struct {
	RuleID          int
	RuleDescription string
	InstrumentID    int
	LocationID      int
	Active          int // tinyint or bool
}

type LogTable struct {
	FileID     int
	RuleID     int
	BackupDate string // datetime?
	IsCopying  int    // tiny int or bool
}
