package processing

type Status uint8

const (
	StatusFailed     Status = 0
	StatusSuccess    Status = 1
	StatusNotFound   Status = 2
	StatusPending    Status = 3
	StatusProcessing Status = 4
	StatusUnknown    Status = 255
)

var statuses = map[Status]string{
	StatusFailed:     "Failed",
	StatusSuccess:    "Success",
	StatusNotFound:   "NotFound",
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
