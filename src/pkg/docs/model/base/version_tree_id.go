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
	for i := 0; i < len(parts) && i < 3; i++ {
		switch i {
		case 0:
			v.trunkVersion = parts[i]
		case 1:
			v.branchNumber = parts[i]
		case 2:
			v.branchVersion = parts[i]
		}
	}

	return nil
}

func (v *VersionTreeID) String() string {
	result := []string{}

	if v.trunkVersion != "" {
		result = append(result, v.trunkVersion)
	}

	if v.branchNumber != "" {
		result = append(result, v.branchNumber)
	}

	if v.branchVersion != "" {
		result = append(result, v.branchVersion)
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
