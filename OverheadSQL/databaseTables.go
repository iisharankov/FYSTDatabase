package overheadsql

type backupLocation struct {
	LocationID      int
	LocationName    string
	LocationAddress string
	FTPIPAddress    string
}

// InstrumentTable is a SQL table changed to a struct
// type InstrumentTable struct {
// 	InstrumentID   int
// 	InstrumentName string
// 	FullName       string
// 	Description    string
// 	NumberOfPixels int
// 	FrequencyMin   int
// 	FrequencyMax   int
// 	TempRange      int
// }
