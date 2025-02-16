package scanner

import (
	"fmt"
	"strconv"

	"github.com/sjsanc/golox/internal/token"
)

var keywords = map[string]token.TokenType{
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type Scanner struct {
	source   string
	tokens   []*token.Token
	keywords map[string]token.TokenType
	start    int
	current  int
	line     int
	hadError bool
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		source:   src,
		tokens:   make([]*token.Token, 0),
		keywords: keywords,
		line:     1,
	}
}

func (s *Scanner) ScanTokens() ([]*token.Token, bool) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, s.line))
	return s.tokens, s.hadError
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case "(":
		s.addToken(token.LEFT_PAREN, nil)
	case ")":
		s.addToken(token.RIGHT_PAREN, nil)
	case "{":
		s.addToken(token.LEFT_BRACE, nil)
	case "}":
		s.addToken(token.RIGHT_BRACE, nil)
	case ",":
		s.addToken(token.COMMA, nil)
	case ".":
		s.addToken(token.DOT, nil)
	case "-":
		s.addToken(token.MINUS, nil)
	case "+":
		s.addToken(token.PLUS, nil)
	case ";":
		s.addToken(token.SEMICOLON, nil)
	case "*":
		s.addToken(token.STAR, nil)
	case "!":
		if s.match("=") {
			s.addToken(token.BANG_EQUAL, nil)
		} else {
			s.addToken(token.BANG, nil)
		}
	case "=":
		if s.match("=") {
			s.addToken(token.EQUAL_EQUAL, nil)
		} else {
			s.addToken(token.EQUAL, nil)
		}
	case "<":
		if s.match("=") {
			s.addToken(token.LESS_EQUAL, nil)
		} else {
			s.addToken(token.LESS, nil)
		}
	case ">":
		if s.match("=") {
			s.addToken(token.GREATER_EQUAL, nil)
		} else {
			s.addToken(token.GREATER, nil)
		}
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
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
			fmt.Println("Unexpected character")
			s.hadError = true
		}
	}
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tt, ok := s.keywords[text]
	if !ok {
		tt = token.IDENTIFIER
	}
	s.addToken(tt, nil)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == "." && isDigit(s.peekNext()) {
		s.advance() // consume the "."
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.Atoi(s.source[s.start:s.current])
	if err != nil {
		s.err(s.line, "Error parsing number")
		return
	}
	s.addToken(token.NUMBER, value)
}

func (s *Scanner) string() {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.err(s.line, "Unterminated string.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

func (s *Scanner) addToken(tt token.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.New(tt, text, literal, s.line))
}
func (s *Scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}
	if string(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
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

func (s *Scanner) err(line int, message string) {
	fmt.Println("Error: ", line, message)
	s.hadError = true
}
