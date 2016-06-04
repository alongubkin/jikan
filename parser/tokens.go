package parser

import "strings"

// Token is a lexical token of the Jikan time expression language.
type Token int

// These are a comprehensive list of Jikan time expression language tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	INTEGER         // 12
	SUFFIXEDINTEGER // 1st, 2nd, 3rd, 4th

	// General keywords
	EVERY // every
	FROM  // from
	TO    // to
	ON    // on
	OF    // of

	COMMA // ,
	COLON // :

	// Units of Time
	SECOND // second, seconds, sec
	MINUTE // minute, minutes, min
	HOUR   // hour, hours
	DAY    // day, days
	MONTH  // month, months
	YEAR   // year, years, yr

	// Days of Week
	SUNDAY    // sun, sunday
	MONDAY    // mon, monday
	TUESDAY   // tue, tuesday
	WEDNESDAY // wed, wednesday
	THURSDAY  // thu, thursday
	FRIDAY    // fri, friday
	SATURDAY  // sat, saturday

	// Months
	JANUARY   // jan, january
	FEBRUARY  // feb, february
	MARCH     // mar, march
	APRIL     // apr, april
	MAY       // may
	JUNE      // june
	JULY      // july
	AUGUST    // aug, august
	SEPTEMBER // sep, september
	OCTOBER   // oct, october
	NOVEMBER  // nov, november
	DECEMBER  // dec, december
)

var keywords = map[string]Token{
	// General keywords
	"every": EVERY,
	"from":  FROM,
	"to":    TO,
	"on":    ON,
	"of":    OF,

	// Units of time
	"sec":     SECOND,
	"second":  SECOND,
	"seconds": SECOND,
	"min":     MINUTE,
	"minute":  MINUTE,
	"minutes": MINUTE,
	"h":       HOUR,
	"hour":    HOUR,
	"hours":   HOUR,
	"day":     DAY,
	"days":    DAY,
	"month":   MONTH,
	"months":  MONTH,
	"yr":      YEAR,
	"year":    YEAR,
	"years":   YEAR,

	// Days of week
	"sun":       SUNDAY,
	"sunday":    SUNDAY,
	"mon":       MONDAY,
	"monday":    MONDAY,
	"tue":       TUESDAY,
	"tuesday":   TUESDAY,
	"wed":       WEDNESDAY,
	"wednesday": WEDNESDAY,
	"thu":       THURSDAY,
	"thursday":  THURSDAY,
	"fri":       FRIDAY,
	"friday":    FRIDAY,
	"sat":       SATURDAY,
	"saturday":  SATURDAY,

	// Months
	"jan":       JANUARY,
	"january":   JANUARY,
	"feb":       FEBRUARY,
	"february":  FEBRUARY,
	"mar":       MARCH,
	"march":     MARCH,
	"apr":       APRIL,
	"april":     APRIL,
	"may":       MAY,
	"june":      JUNE,
	"july":      JULY,
	"aug":       AUGUST,
	"august":    AUGUST,
	"sep":       SEPTEMBER,
	"september": SEPTEMBER,
	"oct":       OCTOBER,
	"october":   OCTOBER,
	"nov":       NOVEMBER,
	"november":  NOVEMBER,
	"dec":       DECEMBER,
	"december":  DECEMBER,
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
	if token, ok := keywords[strings.ToLower(ident)]; ok {
		return token
	}

	return ILLEGAL
}
