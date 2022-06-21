package base

// UidBasedId
// Abstract model of UID-based identifiers consisting of a root part and an optional extension;
// lexical form: root '::' extension
// https://specifications.openehr.org/releases/RM/Release-1.0.2/support.html#_uid_based_id_class
type UidBasedId struct {
	ObjectId
}
