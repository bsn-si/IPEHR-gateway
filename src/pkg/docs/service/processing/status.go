package processing

type Status uint8

const (
	StatusFailed     Status = 0
	StatusSuccess    Status = 1
	StatusPending    Status = 2
	StatusProcessing Status = 3
	StatusUnknown    Status = 255
)

var statuses = map[Status]string{
	StatusFailed:     "Failed",
	StatusSuccess:    "Success",
	StatusPending:    "Pending",
	StatusProcessing: "Processing",
	StatusUnknown:    "Unknown",
}

func (s Status) String() string {
	if status, ok := statuses[s]; ok {
		return status
	}

	return statuses[StatusUnknown]
}
