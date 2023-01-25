package model

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/users"
)

type Attribute = uint8

const (
	AttributeID              Attribute = 1
	AttributeIDEncr          Attribute = 2
	AttributeKeyEncr         Attribute = 3
	AttributeDocUIDHash      Attribute = 4
	AttributeDocUIDEncr      Attribute = 5
	AttributeDealCid         Attribute = 6
	AttributeMinerAddress    Attribute = 7
	AttributeContent         Attribute = 8
	AttributeContentEncr     Attribute = 9
	AttributeDescriptionEncr Attribute = 10
	AttributePasswordHash    Attribute = 11
	AttributeTimestamp       Attribute = 12
	AttributeNameEncr        Attribute = 13
)

type (
	AttributesEhr   []ehrIndexer.AttributesAttribute
	AttributesUsers []users.AttributesAttribute
)

func (a AttributesEhr) GetByCode(code Attribute) []byte {
	for _, attr := range a {
		if attr.Code == code {
			return attr.Value
		}
	}
	return nil
}

func (a AttributesUsers) GetByCode(code Attribute) []byte {
	for _, attr := range a {
		if attr.Code == code {
			return attr.Value
		}
	}
	return nil
}
