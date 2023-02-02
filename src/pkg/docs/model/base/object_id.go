package base

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

// ObjectID
// Ancestor class of identifiers of informational objects. Ids may be completely meaningless, in which
// case their only job is to refer to something, or may carry some information to do with the identified
// object.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_object_id_class
type ObjectID struct {
	Type  ItemType `json:"_type,omitempty"`
	Value string   `json:"value"`
}

// UIDBasedID
// Abstract model of UID-based identifiers consisting of a root part and an optional extension;
// lexical form: root '::' extension
// https://specifications.openehr.org/releases/RM/Release-1.0.2/support.html#_uid_based_id_class
type UIDBasedID struct {
	ObjectID
}

// HierObjectID
// Concrete type corresponding to hierarchical identifiers of the form defined by UID_BASED_ID.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_hier_object_id_class
type HierObjectID struct {
	UIDBasedID
}

// ObjectVersionID
// Globally unique identifier for one version of a versioned object; lexical form: object_id '::' creating_system_id '::' version_tree_id.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_object_version_id_class
type ObjectVersionID struct {
	UID *UIDBasedID `json:"uid,omitempty"`

	objectID         uuid.UUID
	creatingSystemID string
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
	CreatingSystemID() string
	VersionTreeID() string
	//IsBranch() bool
}

func NewObjectVersionID(UID string, creatingSystemID string) (*ObjectVersionID, error) {
	o := &ObjectVersionID{
		creatingSystemID: creatingSystemID,
	}

	if err := o.parseUID(UID); err != nil {
		return nil, fmt.Errorf("parseUID error: %w", err)
	}

	o.UID = &UIDBasedID{
		ObjectID{
			Value: o.String(),
		},
	}

	return o, nil
}

func (o *ObjectVersionID) String() string {
	return o.ObjectID().String() + uidDelimiter + o.CreatingSystemID() + uidDelimiter + o.VersionString()
}

func (o *ObjectVersionID) BasedID() string {
	return o.ObjectID().String() + uidDelimiter + o.CreatingSystemID()
}

func (o *ObjectVersionID) BaseIDHash() *[32]byte {
	h := sha3.Sum256([]byte(o.BasedID()))
	return &h
}

func (o *ObjectVersionID) ObjectID() uuid.UUID {
	return o.objectID
}

func (o *ObjectVersionID) CreatingSystemID() string {
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

func (o *ObjectVersionID) parseUID(UID string) error {
	parts := strings.Split(UID, uidDelimiter)

	switch len(parts) {
	case 0:
		return nil
	case 1:
		o.setVersionTreeID("")
	case 2:
		ver := parts[1]
		if !o.isVersion(ver) {
			// Using default systemID if it is empty
			if ver == "" {
				ver = common.EhrSystemID
			}

			if o.creatingSystemID != ver {
				return fmt.Errorf("%w creatingSystemID mismatch", errors.ErrIncorrectFormat)
			}

			ver = ""
		}

		o.setVersionTreeID(ver)
	case 3:
		creatingSystemID := parts[1]

		// Using default systemID if it is empty
		if creatingSystemID == "" {
			creatingSystemID = common.EhrSystemID
		}

		if o.creatingSystemID != creatingSystemID {
			return fmt.Errorf("%w creatingSystemID mismatch", errors.ErrIncorrectFormat)
		}

		o.setVersionTreeID(parts[2])
	}

	objectID := parts[0]

	objectUUID, err := uuid.Parse(objectID)
	if err != nil {
		return fmt.Errorf("uuid.Parse error: %w objectID: %s", err, objectID)
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

func (o *ObjectVersionID) BaseUIDHash() *[32]byte {
	UIDHash := sha3.Sum256([]byte(o.BasedID()))
	return &UIDHash
}
