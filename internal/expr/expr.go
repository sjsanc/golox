package expr

import (
	"github.com/sjsanc/golox/internal/token"
)

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type Visitor interface {
	VisitAssignExpr(e Assign) (interface{}, error)
	VisitBinaryExpr(e Binary) (interface{}, error)
	VisitCallExpr(e Call) (interface{}, error)
	VisitGroupingExpr(e Grouping) (interface{}, error)
	VisitLiteralExpr(e Literal) (interface{}, error)
	VisitLogicalExpr(e Logical) (interface{}, error)
	VisitUnaryExpr(e Unary) (interface{}, error)
	VisitVariableExpr(e Variable) (interface{}, error)
}

type Assign struct {
	Name  *token.Token
	Value Expr
}

func (e Assign) Accept(v Visitor) (interface{}, error) {
	return v.VisitAssignExpr(e)
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (e Binary) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinaryExpr(e)
}

type Call struct {
	Callee Expr
	Paren  *token.Token
	Args   []Expr
}

func (e Call) Accept(v Visitor) (interface{}, error) {
	return v.VisitCallExpr(e)
}

type Grouping struct {
	Expr Expr
}

func (e Grouping) Accept(v Visitor) (interface{}, error) {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value interface{}
}

func (e Literal) Accept(v Visitor) (interface{}, error) {
	return v.VisitLiteralExpr(e)
}

type Logical struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (e Logical) Accept(v Visitor) (interface{}, error) {
	return v.VisitLogicalExpr(e)
}

type Unary struct {
	Operator *token.Token
	Right    Expr
}

func (e Unary) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnaryExpr(e)
}

type Variable struct {
	Name *token.Token
}

func (e Variable) Accept(v Visitor) (interface{}, error) {
	return v.VisitVariableExpr(e)
}
