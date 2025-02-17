package expr

import (
	"github.com/sjsanc/golox/internal/token"
)

type Expr interface {
	Accept(v Visitor) interface{}
}

type Visitor interface {
	VisitAssignExpr(e Assign) interface{}
	VisitBinaryExpr(e Binary) interface{}
	VisitGroupingExpr(e Grouping) interface{}
	VisitLiteralExpr(e Literal) interface{}
	VisitLogicalExpr(e Logical) interface{}
	VisitUnaryExpr(e Unary) interface{}
	VisitVariableExpr(e Variable) interface{}
}

type Assign struct {
	Name  *token.Token
	Value Expr
}

func (e Assign) Accept(v Visitor) interface{} {
	return v.VisitAssignExpr(e)
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (e Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(e)
}

type Grouping struct {
	Expr Expr
}

func (e Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value interface{}
}

func (e Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(e)
}

type Logical struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (e Logical) Accept(v Visitor) interface{} {
	return v.VisitLogicalExpr(e)
}

type Unary struct {
	Operator *token.Token
	Right    Expr
}

func (e Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(e)
}

type Variable struct {
	Name *token.Token
}

func (e Variable) Accept(v Visitor) interface{} {
	return v.VisitVariableExpr(e)
}
