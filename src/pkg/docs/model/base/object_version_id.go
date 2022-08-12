package base

import (
	"fmt"
	"hms/gateway/pkg/errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ObjectVersionID
// Globally unique identifier for one version of a versioned object; lexical form: object_id '::' creating_system_id '::' version_tree_id.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_object_version_id_class
type ObjectVersionID struct {
	UID *UIDBasedID `json:"uid,omitempty"`

	objectID         uuid.UUID
	creatingSystemID EhrSystemID
	versionTreeID    string
	versionBytes     *[32]byte
}

const (
	startTrunkNumber = "1"
	uidDelimiter     = "::"
)

type UID interface {
	String() string
	ObjectID() uuid.UUID
	CreatingSystemID() EhrSystemID
	VersionTreeID() string
	//IsBranch() bool
}

func NewObjectVersionID(UID string, creatingSystemID EhrSystemID) (*ObjectVersionID, error) {
	o := &ObjectVersionID{
		creatingSystemID: creatingSystemID,
	}

	if err := o.parseUID(UID); err != nil {
		return nil, err
	}

	o.UID = &UIDBasedID{ObjectID{Value: o.String()}}

	return o, nil
}

func (o *ObjectVersionID) String() string {
	uid := []string{o.ObjectID().String(), o.CreatingSystemID().String(), o.VersionString()}
	return strings.Join(uid, uidDelimiter)
}

func (o *ObjectVersionID) BasedID() string {
	UID := []string{o.ObjectID().String(), o.CreatingSystemID().String()}
	return strings.Join(UID, uidDelimiter)
}

func (o *ObjectVersionID) ObjectID() uuid.UUID {
	return o.objectID
}

func (o *ObjectVersionID) CreatingSystemID() EhrSystemID {
	return o.creatingSystemID
}

func (o *ObjectVersionID) VersionString() string {
	return o.versionTreeID
}

func (o *ObjectVersionID) VersionBytes() *[32]byte {
	return o.versionBytes
}

func (o *ObjectVersionID) Equal(ver string) bool {
	return o.VersionString() == ver
}

func (o *ObjectVersionID) setVersionTreeID(ver string) {
	if ver == "" {
		ver = startTrunkNumber
	}

	o.versionTreeID = ver

	o.versionBytes = &[32]byte{}
	copy(o.versionBytes[:], []byte(o.versionTreeID))
}

func (o *ObjectVersionID) parseUID(UID string) (err error) {
	re := regexp.MustCompile(uidDelimiter)
	parts := re.Split(UID, -1)

	if length := len(parts); length == 0 {
		return
	} else if length == 1 {
		o.setVersionTreeID("")
	} else if length == 2 {
		ver := strings.Join(parts[1:2], "")
		if !o.isVersion(ver) {
			if !o.creatingSystemID.Equal(ver) {
				return errors.ErrIncorrectFormat
			}

			ver = ""
		}
		o.setVersionTreeID(ver)
	} else if length == 3 {
		creatingSystemID := strings.Join(parts[1:2], "")

		if !o.creatingSystemID.Equal(creatingSystemID) {
			return errors.ErrIncorrectFormat
		}

		o.setVersionTreeID(strings.Join(parts[2:3], ""))
	}

	objectID := strings.Join(parts[0:1], "")

	objectUUID, err := uuid.Parse(objectID)
	if err != nil {
		return errors.ErrIncorrectFormat
	}

	o.objectID = objectUUID

	return nil
}

func (o *ObjectVersionID) isVersion(ver string) bool {
	re := regexp.MustCompile(`^(\d+)+(\.\d+)?$`)
	return re.MatchString(ver)
}

func (o *ObjectVersionID) IncreaseUIDVersion() (string, error) {
	if o.VersionString() == "" {
		return "", errors.ErrObjectNotInit
	}

	parts := strings.Split(o.VersionString(), ".")
	last := len(parts) - 1

	verInt, err := strconv.Atoi(parts[last])
	if err != nil {
		return "", fmt.Errorf("IncreaseUIDVersion error: %w o.VersionTreeID %s", err, o.VersionString())
	}

	parts[last] = strconv.Itoa(verInt + 1)

	o.setVersionTreeID(strings.Join(parts, "."))

	return o.VersionString(), nil
}
