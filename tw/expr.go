package tw

type Expr interface {
	Accept(v ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	visitAssignExpr(expr AssignExpr) (interface{}, error)
	visitBinaryExpr(expr BinaryExpr) (interface{}, error)
	visitCallExpr(expr CallExpr) (interface{}, error)
	visitGroupingExpr(expr GroupingExpr) (interface{}, error)
	visitLiteralExpr(expr LiteralExpr) (interface{}, error)
	visitLogicalExpr(expr LogicalExpr) (interface{}, error)
	visitUnaryExpr(expr UnaryExpr) (interface{}, error)
	visitVariableExpr(expr VariableExpr) (interface{}, error)
}

// ================================================================================
// ### ASSIGN
// ================================================================================

type AssignExpr struct {
	name  *Token
	value Expr
}

func (expr AssignExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitAssignExpr(expr)
}

// ================================================================================
// ### BINARY
// ================================================================================

type BinaryExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

func (expr BinaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitBinaryExpr(expr)
}

// ================================================================================
// ### CALL
// ================================================================================

type CallExpr struct {
	callee Expr
	paren  *Token
	args   []Expr
}

func (expr CallExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitCallExpr(expr)
}

// ================================================================================
// ### GROUPING
// ================================================================================

type GroupingExpr struct {
	expr Expr
}

func (expr GroupingExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitGroupingExpr(expr)
}

// ================================================================================
// ### LITERAL
// ================================================================================

type LiteralExpr struct {
	value interface{}
}

func (expr LiteralExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitLiteralExpr(expr)
}

// ================================================================================
// ### LOGICAL
// ================================================================================

type LogicalExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

func (expr LogicalExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitLogicalExpr(expr)
}

// ================================================================================
// ### UNARY
// ================================================================================

type UnaryExpr struct {
	operator *Token
	right    Expr
}

func (expr UnaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitUnaryExpr(expr)
}

// ================================================================================
// ### VARIABLE
// ================================================================================

type VariableExpr struct {
	name *Token
}

func (expr VariableExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.visitVariableExpr(expr)
}
