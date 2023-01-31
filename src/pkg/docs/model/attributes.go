package model

type (
	Attribute = uint8

	Attributes []struct {
		Code  uint8
		Value []byte
	}
)

const (
	AttributeID                 Attribute = 1
	AttributeIDEncr             Attribute = 2
	AttributeKeyEncr            Attribute = 3
	AttributeDocUIDHash         Attribute = 4
	AttributeDocUIDEncr         Attribute = 5
	AttributeDealCid            Attribute = 6
	AttributeMinerAddress       Attribute = 7
	AttributeContent            Attribute = 8
	AttributeContentEncr        Attribute = 9
	AttributeDescriptionEncr    Attribute = 10
	AttributePasswordHash       Attribute = 11
	AttributeTimestamp          Attribute = 12
	AttributeNameEncr           Attribute = 13
	AttributeGroupDoctorsIDHash Attribute = 14
	AttributeGroupAllDocsIDHash Attribute = 15
	AttributeDataIndexID        Attribute = 16
)

func AttributeGetByCode(attrs any, code uint8) []byte {
	for _, attr := range attrs.(Attributes) {
		if attr.Code == code {
			return attr.Value
		}
	}
	return nil
}
