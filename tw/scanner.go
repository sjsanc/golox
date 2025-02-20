package tw

import (
	"log"
	"strconv"
)

var keywords = map[string]TokenType{
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source   string
	tokens   []*Token
	keywords map[string]TokenType
	start    int
	current  int
	line     int
	hadErr   bool
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		source:   src,
		tokens:   make([]*Token, 0),
		keywords: keywords,
		line:     1,
	}
}

func (s *Scanner) Scan() ([]*Token, bool) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens, s.hadErr
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case "(":
		s.addToken(LEFT_PAREN, nil)
	case ")":
		s.addToken(RIGHT_PAREN, nil)
	case "{":
		s.addToken(LEFT_BRACE, nil)
	case "}":
		s.addToken(RIGHT_BRACE, nil)
	case ",":
		s.addToken(COMMA, nil)
	case ".":
		s.addToken(DOT, nil)
	case "-":
		s.addToken(MINUS, nil)
	case "+":
		s.addToken(PLUS, nil)
	case ";":
		s.addToken(SEMICOLON, nil)
	case "*":
		s.addToken(STAR, nil)
	case "!":
		s.matchElse("=", BANG_EQUAL, BANG)
	case "=":
		s.matchElse("=", EQUAL_EQUAL, EQUAL)
	case "<":
		s.matchElse("=", LESS_EQUAL, LESS)
	case ">":
		s.matchElse("=", GREATER_EQUAL, GREATER)
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case " ", "\r", "\t":
	case "\n":
		s.line++
	case "\"":
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.error(s.line, "Unexpected character.")
		}
	}
}
func (s *Scanner) addToken(tt TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tt, text, literal, s.line))
}

// ================================================================================
// ### TOKENS
// ================================================================================

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	ttype, ok := s.keywords[text]
	if !ok {
		ttype = IDENTIFIER
	}
	s.addToken(ttype, nil)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == "." && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	if value, err := strconv.Atoi(s.source[s.start:s.current]); err != nil {
		s.error(s.line, "error parsing number")
		return
	} else {
		s.addToken(NUMBER, value)
	}
}

func (s *Scanner) string() {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.error(s.line, "unterminated string")
		return
	}
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

// ================================================================================
// ### HELPERS
// ================================================================================

func (s *Scanner) match(exp string) bool {
	if s.isAtEnd() {
		return false
	}
	if string(s.source[s.current]) != exp {
		return false
	}
	s.current++
	return true
}
func (s *Scanner) matchElse(exp string, then, els TokenType) {
	if s.match(exp) {
		s.addToken(then, nil)
	} else {
		s.addToken(els, nil)
	}
}
func (s *Scanner) peek() string {
	if s.isAtEnd() {
		return "\000"
	}
	return string(s.source[s.current])
}
func (s *Scanner) peekNext() string {
	if s.current+1 >= len(s.source) {
		return "\000"
	}
	return string(s.source[s.current+1])
}
func (s *Scanner) advance() string {
	s.current++
	return string(s.source[s.current-1])
}
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || c == "_"
}
func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isDigit(c)
}
func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func (s *Scanner) error(line int, message string) {
	log.Printf("[line %d] Error: %s", line, message)
	s.hadErr = true
}
