package stmt

import (
	"github.com/sjsanc/golox/internal/expr"
	"github.com/sjsanc/golox/internal/token"
)

type Stmt interface {
	Accept(v Visitor) interface{}
}

type Visitor interface {
	VisitBlockStmt(s Block) interface{}
	VisitExpressionStmt(s Expression) interface{}
	VisitIfStmt(s If) interface{}
	VisitPrintStmt(s Print) interface{}
	VisitVarStmt(s Var) interface{}
	VisitWhileStmt(s While) interface{}
}

type Block struct {
	Statements []Stmt
}

func (s Block) Accept(v Visitor) interface{} {
	return v.VisitBlockStmt(s)
}

type Expression struct {
	Expression expr.Expr
}

func (s Expression) Accept(v Visitor) interface{} {
	return v.VisitExpressionStmt(s)
}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (s If) Accept(v Visitor) interface{} {
	return v.VisitIfStmt(s)
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

type While struct {
	Condition expr.Expr
	Body      Stmt
}

func (s While) Accept(v Visitor) interface{} {
	return v.VisitWhileStmt(s)
}
