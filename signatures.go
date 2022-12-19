package dataversego

type RetrieveSignature struct {
	Auth          Authorization
	TableName     string
	Id            string
	Columns       []string
	ColumnsString string
	Printerror    bool
}

type RetrieveMultipleSignature struct {
	Auth          Authorization
	TableName     string
	Columns       []string
	ColumnsString string
	Filter        Filter
	FilterString  string
	Printerror    bool
}
