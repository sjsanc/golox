package resolver

import (
	"fmt"

	"github.com/sjsanc/golox/internal/expr"
	"github.com/sjsanc/golox/internal/interpreter"
	"github.com/sjsanc/golox/internal/stmt"
	"github.com/sjsanc/golox/internal/token"
)

type Stack[T any] struct {
	values []T
}

func (s *Stack[T]) get(index int) T {
	return s.values[index]
}

func (s *Stack[T]) push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) pop() T {
	value := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return value
}

func (s *Stack[T]) peek() T {
	return s.values[len(s.values)-1]
}

func (s *Stack[T]) isEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) size() int {
	return len(s.values)
}

type FunctionType string

const (
	FunctionNone     FunctionType = "none"
	FunctionFunction FunctionType = "function"
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          Stack[map[string]bool]
	currentFunction FunctionType
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          Stack[map[string]bool]{},
		currentFunction: FunctionNone,
	}
}

// ================================================================================
// ### STMT VISITORS
// ================================================================================

func (r *Resolver) VisitBlockStmt(s stmt.Block) (stmt.ReturnValue, error) {
	r.beginScope()
	r.Resolve(s.Statements)
	r.endScope()
	return stmt.ReturnValue{}, nil
}

func (r *Resolver) VisitVarStmt(s stmt.Var) (stmt.ReturnValue, error) {
	err := r.declare(s.Name)
	if err != nil {
		return stmt.ReturnValue{}, err
	}

	if s.Initializer != nil {
		r.resolveExpr(s.Initializer)
	}
	r.define(s.Name)
	return stmt.ReturnValue{}, nil
}

func (r *Resolver) VisitFunctionStmt(s stmt.Function) (stmt.ReturnValue, error) {
	err := r.declare(s.Name)
	if err != nil {
		return stmt.ReturnValue{}, err
	}
	r.define(s.Name)

	err = r.resolveFunction(s, FunctionFunction)
	if err != nil {
		return stmt.ReturnValue{}, err
	}

	return stmt.ReturnValue{}, nil
}

func (r *Resolver) VisitExpressionStmt(s stmt.Expression) (stmt.ReturnValue, error) {
	r.resolveExpr(s.Expression)
	return stmt.ReturnValue{}, nil
}

func (r *Resolver) VisitIfStmt(s stmt.If) (stmt.ReturnValue, error) {
	r.resolveExpr(s.Condition)
	r.resolveStmt(s.ThenBranch)
	if s.ElseBranch != nil {
		r.resolveStmt(s.ElseBranch)
	}
	return stmt.ReturnValue{}, nil
}

func (r *Resolver) VisitPrintStmt(s stmt.Print) (stmt.ReturnValue, error) {
	r.resolveExpr(s.Expression)
	return stmt.ReturnValue{}, nil
}

func (r *Resolver) VisitReturnStmt(s stmt.Return) (stmt.ReturnValue, error) {
	if r.currentFunction == FunctionNone {
		return stmt.ReturnValue{}, fmt.Errorf("can't return from top-level code")
	}
	if s.Value != nil {
		r.resolveExpr(s.Value)
	}
	return stmt.ReturnValue{}, nil
}

func (r *Resolver) VisitWhileStmt(s stmt.While) (stmt.ReturnValue, error) {
	r.resolveExpr(s.Condition)
	r.resolveStmt(s.Body)
	return stmt.ReturnValue{}, nil
}

// ================================================================================
// ### EXPR VISITORS
// ================================================================================

func (r *Resolver) VisitVariableExpr(e expr.Variable) (interface{}, error) {
	_, ok := r.scopes.peek()[e.Name.Lexeme]
	if !r.scopes.isEmpty() && !ok {
		return nil, fmt.Errorf("can't read local variable in its own initializer")
	}
	r.resolveLocal(e, e.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(e expr.Assign) (interface{}, error) {
	r.resolveExpr(e.Value)
	r.resolveLocal(e, e.Name)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(e expr.Binary) (interface{}, error) {
	r.resolveExpr(e.Left)
	r.resolveExpr(e.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(e expr.Call) (interface{}, error) {
	r.resolveExpr(e.Callee)
	for _, argument := range e.Args {
		r.resolveExpr(argument)
	}
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(e expr.Grouping) (interface{}, error) {
	r.resolveExpr(e.Expr)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(e expr.Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(e expr.Logical) (interface{}, error) {
	r.resolveExpr(e.Left)
	r.resolveExpr(e.Right)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(e expr.Unary) (interface{}, error) {
	r.resolveExpr(e.Right)
	return nil, nil
}

// ================================================================================
// ### HELPERS
// ================================================================================

func (r *Resolver) Resolve(statements []stmt.Stmt) error {
	for _, statement := range statements {
		err := r.resolveStmt(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(statement stmt.Stmt) error {
	_, err := statement.Accept(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) resolveExpr(expression expr.Expr) error {
	_, err := expression.Accept(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) resolveLocal(e expr.Expr, name *token.Token) {
	for i := r.scopes.size() - 1; i >= 0; i-- {
		if r.scopes.get(i)[name.Lexeme] {
			r.interpreter.Resolve(e, r.scopes.size()-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function stmt.Function, ftype FunctionType) error {
	enclosingFn := r.currentFunction
	r.currentFunction = ftype

	r.beginScope()
	for _, param := range function.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		r.define(param)
	}
	r.Resolve(function.Body)
	r.endScope()

	r.currentFunction = enclosingFn
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes.pop()
}

func (r *Resolver) declare(name *token.Token) error {
	if r.scopes.isEmpty() {
		return nil
	}
	scope := r.scopes.peek()
	if _, ok := scope[name.Lexeme]; ok {
		return fmt.Errorf("already a variable with this name in this scope")
	}
	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name *token.Token) {
	if r.scopes.isEmpty() {
		return
	}
	scope := r.scopes.peek()
	scope[name.Lexeme] = true
}
