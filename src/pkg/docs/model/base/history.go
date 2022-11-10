package base

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
	Events   []Event[T]     `json:"event,omitempty"`
}

func (h History[T]) GetType() ItemType {
	return HistoryItemType
}

// Event
// Defines the abstract notion of a single event in a series. This class is generic,
// allowing types to be generated which are locked to particular spatial types, such as EVENT<ITEM_LIST>. Subtypes express point or intveral data.
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_event_class
type Event[T any] struct {
	Locatable
	Time  DvDateTime     `json:"time"`
	State *ItemStructure `json:"state"`
	Data  T              `json:"data"`
}
