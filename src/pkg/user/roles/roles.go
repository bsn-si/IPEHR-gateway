package roles

type Role uint8

const (
	Patient Role = iota
	Doctor
)

var typeNames = map[Role]string{
	Patient: "Patient",
	Doctor:  "Doctor",
}

func (t Role) String() string {
	return typeNames[t]
}
