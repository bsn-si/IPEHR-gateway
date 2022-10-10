package roles

type RoleType uint8

const (
	Patient RoleType = iota
	Doctor
)

var typeNames = map[RoleType]string{
	Patient: "PATIENT",
	Doctor:  "DOCTOR",
}

func (t RoleType) String() string {
	return typeNames[t]
}
