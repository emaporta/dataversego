package dataversego

type RetrieveSignature struct {
	auth          Authorization
	tableName     string
	id            string
	columns       []string
	columnsString string
	printerror    bool
}

type RetrieveMultipleSignature struct {
	auth          Authorization
	tableName     string
	columns       []string
	columnsString string
	filter        Filter
	filterString  string
	printerror    bool
}
