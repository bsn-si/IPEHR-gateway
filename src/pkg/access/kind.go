package access

import "strings"

type Kind = uint8

const (
	NoKind Kind = iota
	Doc
	DocGroup
	UserGroup
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
		return NoKind
	}
}
