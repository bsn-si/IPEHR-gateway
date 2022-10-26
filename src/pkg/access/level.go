package access

type Level uint8

const (
	NoAccess Level = iota
	Owner
	Admin
	Read
)

var levelNames = map[Level]string{
	NoAccess: "NoAccess",
	Owner:    "Owner",
	Admin:    "Admin",
	Read:     "Read",
}

func (l Level) String() string {
	return levelNames[l]
}
