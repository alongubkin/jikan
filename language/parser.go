package language

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// Parser represents a parser.
type Parser struct {
	scanner *Scanner
	buffer  struct {
		token    Token // last read token
		position int
		literal  string // last read literal
		size     int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{scanner: NewScanner(r)}
}

// ParseIntervalsSchedule is
func (p *Parser) ParseIntervalsSchedule() (*IntervalsSchedule, error) {
	// EVERY NUMBER? TimeUnit [OF|ON MonthSpec]? [TimeOfDay|TimeRange]?

	schedule := &IntervalsSchedule{Interval: 1}
	var err error

	// Scan the EVERY token.
	if token, position, literal := p.scanIgnoreWhitespace(); token != EVERY {
		return nil, fmt.Errorf("%d: Found %q, expected \"every\"",
			position, literal)
	}

	// Scan the optional time interval INTEGER token.
	if token, position, literal := p.scanIgnoreWhitespace(); token == INTEGER {
		if schedule.Interval, err = strconv.Atoi(literal); err != nil {
			return nil, fmt.Errorf("%d: Cannot convert %q to integer",
				position, literal)
		}
	} else {
		p.unscan()
	}

	// Scan time unit.
	if schedule.TimeUnit, err = p.parseTimeUnit(); err != nil {
		return nil, err
	}

	// Try to scan optional month spec.
	if token, position, _ := p.scanIgnoreWhitespace(); token == OF || token == ON {
		// Make sure that the current time unit is allowed for month specification
		switch schedule.TimeUnit {
		case Seconds, Minutes, Hours, Days:
			break

		default:
			return nil, fmt.Errorf("%d: Cannot use \"of [month, ...]\" with this time unit.",
				position)
		}
		// Try to scan "month" keyword, which means all months.
		if token, _, _ := p.scanIgnoreWhitespace(); token == MONTH {
			schedule.Months = []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		} else {
			// The "month" keyword doesn't exist, so there should be a list of months.
			p.unscan()
			if schedule.Months, err = p.parseMonthList(); err != nil {
				return nil, err
			}
		}
	} else {
		p.unscan()
	}

	return schedule, nil
}

func (p *Parser) parseTimeUnit() (TimeUnit, error) {
	switch token, position, literal := p.scanIgnoreWhitespace(); token {
	case SECOND:
		return Seconds, nil
	case MINUTE:
		return Minutes, nil
	case HOUR:
		return Hours, nil
	case DAY:
		return Days, nil
	case MONTH:
		return Months, nil
	case YEAR:
		return Years, nil

	default:
		return "", fmt.Errorf("%d: Found %q, expected time unit (sec, hours, days, etc...)",
			position, literal)
	}
}

func (p *Parser) parseMonthList() ([]time.Month, error) {
	months := make([]time.Month, 1, 12)
	var err error

	// Scan first month, which is required
	if months[0], err = p.parseMonth(); err != nil {
		return nil, err
	}

	for {
		// Scan next token, if it's not a comma then stop
		if token, _, _ := p.scanIgnoreWhitespace(); token != COMMA {
			p.unscan()
			break
		}

		// Parse next month
		month, err := p.parseMonth()
		if err != nil {
			return nil, err
		}

		months = append(months, month)
	}

	return months, nil
}

func (p *Parser) parseMonth() (time.Month, error) {
	switch token, position, literal := p.scanIgnoreWhitespace(); token {
	case JANUARY:
		return time.January, nil
	case FEBRUARY:
		return time.February, nil
	case MARCH:
		return time.March, nil
	case APRIL:
		return time.April, nil
	case MAY:
		return time.May, nil
	case JUNE:
		return time.June, nil
	case JULY:
		return time.July, nil
	case AUGUST:
		return time.August, nil
	case SEPTEMBER:
		return time.September, nil
	case OCTOBER:
		return time.October, nil
	case NOVEMBER:
		return time.November, nil
	case DECEMBER:
		return time.December, nil

	default:
		return 0, fmt.Errorf("%d: Found %q, expected month (july, may, etc...)",
			position, literal)
	}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (token Token, position int, literal string) {
	// If we have a token on the buffer, then return it.
	if p.buffer.size != 0 {
		p.buffer.size = 0
		return p.buffer.token, p.buffer.position, p.buffer.literal
	}

	// Otherwise read the next token from the scanner
	token, position, literal = p.scanner.Scan()

	// Save it to the buffer in case we unscan later.
	p.buffer.token, p.buffer.position, p.buffer.literal = token, position, literal

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.buffer.size = 1
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (token Token, position int, literal string) {
	token, position, literal = p.scan()
	if token == WS {
		token, position, literal = p.scan()
	}
	return
}
