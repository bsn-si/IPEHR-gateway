package status

type DocumentStatus uint8

const (
	ACTIVE  DocumentStatus = 0
	DELETED DocumentStatus = 1
)

var statusNames = map[DocumentStatus]string{
	ACTIVE:  "ACTIVE",
	DELETED: "DELETED",
}

func (t DocumentStatus) String() string {
	return statusNames[t]
}
