package ast

// TimeUnit represents any particular time interval.
type TimeUnit string

// Time unit definitions
const (
	Seconds TimeUnit = "seconds"
	Minutes TimeUnit = "minutes"
	Hours   TimeUnit = "hours"
	Days    TimeUnit = "days"
	Months  TimeUnit = "months"
	Years   TimeUnit = "years"
)

// TimeDuration is.
type TimeDuration struct {
	Unit   TimeUnit `json:"time_unit"`
	Length int      `json:"length"`
}

// EveryStatement represents
type EveryStatement struct {
	Duration interface{} `json:"duration"`
}
