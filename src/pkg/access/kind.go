package access

import "strings"

type Kind = uint8

const (
	Doc Kind = iota
	DocGroup
	UserGroup
	Unknown = 255
)

func KindFromString(str string) Kind {
	switch strings.ToLower(str) {
	case "doc":
		return Doc
	case "docgroup":
		return DocGroup
	case "usergroup":
		return UserGroup
	default:
		return Unknown
	}
}
