package base

// ContentItem
// Abstract ancestor of all concrete content types.
// https://specifications.openehr.org/releases/RM/Release-1.0.2/ehr.html#_content_item_class
type ContentItem interface {
	GetType() string
}

type ContentItemType string

func (cit ContentItemType) ToString() string {
	return string(cit)
}

const (
	ActionContentItemType          ContentItemType = "ACTION"
	ActivityContentItemType        ContentItemType = "ACTIVITY"
	ArchetypedContentItemType      ContentItemType = "ARCHETYPED"
	ArchetypeIDContentItemType     ContentItemType = "ARCHETYPE_ID"
	ClusterContentItemType         ContentItemType = "CLUSTER"
	CodePhraseContentItemType      ContentItemType = "CODE_PHRASE"
	CompositionContentItemType     ContentItemType = "COMPOSITION"
	DvBooleanContentItemType       ContentItemType = "DV_BOOLEAN"
	DvCodedTextContentItemType     ContentItemType = "DV_CODED_TEXT"
	DvCountContentItemType         ContentItemType = "DV_COUNT"
	DvDateContentItemType          ContentItemType = "DV_DATE"
	DvDateTimeContentItemType      ContentItemType = "DV_DATE_TIME"
	DvDurationContentItemType      ContentItemType = "DV_DURATION"
	DvIdentifierContentItemType    ContentItemType = "DV_IDENTIFIER"
	DvMultimediaContentItemType    ContentItemType = "DV_MULTIMEDIA"
	DvParsableContentItemType      ContentItemType = "DV_PARSABLE"
	DvProportionContentItemType    ContentItemType = "DV_PROPORTION"
	DvQuantityContentItemType      ContentItemType = "DV_QUANTITY"
	DvTextContentItemType          ContentItemType = "DV_TEXT"
	DvTimeContentItemType          ContentItemType = "DV_TIME"
	DvURIContentItemType           ContentItemType = "DV_URI"
	ElementContentItemType         ContentItemType = "ELEMENT"
	EvaluationContentItemType      ContentItemType = "EVALUATION"
	EventContextContentItemType    ContentItemType = "EVENT_CONTEXT"
	GenericIDContentItemType       ContentItemType = "GENERIC_ID"
	HierObjectIDContentItemType    ContentItemType = "HIER_OBJECT_ID"
	HistoryContentItemType         ContentItemType = "HISTORY"
	InstructionContentItemType     ContentItemType = "INSTRUCTION"
	IsmTransitionContentItemType   ContentItemType = "ISM_TRANSITION"
	ItemTreeContentItemType        ContentItemType = "ITEM_TREE"
	ObjectVersionIDContentItemType ContentItemType = "OBJECT_VERSION_ID"
	ObservationContentItemType     ContentItemType = "OBSERVATION"
	ParticipationContentItemType   ContentItemType = "PARTICIPATION"
	PartyIdentifiedContentItemType ContentItemType = "PARTY_IDENTIFIED"
	PartyRefContentItemType        ContentItemType = "PARTY_REF"
	PartySelfContentItemType       ContentItemType = "PARTY_SELF"
	PointEventContentItemType      ContentItemType = "POINT_EVENT"
	SectionContentItemType         ContentItemType = "SECTION"
	TemplateIDContentItemType      ContentItemType = "TEMPLATE_ID"
	TerminologyIDContentItemType   ContentItemType = "TERMINOLOGY_ID"
)
