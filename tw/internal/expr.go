package internal

type expr interface {
	accept(v visitor) interface{}
}

type visitor interface {
	visitBinaryExpr(e binaryExpr) interface{}
	visitGroupingExpr(e groupingExpr) interface{}
	visitLiteralExpr(e literalExpr) interface{}
	visitUnaryExpr(e unaryExpr) interface{}
}

type binaryExpr struct {
	left     expr
	operator token
	right    expr
}

func (e binaryExpr) accept(v visitor) interface{} {
	return v.visitBinaryExpr(e)
}

type groupingExpr struct {
	expr expr
}

func (e groupingExpr) accept(v visitor) interface{} {
	return v.visitGroupingExpr(e)
}

type literalExpr struct {
	value interface{}
}

func (e literalExpr) accept(v visitor) interface{} {
	return v.visitLiteralExpr(e)
}

type unaryExpr struct {
	operator token
	right    expr
}

func (e unaryExpr) accept(v visitor) interface{} {
	return v.visitUnaryExpr(e)
}
