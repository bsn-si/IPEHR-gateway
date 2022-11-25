package base

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// History
// Root object of a linear history, i.e. time series structure.
// This is a generic class whose type parameter must be a descendant of ITEM_STRUCTURE,
// ensuring that each Event in the events of a given instance is of the same structural type, i.e. ITEM_TREE, ITEM_LIST etc.
//
// For a periodic series of events, period will be set, and the time of each Event in the History must correspond;
// i.e. the EVENT.offset must be a multiple of period for each Event. Missing events in a period History are however allowed.
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_history_class
type History[T any] struct {
	DataStructure
	Origin   DvDateTime     `json:"origin"`
	Period   *DvDuration    `json:"period,omitempty"`
	Duration *DvDuration    `json:"duration,omitempty"`
	Summary  *ItemStructure `json:"summary,omitempty"`
	Events   []Event[T]     `json:"events,omitempty"`
}

func (h History[T]) GetType() ItemType {
	return HistoryItemType
}

// Event
// Defines the abstract notion of a single event in a series. This class is generic,
// allowing types to be generated which are locked to particular spatial types, such as EVENT<ITEM_LIST>. Subtypes express point or intveral data.
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_event_class
type Event[T any] struct {
	Data Root `json:"-"`
}

func (e Event[T]) GetType() ItemType {
	return e.Data.GetType()
}

func (e Event[T]) GetLocatable() Locatable {
	return e.Data.GetLocatable()
}

func (e Event[T]) GetArchetypeNodeID() string {
	return e.Data.GetArchetypeNodeID()
}

func (e Event[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Data)
}

func (e *Event[T]) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type ItemType `json:"_type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "cannot unmarshal Event type")
	}

	switch tmp.Type {
	case PointEventItemType:
		e.Data = &PointEvent[T]{}
	case IntervalEventItemType:
		e.Data = &IntervalEvent[T]{}
	default:
		return fmt.Errorf("unexpected Event type: %v", tmp.Type) // nolint
	}

	if err := json.Unmarshal(data, e.Data); err != nil {
		return errors.Wrap(err, "cannot unmarshal item structure instance")
	}

	return nil
}

type BaseEvent[T any] struct { // nolint
	Locatable
	Time  DvDateTime     `json:"time"`
	State *ItemStructure `json:"state,omitempty"`
	Data  T              `json:"data"`
}

// PointEvent
//
// Defines a single point event in a series.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_point_event_class
type PointEvent[T any] struct {
	BaseEvent[T]
}

// IntervalEvent
//
// Defines a single interval event in a series.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_interval_event_class
type IntervalEvent[T any] struct {
	BaseEvent[T]
	Width        DvDuration  `json:"width"`
	SampleCount  *int        `json:"sample_count"`
	MathFunction DvCodedText `json:"math_function"`
}
