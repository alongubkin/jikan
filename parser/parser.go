package parser

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// Parser represents a parser.
type Parser struct {
	lexer  *Lexer
	buffer struct {
		token    Token // last read token
		position int
		literal  string // last read literal
		size     int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{lexer: NewLexer(r)}
}

// ParseIntervalsSchedule parses the following syntaxes:
//   every month
//   every 5 minutes
//   every 3 sec of jan
//   every 7 hours of jan, feb, april
//   every 5 hours of month (identical to "every 5 hours")
func (p *Parser) ParseIntervalsSchedule() (*IntervalsSchedule, error) {
	schedule := &IntervalsSchedule{Interval: 1}
	var err error

	// An intervals schedule starts with an "every" token.
	if token, position, literal := p.scanIgnoreWhitespace(); token != EVERY {
		return nil, fmt.Errorf("%d: Found %q, expected \"every\"",
			position, literal)
	}

	// Scan the optional time interval (the "5" in "every 5 minutes").
	// It's optional because the following syntaxes are allowed:
	// 		Every minute
	//		Every 5 minutes
	if token, position, literal := p.scanIgnoreWhitespace(); token == INTEGER {
		if schedule.Interval, err = strconv.Atoi(literal); err != nil {
			return nil, fmt.Errorf("%d: Cannot convert %q to integer",
				position, literal)
		}
	} else {
		p.unscan()
	}

	// Scan time unit (the "sec" in "every 5 sec").
	if schedule.TimeUnit, err = p.parseTimeUnit(); err != nil {
		return nil, err
	}

	// Try to scan optional month spec (e.g: "every 5 hours of january")
	if token, position, _ := p.scanIgnoreWhitespace(); token == OF || token == ON {
		// Make sure that the current time unit is allowed for month specification
		// (e.g: "every 5 years of jan, feb" is not allowed)
		switch schedule.TimeUnit {
		case Seconds, Minutes, Hours, Days:
			break

		default:
			return nil, fmt.Errorf("%d: Cannot use \"of [month, ...]\" with this time unit",
				position)
		}

		// Check if the "month" keyword exists (e.g: "every 5 minutes of month").
		// It means all months.
		if token, _, _ := p.scanIgnoreWhitespace(); token != MONTH {
			// The "month" keyword doesn't exist, so there should be a list of months
			// (e.g: "every 5 minutes of jan, feb, march")
			p.unscan()
			if schedule.Months, err = p.parseMonthList(); err != nil {
				return nil, err
			}

			// If the month list is a list of all months (e.g: "jan, feb, march, ..., dec"),
			// then just turn it to nil, which is the same thing.
			if isListOfAllMonths(schedule.Months) {
				schedule.Months = nil
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

// scan returns the next token from the underlying lexer.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (token Token, position int, literal string) {
	// If we have a token on the buffer, then return it.
	if p.buffer.size != 0 {
		p.buffer.size = 0
		return p.buffer.token, p.buffer.position, p.buffer.literal
	}

	// Otherwise read the next token from the lexer
	token, position, literal = p.lexer.Scan()

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

func isListOfAllMonths(months []time.Month) bool {
	for month := time.January; month <= time.December; month++ {
		found := false

		for _, value := range months {
			if month == value {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
