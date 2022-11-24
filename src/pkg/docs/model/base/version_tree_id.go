package base

import (
	"fmt"
	"hms/gateway/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

//	Version tree identifier for one version. Lexical form: trunk_version [ '.' branch_number '.' branch_version ]
//
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_version_tree_id_class
type VersionTreeID struct {
	trunkVersion  string
	branchNumber  string
	branchVersion string
	delimiter     string
}

func NewVersionTreeID(version string) (*VersionTreeID, error) {
	v := &VersionTreeID{}

	if len(version) == 0 {
		return nil, errors.ErrIsEmpty
	}

	if ok := v.isVersion(version); !ok {
		return nil, errors.ErrIncorrectFormat
	}

	v.delimiter = "."
	if err := v.parse(v.split(version)); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	return v, nil
}

func (v *VersionTreeID) isVersion(ver string) bool {
	re := regexp.MustCompile(`^(\d+)+(\.\d+)*$`)
	return re.MatchString(ver)
}

func (v *VersionTreeID) split(ver string) []string {
	return strings.Split(ver, v.delimiter)
}

func (v *VersionTreeID) parse(parts []string) error {
	if length := len(parts); length == 1 {
		v.trunkVersion = strings.Join(parts[:1], "")
	} else if length == 2 {
		v.trunkVersion = strings.Join(parts[:1], "")
		v.branchNumber = strings.Join(parts[1:2], "")
	} else if length == 3 {
		v.trunkVersion = strings.Join(parts[:1], "")
		v.branchNumber = strings.Join(parts[1:2], "")
		v.branchVersion = strings.Join(parts[2:3], "")
	}

	return nil
}

func (v *VersionTreeID) String() string {
	ver := [3]string{v.branchVersion, v.branchNumber, v.trunkVersion}
	result := []string{}

	for _, p := range ver {
		if p == "" && len(result) == 0 {
			continue
		}

		result = append([]string{p}, result...)
	}

	return strings.Join(result, v.delimiter)
}

func (v *VersionTreeID) Equal(ver string) bool {
	return v.String() == ver
}

func (v *VersionTreeID) Increase() (string, error) {
	parts := v.split(v.String())
	last := len(parts) - 1

	verInt, err := strconv.Atoi(parts[last])
	if err != nil {
		return "", fmt.Errorf("Increase error: %w, ver is %s", err, v.String())
	}

	parts[last] = strconv.Itoa(verInt + 1)

	if err := v.parse(parts); err != nil {
		return "", fmt.Errorf("Increase error: %w, ver is %s", err, v.String())
	}

	return v.String(), nil
}
