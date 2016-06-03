package language

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var eof = rune(0)

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
func (s *Scanner) read() (ch rune, position int) {
	ch, position, err := s.reader.ReadRune()
	if err != nil {
		return eof, position
	}
	return ch, position
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
}

// Scan returns the next token and position from the underlying reader.
// Also returns the literal text read for integers and suffixed integers tokens
// since these token types can have different literal representations.
func (s *Scanner) Scan() (token Token, position int, literal string) {
	// Read next code point.
	ch, position := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter, or certain acceptable special characters, then consume
	// as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanKeyword()
	} else if isDigit(ch) {
		s.unread()
		return s.scanInteger()
	}

	// Otherwise parse individual characters.
	switch ch {
	case eof:
		return EOF, position, ""
	case ',':
		return COMMA, position, ""
	case ':':
		return COLON, position, ""
	}

	return ILLEGAL, position, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (token Token, position int, literal string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	ch, pos := s.read()
	buf.WriteRune(ch)

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch, _ = s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, pos, buf.String()
}

// scanKeyword consumes the current rune and all contiguous letter runes.
func (s *Scanner) scanKeyword() (token Token, position int, literal string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	ch, position := s.read()
	buf.WriteRune(ch)

	// Read every subsequent letter characters into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch, _ = s.read(); ch == eof {
			break
		} else if !isLetter(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	// If the literal matches a keyword then return that keyword.
	if token = Lookup(strings.ToLower(buf.String())); token != ILLEGAL {
		return token, position, buf.String()
	}

	return ILLEGAL, position, buf.String()
}

// scaInteger consumes the current digit and all contiguous digit runes.
func (s *Scanner) scanInteger() (token Token, position int, literal string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	ch, pos := s.read()
	buf.WriteRune(ch)

	// Read every subsequent digit characters into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch, _ = s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	if suffix, err := s.reader.Peek(2); err != nil {
		lowercaseSuffix := strings.ToLower(string(suffix))

		if lowercaseSuffix == "st" || lowercaseSuffix == "nd" ||
			lowercaseSuffix == "rd" || lowercaseSuffix == "th" {
			return SUFFIXEDINTEGER, pos, buf.String()
		}
	}

	return INTEGER, pos, buf.String()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}
