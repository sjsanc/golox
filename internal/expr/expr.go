package expr

import (
	"github.com/sjsanc/golox/internal/token"
)

type Expr interface {
	Accept(v Visitor) interface{}
}

type Visitor interface {
	VisitBinaryExpr(e Binary) interface{}
	VisitGroupingExpr(e Grouping) interface{}
	VisitLiteralExpr(e Literal) interface{}
	VisitUnaryExpr(e Unary) interface{}
	VisitVariableExpr(e Variable) interface{}
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
