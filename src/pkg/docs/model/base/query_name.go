package base

import (
	"hms/gateway/pkg/errors"
	"regexp"
	"strings"
)

// QueryName
// The (fully qualified) name of the query (when is registered as a stored query), in a format of [{namespace}::]{query-name}. The namespace prefix is optional, and when used it should be in a form of a reverse domain name.
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/StoredQuery_schema
type QueryName struct {
	namespace string
	name      string
}

const (
	delimiter = "::"
)

func NewQueryName(s string) (*QueryName, error) {
	o := &QueryName{}

	if s == "" {
		return nil, errors.ErrIsEmpty
	}

	re := regexp.MustCompile(delimiter)
	parts := re.Split(s, -1)

	o.name = s

	if length := len(parts); length == 2 {
		o.namespace = strings.Join(parts[:1], "")
		o.name = strings.Join(parts[1:], "")
	}

	return o, nil
}

type QueryNameI interface {
	String() string
	Validate() func() error
}

func (o *QueryName) String() string {
	if o.namespace == "" {
		return o.name
	}

	name := []string{o.namespace, o.name}
	return strings.Join(name, delimiter)
}
