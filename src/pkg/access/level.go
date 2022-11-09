package access

import "strings"

type Level = uint8

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

func LevelFromString(str string) Level {
	switch strings.ToLower(str) {
	case "owner":
		return Owner
	case "admin":
		return Admin
	case "read":
		return Read
	default:
		return Unknown
	}
}

func LevelToString(l Level) string {
	return levelNames[l]
}
