package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Token represents a lexical token.
type Token int

// once a day at noon
// every Saturday, and on the 30th day of the month
// every hour on the hour
// X every 12 hours
// X every 5 minutes from 10:00 to 14:00
// X every day 00:00
// X every monday 09:00
// 2nd,third mon,wed,thu of march 17:00
// 1st monday of sep,oct,nov 17:00
// 1 of jan,april,july,oct 00:00
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	INTEGER

	// Durations
	EVERY

	// Days
	SUNDAY
	MONDAY
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY

	// Time units
	SECONDS
	MINUTES
	HOURS
	DAYS
	MONTHS
	YEARS
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

// Scanner represents a lexical scanner.
type Scanner struct {
	reader *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{reader: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (token Token, literal string) {
	// Read the next rune.
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanKeyword()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (token Token, literal string) {
	// Create a buffer to store the whitespace in and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanKeyword() (token Token, literal string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToLower(buf.String()) {
	case "every":
		return EVERY, buf.String()

	// Days
	case "sunday", "sun":
		return SUNDAY, buf.String()
	case "monday", "mon":
		return MONDAY, buf.String()
	case "tuesday", "tue":
		return TUESDAY, buf.String()
	case "wednesday", "wed":
		return WEDNESDAY, buf.String()
	case "thursday", "thu":
		return THURSDAY, buf.String()
	case "friday", "fri":
		return FRIDAY, buf.String()
	case "saturday", "sat":
		return SATURDAY, buf.String()

	// Time units
	case "seconds", "second", "sec":
		return SECONDS, buf.String()
	case "minutes", "minute", "min":
		return MINUTES, buf.String()
	case "hours", "hour", "h":
		return HOURS, buf.String()
	case "days", "day":
		return DAYS, buf.String()
	case "months", "month":
		return MONTHS, buf.String()
	case "years", "year", "yr":
		return YEARS, buf.String()
	}

	return ILLEGAL, strings.ToUpper(buf.String())
}

// scanNumber consumes the current rune and all contiguous digit runes.
func (s *Scanner) scanNumber() (token Token, literal string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent digit character into the buffer.
	// Non-digit characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return INTEGER, buf.String()
}
