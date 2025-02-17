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
	return p.assignment()
}

func (p *Parser) declaration() stmt.Stmt {
	if p.hadError {
		p.synchronize()
		return nil
	}

	if p.match(token.FUN) {
		return p.function("function")
	}

	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) statement() stmt.Stmt {
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.RETURN) {
		return p.returnStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return stmt.Block{Statements: p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() stmt.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer stmt.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition expr.Expr
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment expr.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = stmt.Block{Statements: []stmt.Stmt{body, stmt.Expression{Expression: increment}}}
	}

	if condition == nil {
		condition = expr.Literal{Value: true}
	}

	body = stmt.While{Condition: condition, Body: body}

	if initializer != nil {
		body = stmt.Block{Statements: []stmt.Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) ifStatement() stmt.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch stmt.Stmt
	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return stmt.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) printStatement() stmt.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return stmt.Print{Expression: value}
}

func (p *Parser) returnStatement() stmt.Stmt {
	keyword := p.previous()
	var value expr.Expr
	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after return value.")
	return stmt.Return{Keyword: keyword, Value: value}
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

func (p *Parser) whileStatement() stmt.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return stmt.While{Condition: condition, Body: body}
}

func (p *Parser) expressionStatement() stmt.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return stmt.Expression{Expression: value}
}

func (p *Parser) function(kind string) stmt.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	params := make([]*token.Token, 0)
	if !p.check(token.RIGHT_PAREN) {
		params = append(params, p.consume(token.IDENTIFIER, "Expect parameter name."))
		for p.match(token.COMMA) {
			if len(params) >= 255 {
				p.error(p.peek().Line, "Cannot have more than 255 parameters.")
			}
			params = append(params, p.consume(token.IDENTIFIER, "Expect parameter name."))
		}
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.block()
	return stmt.Function{Name: name, Params: params, Body: body}
}

func (p *Parser) block() []stmt.Stmt {
	stmts := make([]stmt.Stmt, 0)

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	return stmts
}

func (p *Parser) assignment() expr.Expr {
	e := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if e, ok := e.(expr.Variable); ok {
			name := e.Name
			return expr.Assign{Name: name, Value: value}
		}

		p.error(equals.Line, "Invalid assignment target.")
	}

	return e
}

func (p *Parser) or() expr.Expr {
	e := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		e = expr.Logical{Left: e, Operator: operator, Right: right}
	}

	return e
}

func (p *Parser) and() expr.Expr {
	e := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		e = expr.Logical{Left: e, Operator: operator, Right: right}
	}

	return e
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
	return p.call()
}

func (p *Parser) finishCall(callee expr.Expr) expr.Expr {
	arguments := make([]expr.Expr, 0)
	if !p.check(token.RIGHT_PAREN) {
		arguments = append(arguments, p.expression())
		for p.match(token.COMMA) {
			if len(arguments) >= 255 {
				p.error(p.peek().Line, "Cannot have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
		}
	}
	paren := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	return expr.Call{Callee: callee, Paren: paren, Args: arguments}
}

func (p *Parser) call() expr.Expr {
	e := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			e = p.finishCall(e)
		} else {
			break
		}
	}

	return e
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

	p.error(p.peek().Line, "Expect expression.")
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

func (s *Parser) error(line int, message string) {
	fmt.Println("Error: ", line, message)
	s.hadError = true
}
