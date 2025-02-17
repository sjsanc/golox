package stmt

import (
	"github.com/sjsanc/golox/internal/expr"
	"github.com/sjsanc/golox/internal/token"
)

type ReturnValue struct {
	Value    interface{}
	IsReturn bool
}

type Stmt interface {
	Accept(v Visitor) (ReturnValue, error)
}

type Visitor interface {
	VisitBlockStmt(s Block) (ReturnValue, error)
	VisitExpressionStmt(s Expression) (ReturnValue, error)
	VisitFunctionStmt(s Function) (ReturnValue, error)
	VisitIfStmt(s If) (ReturnValue, error)
	VisitPrintStmt(s Print) (ReturnValue, error)
	VisitReturnStmt(s Return) (ReturnValue, error)
	VisitVarStmt(s Var) (ReturnValue, error)
	VisitWhileStmt(s While) (ReturnValue, error)
}

type Block struct {
	Statements []Stmt
}

func (s Block) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitBlockStmt(s)
}

type Expression struct {
	Expression expr.Expr
}

func (s Expression) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitExpressionStmt(s)
}

type Function struct {
	Name   *token.Token
	Params []*token.Token
	Body   []Stmt
}

func (s Function) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitFunctionStmt(s)
}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (s If) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitIfStmt(s)
}

type Print struct {
	Expression expr.Expr
}

func (s Print) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitPrintStmt(s)
}

type Return struct {
	Keyword *token.Token
	Value   expr.Expr
}

func (s Return) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitReturnStmt(s)
}

type Var struct {
	Name        *token.Token
	Initializer expr.Expr
}

func (s Var) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitVarStmt(s)
}

type While struct {
	Condition expr.Expr
	Body      Stmt
}

func (s While) Accept(v Visitor) (ReturnValue, error) {
	return v.VisitWhileStmt(s)
}
