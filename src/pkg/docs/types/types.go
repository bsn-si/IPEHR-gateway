package types

type DocumentType uint8

const (
	Ehr DocumentType = iota
	EhrAccess
	EhrStatus
	Composition
	Query
	Template
	Directory
)

var typeNames = map[DocumentType]string{
	Ehr:         "EHR",
	EhrAccess:   "EHR_ACCESS",
	EhrStatus:   "EHR_STATUS",
	Composition: "COMPOSITION",
	Query:       "QUERY",
	Template:    "TEMPLATE",
	Directory:   "DIRECTORY",
}

func (t DocumentType) String() string {
	return typeNames[t]
}
