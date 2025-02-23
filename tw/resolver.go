package tw

import "fmt"

type FunctionType string

const (
	FunctionNone        FunctionType = "none"
	FunctionFunction    FunctionType = "function"
	FunctionInitializer FunctionType = "initializer"
	FunctionMethod      FunctionType = "method"
)

type ClassType string

const (
	ClassNone     ClassType = "none"
	ClassClass    ClassType = "class"
	ClassSubclass ClassType = "subclass"
)

type Resolver struct {
	interpreter  *Interpreter
	scopes       Stack[map[string]bool]
	currentFn    FunctionType
	currentClass ClassType
	hadErr       bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:  interpreter,
		scopes:       Stack[map[string]bool]{},
		currentFn:    FunctionNone,
		currentClass: ClassNone,
		hadErr:       false,
	}
}

func (r *Resolver) Resolve(stmts []Stmt) bool {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
		if r.hadErr {
			return true
		}
	}
	return r.hadErr
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.Get(i)[name.lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(stmt *FunctionStmt, ftype FunctionType) {
	enclosingFn := r.currentFn
	r.currentFn = ftype

	r.beginScope()
	for _, param := range stmt.params {
		r.declare(param)
		r.define(param)
	}
	r.Resolve(stmt.body)
	r.endScope()

	r.currentFn = enclosingFn
}

// ================================================================================
// ### STMT VISITORS
// ================================================================================

func (r *Resolver) visitBlockStmt(stmt *BlockStmt) (StmtReturn, error) {
	r.beginScope()
	r.Resolve(stmt.stmts)
	r.endScope()
	return StmtReturn{}, nil
}

func (r *Resolver) visitClassStmt(stmt *ClassStmt) (StmtReturn, error) {
	enclosingClass := r.currentClass
	r.currentClass = ClassClass

	r.declare(stmt.name)
	r.define(stmt.name)

	if stmt.superclass != nil && stmt.name.lexeme == stmt.superclass.name.lexeme {
		r.error(stmt.superclass.name, "A class can't inherit from itself")
	}

	if stmt.superclass != nil {
		r.currentClass = ClassSubclass
		r.resolveExpr(stmt.superclass)
	}

	if stmt.superclass != nil {
		r.beginScope()
		r.scopes.Peek()["super"] = true
	}

	r.beginScope()
	r.scopes.Peek()["this"] = true
	for _, method := range stmt.methods {
		declaration := FunctionMethod
		if method.name.lexeme == "init" {
			declaration = FunctionInitializer
		}
		r.resolveFunction(method, declaration)
	}
	r.endScope()
	if stmt.superclass != nil {
		r.endScope()
	}
	r.currentClass = enclosingClass
	return StmtReturn{}, nil
}

func (r *Resolver) visitExpressionStmt(stmt *ExpressionStmt) (StmtReturn, error) {
	r.resolveExpr(stmt.expr)
	return StmtReturn{}, nil
}

func (r *Resolver) visitFunctionStmt(stmt *FunctionStmt) (StmtReturn, error) {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FunctionFunction)
	return StmtReturn{}, nil
}

func (r *Resolver) visitIfStmt(stmt *IfStmt) (StmtReturn, error) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
	return StmtReturn{}, nil
}

func (r *Resolver) visitPrintStmt(stmt *PrintStmt) (StmtReturn, error) {
	r.resolveExpr(stmt.expr)
	return StmtReturn{}, nil
}

func (r *Resolver) visitReturnStmt(stmt *ReturnStmt) (StmtReturn, error) {
	if r.currentFn == FunctionNone {
		r.error(stmt.keyword, "Cannot return from top-level code.")
	}
	if stmt.value != nil {
		if r.currentFn == FunctionInitializer {
			r.error(stmt.keyword, "Cannot return a value from an initializer.")
		}
		r.resolveExpr(stmt.value)
	}
	return StmtReturn{}, nil
}

func (r *Resolver) visitVarStmt(stmt *VarStmt) (StmtReturn, error) {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpr(stmt.initializer)
	}
	r.define(stmt.name)
	return StmtReturn{}, nil
}

func (r *Resolver) visitWhileStmt(stmt *WhileStmt) (StmtReturn, error) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
	return StmtReturn{}, nil
}

// ================================================================================
// ### EXPR VISITORS
// ================================================================================

func (r *Resolver) visitAssignExpr(expr *AssignExpr) (interface{}, error) {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil, nil
}

func (r *Resolver) visitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil, nil
}

func (r *Resolver) visitCallExpr(expr *CallExpr) (interface{}, error) {
	r.resolveExpr(expr.callee)
	for _, arg := range expr.args {
		r.resolveExpr(arg)
	}
	return nil, nil
}

func (r *Resolver) visitGetExpr(expr *GetExpr) (interface{}, error) {
	r.resolveExpr(expr.object)
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	r.resolveExpr(expr.expr)
	return nil, nil
}

func (r *Resolver) visitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(expr *LogicalExpr) (interface{}, error) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil, nil
}

func (r *Resolver) visitSetExpr(expr *SetExpr) (interface{}, error) {
	r.resolveExpr(expr.value)
	r.resolveExpr(expr.object)
	return nil, nil
}

func (r *Resolver) visitSuperExpr(expr *SuperExpr) (interface{}, error) {
	if r.currentClass == ClassNone {
		r.error(expr.keyword, "can't use 'super' outside of a class")
	} else if r.currentClass != ClassSubclass {
		r.error(expr.keyword, "can't use 'super' in a class with no subclass")
	}
	r.resolveLocal(expr, expr.keyword)
	return nil, nil
}

func (r *Resolver) visitThisExpr(expr *ThisExpr) (interface{}, error) {
	if r.currentClass == ClassNone {
		r.error(expr.keyword, "Cannot use 'this' outside of a class.")
		return nil, nil
	}

	r.resolveLocal(expr, expr.keyword)
	return nil, nil
}

func (r *Resolver) visitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	r.resolveExpr(expr.right)
	return nil, nil
}

func (r *Resolver) visitVariableExpr(expr *VariableExpr) (interface{}, error) {
	if !r.scopes.IsEmpty() {
		if val, ok := r.scopes.Peek()[expr.name.lexeme]; ok && !val {
			r.error(expr.name, "Cannot read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil, nil
}

// ================================================================================
// ### HELPERS
// ================================================================================

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name *Token) {
	if r.scopes.IsEmpty() {
		return
	}
	scope := r.scopes.Peek()
	if _, ok := scope[name.lexeme]; ok {
		r.error(name, "Variable with this name already declared in this scope.")
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if r.scopes.IsEmpty() {
		return
	}
	r.scopes.Peek()[name.lexeme] = true
}

func (r *Resolver) error(token *Token, msg string) {
	fmt.Println("[line", token.line, "] Error", msg)
	r.hadErr = true
}
