package base

import (
	"hms/gateway/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

// ObjectVersionId
// Globally unique identifier for one version of a versioned object; lexical form: object_id '::' creating_system_id '::' version_tree_id.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_object_version_id_class
type ObjectVersionId struct {
	UID *UIDBasedID `json:"uid,omitempty"`

	objectId         string
	creatingSystemId string
	versionTreeId    string
}

const (
	startTrunkNumber = "1"
	uidDelimiter     = "::"
)

type Uid interface {
	String() string
	ObjectId() string
	CreatingSystemId() string
	VersionTreeId() string
	//IsBranch() bool
}

func (o *ObjectVersionId) New(uid string, creatingSystemId string) {
	o.parseUid(uid)
	o.creatingSystemId = creatingSystemId
	o.UID = &UIDBasedID{ObjectID{Value: o.String()}}
}

func (o *ObjectVersionId) String() string {
	uid := []string{o.ObjectId(), o.CreatingSystemId(), o.VersionTreeId()}
	return strings.Join(uid, uidDelimiter)
}

func (o *ObjectVersionId) BasedId() string {
	uid := []string{o.ObjectId(), o.CreatingSystemId()}
	return strings.Join(uid, uidDelimiter)
}

func (o *ObjectVersionId) ObjectId() string {
	return o.objectId
}

func (o *ObjectVersionId) CreatingSystemId() string {
	return o.creatingSystemId
}

func (o *ObjectVersionId) VersionTreeId() string {
	return o.versionTreeId
}

func (o *ObjectVersionId) setVersionTreeId(ver string) {
	if ver == "" {
		ver = startTrunkNumber
	}
	o.versionTreeId = ver
}

func (c *ObjectVersionId) parseUid(uid string) {
	re := regexp.MustCompile(uidDelimiter)
	parts := re.Split(uid, -1)
	length := len(parts)
	if length == 0 {
		return
	} else if length == 1 {
		c.creatingSystemId = ""
		c.setVersionTreeId("")
	} else if length == 2 {
		// TODO спорный момент
		ver := strings.Join(parts[1:2], "")
		if c.isVersion(ver) != true {
			c.creatingSystemId = ver
			ver = ""
		}
		c.setVersionTreeId(ver)
	} else if length == 3 {
		c.creatingSystemId = strings.Join(parts[1:2], "")
		c.setVersionTreeId(strings.Join(parts[2:3], ""))
	}

	c.objectId = strings.Join(parts[0:1], "")
}

func (c *ObjectVersionId) isVersion(ver string) bool {
	re := regexp.MustCompile(`^(\d+.?)+$`)
	return re.MatchString(ver)
}

func (c *ObjectVersionId) IncreaseUidVersion() (err error, ver string) {
	ver = c.VersionTreeId()
	if ver == "" {
		err := errors.ErrObjectNotInit
		return err, ""
	}

	verInt, err := strconv.Atoi(ver)
	if err != nil {
		return err, ""
	}
	verInt++

	ver = strconv.Itoa(verInt)
	c.setVersionTreeId(ver)
	return nil, ver
}
