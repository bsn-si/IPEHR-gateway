package base

import (
	"encoding/json"
)

// DataValue
// Abstract parent of all DV_ data value types. Serves as a common ancestor of all data value types in openEHR models.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_data_value_class
type DataValue interface {
	GetType() ItemType
}

type dataValueWrapper struct {
	dv DataValue
}

func (dvw *dataValueWrapper) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type ItemType `json:"_type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	switch tmp.Type {
	case DvURIItemType:
		dvw.dv = &DvURI{}
	case DvTimeItemType:
		dvw.dv = &DvTime{}
	case DvQuantityItemType:
		dvw.dv = &DvQuantity{}
	case DvStateItemType:
		dvw.dv = &DvState{}
	case DvProportionItemType:
		dvw.dv = &DvProportion{}
	case DvParsableItemType:
		dvw.dv = &DvParsable{}
	case DvParagraphItemType:
		dvw.dv = &DvParagraph{}
	case DvMultimediaItemType:
		dvw.dv = &DvMultimedia{}
	case DvIdentifierItemType:
		dvw.dv = &DvIdentifier{}
	case DvDurationItemType:
		dvw.dv = &DvDuration{}
	case DvDateTimeItemType:
		dvw.dv = &DvDateTime{}
	case DvDateItemType:
		dvw.dv = &DvDate{}
	case DvCountItemType:
		dvw.dv = &DvCount{}
	case DvCodedTextItemType:
		dvw.dv = &DvCodedText{}
	case DvTextItemType:
		dvw.dv = &DvText{}
	case DvBooleanItemType:
		dvw.dv = &DvBoolean{}
	}

	if err := json.Unmarshal(data, dvw.dv); err != nil {
		return err
	}

	return nil
}

type DvValueBase struct {
	Type ItemType `json:"_type"`
}

func (dv DvValueBase) GetType() ItemType {
	return dv.Type
}

// DvBoolean
// Items which are truly boolean data, such as true/false or yes/no answers.
// For such data, it is important to devise the meanings (usually questions in subjective data) carefully,
// so that the only allowed results are in fact true or false.
// Misuse: The DV_BOOLEAN class should not be used as a replacement for naively modelled enumerated types such as male/female etc.
// Such values should be coded, and in any case the enumeration often has more than two values.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_boolean_class
type DvBoolean struct {
	DvValueBase
	Value bool `json:"value"`
}

// DvIdentifier
// Type for representing identifiers of real-world entities. Typical identifiers include drivers licence
// number, social security number, veterans affairs number, prescription id, order id, and so on.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_identifier_class
type DvState struct {
	DvValueBase
	Value      DvCodedText `json:"value"`
	IsTerminal bool        `json:"is_terminal"`
}

// DvEncapsulated
// Abstract class defining the common meta-data of all types of encapsulated data.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_encapsulated_class
type DvEncapsulated struct {
	DvValueBase
	Charset  *CodePhrase `json:"charset,omitempty"`
	Language *CodePhrase `json:"language,omitempty"`
}

// DvMultimedia
//
// A specialisation of DV_ENCAPSULATED for audiovisual and bio-signal types.
// Includes further metadata relating to multimedia types which are not applicable to other subtypes of DV_ENCAPSULATED.
//
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_multimedia_class
type DvMultimedia struct {
	DvEncapsulated
	AlternativeText         string        `json:"alternative_text"`
	URI                     *DvURI        `json:"uri,omitempty"`
	Data                    []byte        `json:"data,omitempty"`
	MediaType               CodePhrase    `json:"media_type"`
	CompressionAlgorithm    *CodePhrase   `json:"compression_algorithm,omitempty"`
	IntegrityCheck          []byte        `json:"integrity_check,omitempty"`
	IntegrityCheckAlgorithm *CodePhrase   `json:"integrity_check_algorithm,omitempty"`
	Thumbnail               *DvMultimedia `json:"thumbnail,omitempty"`
	Size                    int           `json:"size"`
}

// DvParsable
// Encapsulated data expressed as a parsable String. The internal model of the data item is not described in the
// openEHR model in common with other encapsulated types, but in this case, the form of the data is assumed to be plaintext,
// rather than compressed or other types of large binary data.
//
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_parsable_class
type DvParsable struct {
	DvEncapsulated
	Value     string `json:"value"`
	Formalism string `json:"formalism"`
}

// DvIdentifier
// Type for representing identifiers of real-world entities. Typical identifiers include drivers licence
// number, social security number, veterans affairs number, prescription id, order id, and so on.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_identifier_class
type DvIdentifier struct {
	DvValueBase
	Issuer   string `json:"issuer,omitempty"`
	Assigner string `json:"assigner,omitempty"`
	ID       string `json:"id"`
	Type     string `json:"type,omitempty"`
}

// DvInterval
// Generic class defining an interval (i.e. range) of a comparable type.
// An interval is a contiguous subrange of a comparable base type. Used to define intervals of dates, times,
// quantities (whose units match) and so on. The type parameter, T, must be a descendant of the type DV_ORDERED,
// which is necessary (but not sufficient) for instances to be compared (strictly_comparable is also needed).
//
// Without the DV_INTERVAL class, quite a few more DV_ classes would be needed to express logical intervals,
// namely interval versions of all the date/time classes, and of quantity classes.
// Further, it allows the semantics of intervals to be stated in one place unequivocally, including the conditions for strict comparison.
//
// The basic semantics are derived from the class Interval<T>, described in the support RM.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_interval_class
type DvInterval[T any] struct {
	DvValueBase
	Interval[T]
}

func (interval *DvInterval[T]) GetType() ItemType {
	return DvIntervalItemType
}

// DvOrdered
// Abstract class defining the concept of ordered values, which includes ordinals as well as true
// quantities.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_ordered_class
type DvOrdered[T any] struct {
	DvValueBase
	NormalStatus         *CodePhrase         `json:"normal_status,omitempty"`
	NormalRange          *DvInterval[T]      `json:"normal_range,omitempty"`
	OtherReferenceRanges []ReferenceRange[T] `json:"other_reference_ranges,omitempty"`
}

// DvQuantified
// Abstract class defining the concept of true quantified values, i.e. values which are not only ordered,
// but which have a precise magnitude.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_quantified_class
type DvQuantified[T any] struct {
	MagnitudeStatus bool `json:"magnitude_status,omitempty"`
	Accuracy        any  `json:"accuracy,omitempty"`
	DvOrdered[T]
}

// DvAmount
// Abstract class defining the concept of relative quantified 'amounts'. For relative quantities,
// the + and - operators are defined (unlike descendants of DV_ABSOLUTE_QUANTITY, such as the date/time
// types).
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_amount_class
type DvAmount[T any] struct {
	AccuracyIsPercent bool    `json:"accuracy_is_percent,omitempty"`
	Accuracy          float32 `json:"accuracy,omitempty"`
	DvQuantified[T]
}

// DvQuantity
// Quantitified type representing scientific quantities, i.e.
// quantities expressed as a magnitude and units. Units are expressed in the UCUM syntax
// (Unified Code for Units of Measure (UCUM), by Gunther Schadow and Clement J. McDonald of The Regenstrief Institute)
// (case-sensitive form) by default, or another system if units_system is set.
//
// Can also be used for time durations, where it is more convenient to treat these as simply a
// number of seconds rather than days, months, years (in the latter case, DV_DURATION may be used).
//
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_quantity_class
type DvQuantity struct {
	DvAmount[int64]
	Magnitude            float64                      `json:"magnitude"`
	Precision            *int                         `json:"precision,omitempty"`
	Units                *string                      `json:"units,omitempty"`
	UnitsSystem          *string                      `json:"units_system,omitempty"`
	UnitsDisplayName     *string                      `json:"units_display_name,omitempty"`
	NormalRange          *DvInterval[DvQuantity]      `json:"normal_range,omitempty"`
	OtherReferenceRanges []ReferenceRange[DvQuantity] `json:"other_reference_ranges,omitempty"`
}

// DvCount
// Countable quantities. Used for countable types such as pregnancies and steps (taken by a physiotherapy patient),
// number of cigarettes smoked in a day.
//
// Misuse: Not to be used for amounts of physical entities (which all have units).
//
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_count_class
type DvCount struct {
	DvAmount[int64]
	Magnitude            int64                     `json:"magnitude"`
	NormalRange          *DvInterval[DvCount]      `json:"normal_range,omitempty"`
	OtherReferenceRanges []ReferenceRange[DvCount] `json:"other_reference_ranges,omitempty"`
}

// DvProportion
// Models a ratio of values, i.e. where the numerator and denominator are both pure numbers.
// The valid_proportion_kind property of the PROPORTION_KIND class is used to control the type attribute to be one of a defined set.
//
// Used for recording titers (e.g. 1:128), concentration ratios, e.g. Na:K (unitary denominator),
// albumin:creatinine ratio, and percentages, e.g. red cell distribution width (RDW).
//
// Misuse: Should not be used to represent things like blood pressure which are often written using a '/' character,
// giving the misleading impression that the item is a ratio, when in fact it is a structured value.
// Similarly, visual acuity, often written as (e.g.) "6/24" in clinical notes is not a ratio but an ordinal
// (which includes non-numeric symbols like CF = count fingers etc). Should not be used for formulations.
//
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_proportion_class
type DvProportion struct {
	DvAmount[int64]
	Numeration           float64                        `json:"numeration"`
	Denomination         float64                        `json:"denomination"`
	Type                 int                            `json:"type"`
	Precision            *int                           `json:"precision,omitempty"`
	NormalRange          *DvInterval[DvProportion]      `json:"normal_range,omitempty"`
	OtherReferenceRanges []ReferenceRange[DvProportion] `json:"other_reference_ranges,omitempty"`
}

// DvDuration
// Represents a period of time with respect to a notional point in time, which is not specified.
// A sign may be used to indicate the duration is backwards in time rather than forwards.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_duration_class
type DvDuration struct {
	Value string `json:"value"`
	DvAmount[int64]
}

// DvTemporal
// Specialised temporal variant of DV_ABSOLUTE_QUANTITY whose diff type is DV_DURATION.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_temporal_class
type DvTemporal struct {
	DvValueBase
	Accuracy *DvDuration `json:"accuracy,omitempty"`
}

// DvDate
// Represents an absolute point in time, as measured on the Gregorian calendar, and specified only to the day.
// Semantics defined by ISO 8601. Used for recording dates in real world time.
// The partial form is used for approximate birth dates, dates of death, etc.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_date_class
type DvDate struct {
	DvTemporal
	Value string `json:"value"`
}

// DvTime
// Represents an absolute point in time from an origin usually interpreted as meaning the start of the current day,
// specified to a fraction of a second. Semantics defined by ISO 8601.
//
// Used for recording real world times, rather than scientifically measured fine amounts of time.
// The partial form is used for approximate times of events and substance administrations.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_time_class
type DvTime struct {
	DvTemporal
	Value string `json:"value"`
}

// DvDateTime
// Represents an absolute point in time, specified to the second. Semantics defined by ISO 8601.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_date_time_class
type DvDateTime struct {
	DvTemporal
	Value string `json:"value"`
}

// DvText
// A text item, which may contain any amount of legal characters arranged as e.g. words, sentences etc
// (i.e. one DV_TEXT may be more than one word). Visual formatting and hyperlinks may be included via
// markdown.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_text_class
type DvText struct {
	DvValueBase
	Value      string        `json:"value"`
	Hyperlink  *DvURI        `json:"hyperlink,omitempty"`
	Formatting string        `json:"formatting,omitempty"`
	Mappings   []TermMapping `json:"mappings,omitempty"`
	Language   *CodePhrase   `json:"language,omitempty"`
	Encoding   *CodePhrase   `json:"encoding,omitempty"`
}

func NewDvText(value string) DvText {
	return DvText{
		DvValueBase: DvValueBase{Type: DvTextItemType},
		Value:       value,
	}
}

// DvCodedText
// A text item whose value must be the rubric from a controlled terminology, the key (i.e. the 'code') of
// which is the defining_code attribute. In other words: a DV_CODED_TEXT is a combination of a CODE_PHRASE
// (effectively a code) and the rubric of that term, from a terminology service, in the language in which
// the data were authored.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_coded_text_class
type DvCodedText struct {
	DefiningCode CodePhrase `json:"defining_code"`
	DvText
}

func NewDvCodedText(value string, codePhrase CodePhrase) DvCodedText {
	return DvCodedText{
		DefiningCode: codePhrase,
		DvText: DvText{
			DvValueBase: DvValueBase{
				Type: DvCodedTextItemType,
			},
			Value: value,
		},
	}
}

// DvParagraph
//
// DEPRECATED: use markdown formatted DV_TEXT instead.
//
// Original definition:
//
// A logical composite text value consisting of a series of DV_TEXTs, i.e. plain text (optionally coded)
// potentially with simple formatting, to form a larger tract of prose, which may be interpreted for display purposes as a paragraph.
//
// DV_PARAGRAPH is the standard way for constructing longer text items in summaries, reports and so on.
//
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_paragraph_class
type DvParagraph struct {
	DvValueBase
	Items []DvText `json:"items"`
}

// DvURI
// A reference to an object which structurally conforms to the Universal Resource Identifier (URI)
// RFC-3986 standard.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_uri_class
type DvURI struct {
	DvValueBase
	Value string `json:"value"`
}
