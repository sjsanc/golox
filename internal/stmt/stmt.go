package stmt

import (
	"github.com/sjsanc/golox/internal/expr"
	"github.com/sjsanc/golox/internal/token"
)

type Stmt interface {
	Accept(v Visitor) interface{}
}

type Visitor interface {
	VisitExpressionStmt(s Expression) interface{}
	VisitPrintStmt(s Print) interface{}
	VisitVarStmt(s Var) interface{}
}

type Expression struct {
	Expression expr.Expr
}

func (s Expression) Accept(v Visitor) interface{} {
	return v.VisitExpressionStmt(s)
}

type Print struct {
	Expression expr.Expr
}

func (s Print) Accept(v Visitor) interface{} {
	return v.VisitPrintStmt(s)
}

type Var struct {
	Name        *token.Token
	Initializer expr.Expr
}

func (s Var) Accept(v Visitor) interface{} {
	return v.VisitVarStmt(s)
}
