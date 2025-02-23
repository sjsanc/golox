package tw

import (
	"fmt"
)

type Parser struct {
	tokens  []*Token
	current int
	hadErr  bool
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() ([]Stmt, bool) {
	stmts := make([]Stmt, 0)
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts, p.hadErr
}

// ================================================================================
// ### STATEMENTS
// ================================================================================

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return &BlockStmt{stmts: p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}
	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")
	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		body = &BlockStmt{stmts: []Stmt{body, &ExpressionStmt{expr: increment}}}
	}
	if condition == nil {
		condition = &LiteralExpr{value: true}
	}
	body = &WhileStmt{condition: condition, body: body}
	if initializer != nil {
		body = &BlockStmt{stmts: []Stmt{initializer, body}}
	}
	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return &IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &PrintStmt{expr: value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value.")
	return &ReturnStmt{keyword: keyword, value: value}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after while condition.")
	body := p.statement()
	return &WhileStmt{condition: condition, body: body}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return &ExpressionStmt{expr: expr}
}

func (p *Parser) block() []Stmt {
	stmts := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return stmts
}

func (p *Parser) declaration() Stmt {
	if p.hadErr {
		p.synchronize()
		return nil
	}
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")
	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return &VarStmt{name: name, initializer: initializer}
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect class name.")
	p.consume(LEFT_BRACE, "Expect '{' before class body.")
	methods := make([]*FunctionStmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method").(*FunctionStmt))
	}
	p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	return &ClassStmt{name: name, methods: methods}
}

func (p *Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, "Expect "+kind+" name.")
	p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	params := make([]*Token, 0)
	if !p.check(RIGHT_PAREN) {
		if len(params) >= 255 {
			p.error(p.peek(), "Can't have more than 255 parameters.")
		}
		params = append(params, p.consume(IDENTIFIER, "Expect parameter name."))

		for p.match(COMMA) {
			if len(params) >= 255 {
				p.error(p.peek(), "Can't have more than 255 parameters.")
			}
			params = append(params, p.consume(IDENTIFIER, "Expect parameter name."))
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.block()
	return &FunctionStmt{name: name, params: params, body: body}
}

// ================================================================================
// ### EXPRESSIONS
// ================================================================================

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if expr, ok := expr.(*VariableExpr); ok {
			return &AssignExpr{name: expr.name, value: value}
		}
		if expr, ok := expr.(*GetExpr); ok {
			return &SetExpr{object: expr.object, name: expr.name, value: value}
		}
		p.error(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicalExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &LogicalExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return &UnaryExpr{operator: operator, right: right}
	}
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = &GetExpr{object: expr, name: name}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	args := make([]Expr, 0)
	if !p.check(RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(COMMA) {
			if len(args) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			args = append(args, p.expression())
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	return &CallExpr{callee: callee, paren: paren, args: args}
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &LiteralExpr{value: false}
	}
	if p.match(TRUE) {
		return &LiteralExpr{value: true}
	}
	if p.match(NIL) {
		return &LiteralExpr{value: nil}
	}
	if p.match(NUMBER, STRING) {
		return &LiteralExpr{value: p.previous().literal}
	}
	if p.match(THIS) {
		return &ThisExpr{keyword: p.previous()}
	}
	if p.match(IDENTIFIER) {
		return &VariableExpr{name: p.previous()}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{expr: expr}
	}

	p.error(p.peek(), "expected an expression. Last token was: "+p.peek().lexeme)
	return nil
}

// ================================================================================
// ### HELPERS
// ================================================================================

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}
func (p *Parser) consume(ttype TokenType, message string) *Token {
	if p.check(ttype) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return nil
}
func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().ttype == ttype
}
func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}
func (p *Parser) isAtEnd() bool {
	return p.peek().ttype == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().ttype == SEMICOLON {
			return
		}
		switch p.peek().ttype {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) error(token *Token, msg string) {
	fmt.Printf("[line %d] error: %s\n", token.line, msg)
	p.hadErr = true
}
