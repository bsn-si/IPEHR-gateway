package base

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

//	Version tree identifier for one version. Lexical form: trunk_version [ '.' branch_number '.' branch_version ]
//
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_version_tree_id_class
type VersionTreeID struct {
	trunkVersion  string
	branchNumber  string
	branchVersion string
}

const delimiter = "."

var isVersion = regexp.MustCompile(`^(\d{1,8})+(\.\d{1,8}){0,2}$`)

func NewVersionTreeID(version string) (*VersionTreeID, error) {
	if !isVersion.MatchString(version) {
		return nil, errors.ErrIncorrectFormat
	}

	parts := strings.Split(version, delimiter)

	switch len(parts) {
	case 2:
		return &VersionTreeID{
			trunkVersion: parts[0],
			branchNumber: parts[1],
		}, nil
	case 3:
		return &VersionTreeID{
			trunkVersion:  parts[0],
			branchNumber:  parts[1],
			branchVersion: parts[2],
		}, nil
	default:
		return &VersionTreeID{
			trunkVersion: parts[0],
		}, nil
	}
}

func (v *VersionTreeID) String() string {
	result := []string{v.trunkVersion}

	if v.branchNumber != "" {
		result = append(result, v.branchNumber)
	}

	if v.branchVersion != "" {
		result = append(result, v.branchVersion)
	}

	return strings.Join(result, delimiter)
}

func (v *VersionTreeID) Equal(ver string) bool {
	return v.String() == ver
}

func (v *VersionTreeID) Increase() string {
	incr := func(s *string) {
		i, err := strconv.Atoi(*s)
		if err != nil {
			panic(err)
		}

		*s = strconv.Itoa(i + 1)
	}

	switch {
	case v.branchVersion != "":
		incr(&v.branchVersion)
	case v.branchNumber != "":
		incr(&v.branchNumber)
	default:
		incr(&v.trunkVersion)
	}

	return v.String()
}
