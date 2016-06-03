package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/alongubkin/jikan/ast"
)

// DayExpression =
//		SUNDAY|MONDAY|TUESDAY|WEDNESDAY|THURSDAY|FRIDAY|SATURDAY

// TimeUnit =
//		SECONDS|MINUTES|HOURS|DAYS|MONTHS|YEARS

// TimeDuration =
// 		INTEGER? TimeUnit

// Time =
//		INTEGER(0 <= x <= 60) COLON INTEGER(0 <= x <= 60)

// TimeRange =
//		FROM Time TO Time

// OnExpression
//		ON (Time|

// EveryStatement =
//		every (TimeDuration|DayExpression) (Time|TimeRange|OnExpression)?

// Parser represents a parser.
type Parser struct {
	scanner *Scanner
	buffer  struct {
		token   Token  // last read token
		literal string // last read literal
		size    int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{scanner: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (token Token, literal string) {
	// If we have a token on the buffer, then return it.
	if p.buffer.size != 0 {
		p.buffer.size = 0
		return p.buffer.token, p.buffer.literal
	}

	// Otherwise read the next token from the scanner.
	token, literal = p.scanner.Scan()

	// Save it to the buffer in case we unscan later.
	p.buffer.token, p.buffer.literal = token, literal

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buffer.size = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (token Token, literal string) {
	token, literal = p.scan()
	if token == WS {
		token, literal = p.scan()
	}
	return
}

// Parse is
func (p *Parser) Parse() (*ast.EveryStatement, error) {
	statement := &ast.EveryStatement{}
	var err error

	if token, literal := p.scanIgnoreWhitespace(); token != EVERY {
		return nil, fmt.Errorf("Found %q, expected \"every\".", literal)
	}

	// Parse time duration
	if statement.Duration, err = p.parseTimeDuration(); err != nil {
		return nil, err
	}

	return statement, nil
}

func (p *Parser) parseTimeDuration() (*ast.TimeDuration, error) {
	duration := &ast.TimeDuration{}
	var err error

	// Parse optional time length
	if token, literal := p.scanIgnoreWhitespace(); token == INTEGER {
		if duration.Length, err = strconv.Atoi(literal); err != nil {
			return nil, fmt.Errorf("Cannot convert %q to integer.", literal)
		}
	} else {
		p.unscan()
	}

	// Parse time unit
	if duration.Unit, err = p.parseTimeUnit(); err != nil {
		return nil, err
	}

	return duration, nil
}

func (p *Parser) parseTimeUnit() (ast.TimeUnit, error) {
	token, literal := p.scanIgnoreWhitespace()
	switch token {
	case SECONDS:
		return ast.Seconds, nil
	case MINUTES:
		return ast.Minutes, nil
	case HOURS:
		return ast.Hours, nil
	case DAYS:
		return ast.Days, nil
	case MONTHS:
		return ast.Months, nil
	case YEARS:
		return ast.Years, nil
	}

	return "", fmt.Errorf("Found %q, expected a time unit (sec, days, hours, etc).", literal)
}
