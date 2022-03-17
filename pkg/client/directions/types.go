package directions

// DirectionResponse is the response from GetDirections
type DirectionResponse struct {
	Code      string
	Waypoints []Waypoint
	Routes    []Route
}

// Route A route through (potentially multiple) waypoints.
type Route struct {
	Distance float64
	Duration float64
	Geometry string
	Legs     []RouteLeg
}

// Waypoint is an input point snapped to the road network
type Waypoint struct {
	Name     string
	Location []float64
}

// RouteLeg A route between two Waypoints
type RouteLeg struct {
	Distance   float64
	Duration   float64
	Steps      []RouteStep
	Summary    string
	Annotation Annotation
}

// Annotation conains additional details about each line segment
type Annotation struct {
	Distance []float64
	Duration []float64
	Speed    []float64
}

// RouteStep Includes one StepManeuver object and travel to the following RouteStep.
type RouteStep struct {
	Distance      float64
	Duration      float64
	Geometry      string
	Name          string
	Ref           string
	Destinations  string
	Mode          TransportationMode
	Maneuver      StepManeuver
	Intersections []Intersection
}

// TransportationMode indicates the mode of transportation
type TransportationMode string

const (
	ModeDriving      TransportationMode = "driving"
	ModeWalking      TransportationMode = "walking"
	ModeFerry        TransportationMode = "ferry"
	ModeCycling      TransportationMode = "cyling"
	ModeUnaccessible TransportationMode = "unaccessible"
)

// Intersection
type Intersection struct {
	Location []float64
	Bearings []float64
	Entry    []bool
	In       uint
	Out      uint
	Lanes    []Lane
}

// Lane
type Lane struct {
	Valid      bool
	Indicatons []string
}

// StepManeuver
type StepManeuver struct {
	Location      []float64
	BearingBefore float64
	BearingAfter  float64
	Instruction   string
	Type          string
	Modifier      StepModifier
}

// StepModifier indicates the direction change of the maneuver
type StepModifier string

const (
	StepModifierUTurn       StepModifier = "uturn"
	StepModifierSharpRight  StepModifier = "sharp right"
	StepModifierRight       StepModifier = "right"
	StepModifierSlightRight StepModifier = "slight right"
	StepModifierStraight    StepModifier = "straight"
	StepModifierSharpLeft   StepModifier = "sharp left"
	StepModifierLeft        StepModifier = "left"
	StepModifierSlightLeft  StepModifier = "slight left"
)

// Codes are direction response Codes
type Codes string

const (
	CodeOK              Codes = "Ok"
	CodeNoRoute         Codes = "NoRoute"
	CodeNoSegment       Codes = "NoSegment"
	CodeProfileNotFound Codes = "ProfileNotFound"
	CodeInvalidInput    Codes = "InvalidInput"
)
