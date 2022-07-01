package types

type DocumentType uint8

const (
	Ehr         DocumentType = 1
	EhrAccess   DocumentType = 2
	EhrStatus   DocumentType = 3
	Composition DocumentType = 4
)

var typeNames = map[DocumentType]string{
	Ehr:         "EHR",
	EhrAccess:   "EHR_ACCESS",
	EhrStatus:   "EHR_STATUS",
	Composition: "COMPOSITION",
}

func (t DocumentType) String() string {
	return typeNames[t]
}
