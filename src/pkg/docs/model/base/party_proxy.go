package base

// PartyProxy
// Abstract concept of a proxy description of a party, including an optional link to data for this party
// in a demographic or other identity management system.
//
// https://specifications.openehr.org/releases/RM/latest/common.html#_party_proxy_class
type PartyProxy struct {
	ExternalRef ObjectRef `json:"external_ref"`
}
