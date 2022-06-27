package status

type DocumentStatus uint8

const (
	ACTIVE  DocumentStatus = 1
	DELETED DocumentStatus = 2
)

var statusNames map[DocumentStatus]string = map[DocumentStatus]string{
	ACTIVE:  "ACTIVE",
	DELETED: "DELETED",
}

func (t DocumentStatus) String() string {
	return statusNames[t]
}
