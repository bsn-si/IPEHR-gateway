package base

// Pathable
// The PATHABLE class defines the pathing capabilities used by nearly all classes in the openEHR
// reference model, mostly via inheritance of LOCATABLE. The defining characteristics of PATHABLE objects
// are that they can locate child objects using paths, and they know their parent object in a compositional
// hierarchy. The parent feature is defined as abstract in the model, and may be implemented in any way convenient.
// https://specifications.openehr.org/releases/RM/Release-1.0.4/common.html#_pathable_class
type Pathable interface{}
