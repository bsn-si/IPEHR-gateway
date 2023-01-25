package base

import (
	"encoding/json"
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type PartyProxyBase struct {
	Type        ItemType   `json:"_type"`
	ExternalRef *ObjectRef `json:"external_ref,omitempty"`
}

type PartyProxier interface {
	GetExternalRef() *ObjectRef
}

// PartyProxy
// Abstract concept of a proxy description of a party, including an optional link to data for this party
// in a demographic or other identity management system.
//
// https://specifications.openehr.org/releases/RM/latest/common.html#_party_proxy_class
type PartyProxy struct {
	data PartyProxier
}

func NewPartyProxy(pp PartyProxier) PartyProxy {
	return PartyProxy{
		data: pp,
	}
}

func (pp PartyProxy) MarshalJSON() ([]byte, error) {
	return json.Marshal(pp.data)
}

func (pp *PartyProxy) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type ItemType `json:"_type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "cannot unmarshal PartyProxy tmp obj")
	}

	switch tmp.Type {
	case PartySelfItemType:
		pp.data = &PartySelf{}
	case PartyIdentifiedItemType:
		pp.data = &PartyIdentified{}
	case PartyRelatedItemType:
		pp.data = &PartyIdentified{}
	default:
		return fmt.Errorf("unexpected PartyProxy type: '%v'", tmp.Type) //nolint
	}

	if err := json.Unmarshal(data, pp.data); err != nil {
		return errors.Wrap(err, "cannot unmarshal PartyProxy.data instance")
	}

	return nil
}

// Party proxy representing the subject of the record.
// Used to indicate that the party is the owner of the record.
// May or may not have external_ref set.
// https://specifications.openehr.org/releases/RM/latest/common.html#_party_self_class
type PartySelf struct {
	PartyProxyBase
}

func (ps PartySelf) GetExternalRef() *ObjectRef {
	return ps.ExternalRef
}

// PartyIdentified
// Proxy data for an identified party other than the subject of the record, minimally consisting of
// human-readable identifier(s), such as name, formal (and possibly computable) identifiers such as NHS
// number, and an optional link to external data.
// https://specifications.openehr.org/releases/RM/latest/common.html#_party_identified_class
type PartyIdentified struct {
	Name        string         `json:"name"`
	Identifiers []DvIdentifier `json:"identifiers"`
	PartyProxyBase
}

func (pi PartyIdentified) GetExternalRef() *ObjectRef {
	return pi.ExternalRef
}

// https://specifications.openehr.org/releases/RM/latest/common.html#_party_related_class
type PartyRelated struct {
	Name         string      `json:"name"`
	Relationship DvCodedText `json:"relationship"`
	PartyProxyBase
}

func (pi PartyRelated) GetExternalRef() *ObjectRef {
	return pi.ExternalRef
}

//TODO: add realiaston for PartyProxy related objects https://specifications.openehr.org/releases/RM/latest/common.html#_class_descriptions
