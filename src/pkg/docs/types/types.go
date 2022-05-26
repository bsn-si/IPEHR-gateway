package types

type DocumentType uint8

const (
	EHR         DocumentType = 1
	EHR_ACCESS  DocumentType = 2
	EHR_STATUS  DocumentType = 3
	COMPOSITION DocumentType = 4
)

var typeNames map[DocumentType]string = map[DocumentType]string{
	EHR:         "EHR",
	EHR_ACCESS:  "EHR_ACCESS",
	EHR_STATUS:  "EHR_STATUS",
	COMPOSITION: "COMPOSITION",
}

func (t DocumentType) String() string {
	return typeNames[t]
}
