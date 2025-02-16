package internal

import "fmt"

var ErrParser = "Parser error"

type parser struct {
	tokens  []*token
	current int
}

func newParser(tokens []*token) *parser {
	return &parser{
		tokens: tokens,
	}
}

func (p *parser) parse() expr {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error?")
		}
	}()

	return p.expression()
}

func (p *parser) expression() expr {
	return p.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *parser) equality() expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = binaryExpr{expr, operator, right}
	}
	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *parser) comparison() expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = binaryExpr{expr, operator, right}
	}
	return expr
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *parser) term() expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = binaryExpr{expr, operator, right}
	}
	return expr
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *parser) factor() expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = binaryExpr{expr, operator, right}
	}
	return expr
}

// unary → ( "!" | "-" ) unary | primary ;
func (p *parser) unary() expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return unaryExpr{operator, right}
	}
	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *parser) primary() expr {
	if p.match(FALSE) {
		return literalExpr{false}
	}
	if p.match(TRUE) {
		return literalExpr{true}
	}
	if p.match(NIL) {
		return literalExpr{nil}
	}

	if p.match(NUMBER, STRING) {
		return literalExpr{p.previous().literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return groupingExpr{expr}
	}

	panic(p.error(p.peek(), "Expect expression."))
}

func (p *parser) match(ttypes ...TokenType) bool {
	for _, tt := range ttypes {
		if p.check(tt) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) consume(ttype TokenType, message string) *token {
	if p.check(ttype) {
		return p.advance()
	}
	panic(message)
}

func (p *parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().ttype == ttype
}

func (p *parser) advance() *token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().ttype == EOF
}

func (p *parser) peek() *token {
	return p.tokens[p.current]
}

func (p *parser) previous() *token {
	return p.tokens[p.current-1]
}

func (p *parser) error(token *token, message string) error {
	Program.errorToken(token, message)
	return fmt.Errorf(ErrParser)
}

func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().ttype == SEMICOLON {
			return
		}

		switch p.peek().ttype {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
		}

		p.advance()
	}
}
