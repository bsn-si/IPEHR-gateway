package base

import (
	"github.com/google/uuid"
	"hms/gateway/pkg/errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// ObjectVersionID
// Globally unique identifier for one version of a versioned object; lexical form: object_id '::' creating_system_id '::' version_tree_id.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_object_version_id_class
type ObjectVersionID struct {
	UID *UIDBasedID `json:"uid,omitempty"`

	objectID         uuid.UUID
	creatingSystemID string
	versionTreeID    string
}

const (
	startTrunkNumber = "1"
	uidDelimiter     = "::"
)

type UID interface {
	String() string
	ObjectID() uuid.UUID
	CreatingSystemID() string
	VersionTreeID() string
	//IsBranch() bool
}

func (o *ObjectVersionID) New(UID string, creatingSystemID string) {
	o.creatingSystemID = creatingSystemID
	o.parseUID(UID)
	o.UID = &UIDBasedID{ObjectID{Value: o.String()}}
}

func (o *ObjectVersionID) String() string {
	uid := []string{o.ObjectID().String(), o.CreatingSystemID(), o.VersionTreeID()}
	return strings.Join(uid, uidDelimiter)
}

func (o *ObjectVersionID) BasedID() string {
	UID := []string{o.ObjectID().String(), o.CreatingSystemID()}
	return strings.Join(UID, uidDelimiter)
}

func (o *ObjectVersionID) ObjectID() uuid.UUID {
	return o.objectID
}

func (o *ObjectVersionID) CreatingSystemID() string {
	return o.creatingSystemID
}

func (o *ObjectVersionID) VersionTreeID() string {
	return o.versionTreeID
}

func (o *ObjectVersionID) Equal(ver string) bool {
	if o.VersionTreeID() == ver {
		return true
	}
	return false
}

func (o *ObjectVersionID) setVersionTreeID(ver string) {
	if ver == "" {
		ver = startTrunkNumber
	}

	o.versionTreeID = ver
}

func (o *ObjectVersionID) parseUID(UID string) {
	re := regexp.MustCompile(uidDelimiter)
	parts := re.Split(UID, -1)

	if length := len(parts); length == 0 {
		return
	} else if length == 1 {
		o.setVersionTreeID("")
	} else if length == 2 {
		ver := strings.Join(parts[1:2], "")
		if !o.isVersion(ver) {
			o.creatingSystemID = ver
			ver = ""
		}
		o.setVersionTreeID(ver)
	} else if length == 3 {
		o.creatingSystemID = strings.Join(parts[1:2], "")
		o.setVersionTreeID(strings.Join(parts[2:3], ""))
	}

	objectID := strings.Join(parts[0:1], "")

	objectUUID, err := uuid.Parse(objectID)
	if err != nil {
		log.Fatal(errors.ErrIncorrectFormat)
	}

	o.objectID = objectUUID
}

func (o *ObjectVersionID) isVersion(ver string) bool {
	re := regexp.MustCompile(`^(\d+)+$`)
	return re.MatchString(ver)
}

func (o *ObjectVersionID) IncreaseUIDVersion() (ver string, err error) {
	ver = o.VersionTreeID()
	if ver == "" {
		err := errors.ErrObjectNotInit

		return "", err
	}

	verInt, err := strconv.Atoi(ver)
	if err != nil {
		return "", err
	}
	verInt++

	ver = strconv.Itoa(verInt)
	o.setVersionTreeID(ver)

	return
}
