package parser

import "time"

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

// IntervalsSchedule is a representation of a regular intervals schedule, such as:
//   every 12 hours
//   every 5 minutes from 10:00 to 14:00
//   every day 00:00
type IntervalsSchedule struct {
	Interval int          `json:"interval"`
	TimeUnit TimeUnit     `json:"timeUnit"`
	Months   []time.Month `json:"months"`
}
