package internal

import (
	"fmt"
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
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

type scanner struct {
	source   string
	tokens   []*token
	keywords map[string]TokenType

	start   int
	current int
	line    int
}

func newScanner(source string) *scanner {
	fmt.Println(source)

	return &scanner{
		source:   source,
		tokens:   make([]*token, 0),
		keywords: keywords,
		line:     1,
	}
}

func (s *scanner) scanTokens() []*token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, newToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *scanner) scanToken() {
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
		if s.match("=") {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case "=":
		if s.match("=") {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case "<":
		if s.match("=") {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case ">":
		if s.match("=") {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
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
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			Program.error(s.line, "Unexpected character.")
		}
	}
}

func (s *scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	ttype := s.keywords[text]

	if ttype == "" {
		ttype = IDENTIFIER
	}

	s.addToken(ttype, nil)
}

func (s *scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == "." && s.isDigit(s.peekNext()) {
		s.advance() // Consume the "."
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.Atoi(s.source[s.start:s.current])
	if err != nil {
		panic("Could not convert number to integer.")
	}
	s.addToken(NUMBER, value)
}

func (s *scanner) string() {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		Program.error(s.line, "Unterminated string.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}
	if string(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *scanner) peek() string {
	if s.isAtEnd() {
		return "\x00"
	}
	return string(s.source[s.current])
}

func (s *scanner) peekNext() string {
	if s.current+1 >= len(s.source) {
		return "\x00"
	}
	return string(s.source[s.current+1])
}

func (s *scanner) isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || c == "_"
}

func (s *scanner) isAlphaNumeric(c string) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *scanner) isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func (s *scanner) advance() string {
	c := string(s.source[s.current])
	s.current++
	return c
}

func (s *scanner) addToken(tt TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, newToken(tt, text, literal, s.line))
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
