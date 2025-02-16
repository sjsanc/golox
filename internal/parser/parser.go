package parser

import (
	"fmt"

	"github.com/sjsanc/golox/internal/expr"
	"github.com/sjsanc/golox/internal/stmt"
	"github.com/sjsanc/golox/internal/token"
)

type Parser struct {
	tokens   []*token.Token
	current  int
	hadError bool
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() ([]stmt.Stmt, bool) {
	stmts := make([]stmt.Stmt, 0)
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts, p.hadError
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

func (p *Parser) declaration() stmt.Stmt {
	if p.hadError {
		p.synchronize()
		return nil
	}

	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) statement() stmt.Stmt {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() stmt.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return stmt.Print{Expression: value}
}

func (p *Parser) varDeclaration() stmt.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")

	var initializer expr.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return stmt.Var{Name: name, Initializer: initializer}
}

func (p *Parser) expressionStatement() stmt.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return stmt.Expression{Expression: value}
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() expr.Expr {
	e := p.comparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		e = expr.Binary{Left: e, Operator: operator, Right: right}
	}
	return e
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() expr.Expr {
	e := p.term()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		e = expr.Binary{Left: e, Operator: operator, Right: right}
	}
	return e
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() expr.Expr {
	e := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		e = expr.Binary{Left: e, Operator: operator, Right: right}
	}
	return e
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() expr.Expr {
	e := p.unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		e = expr.Binary{Left: e, Operator: operator, Right: right}
	}
	return e
}

// unary → ( "!" | "-" ) unary
//
//	| primary ;
func (p *Parser) unary() expr.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return expr.Unary{Operator: operator, Right: right}
	}
	return p.primary()
}

// primary → NUMBER | STRING | "false" | "true" | "nil"
func (p *Parser) primary() expr.Expr {
	if p.match(token.FALSE) {
		return expr.Literal{Value: false}
	}
	if p.match(token.TRUE) {
		return expr.Literal{Value: true}
	}
	if p.match(token.NIL) {
		return expr.Literal{Value: nil}
	}
	if p.match(token.NUMBER, token.STRING) {
		return expr.Literal{Value: p.previous().Literal}
	}
	if p.match(token.IDENTIFIER) {
		return expr.Variable{Name: p.previous()}
	}

	if p.match(token.LEFT_PAREN) {
		e := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return expr.Grouping{Expr: e}
	}

	p.err(p.peek().Line, "Expect expression.")
	return nil
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t token.TokenType, message string) *token.Token {
	if p.check(t) {
		return p.advance()
	}
	panic(message)
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func (s *Parser) err(line int, message string) {
	fmt.Println("Error: ", line, message)
	s.hadError = true
}
