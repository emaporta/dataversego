package dataversego

// The 'RetrieveSignature' struct represents the signature of a 'Retrieve' function.
// It contains the following fields:
//   - Auth: a struct containing authentication information
//   - TableName: the name of the table to retrieve the entry from
//   - Id: the ID of the entry to be retrieved
//   - Columns: a slice of strings representing the columns to be retrieved
//   - ColumnsString: a string representing the columns to be retrieved
//   - Printerror: a boolean value indicating whether or not to print errors
type RetrieveSignature struct {
	Auth          Authorization
	TableName     string
	Id            string
	Columns       []string
	ColumnsString string
	Printerror    bool
}

// The 'RetrieveMultipleSignature' struct represents the signature of a 'RetrieveMultiple' function.
// It contains the following fields:
//   - Auth: a struct containing authentication information
//   - TableName: the name of the table to retrieve the entries from
//   - Columns: a slice of strings representing the columns to be retrieved
//   - ColumnsString: a string representing the columns to be retrieved
//   - Filter: a struct containing filter criteria for the entries to be retrieved
//   - FilterString: a string representing the filter criteria
//   - Printerror: a boolean value indicating whether or not to print errors
type RetrieveMultipleSignature struct {
	Auth          Authorization
	TableName     string
	Columns       []string
	ColumnsString string
	Filter        Filter
	FilterString  string
	Printerror    bool
}
