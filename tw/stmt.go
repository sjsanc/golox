package tw

type StmtReturn struct {
	value    interface{}
	isReturn bool
}

type Stmt interface {
	Accept(v StmtVisitor) (StmtReturn, error)
}

type StmtVisitor interface {
	visitBlockStmt(stmt *BlockStmt) (StmtReturn, error)
	visitClassStmt(stmt *ClassStmt) (StmtReturn, error)
	visitExpressionStmt(stmt *ExpressionStmt) (StmtReturn, error)
	visitFunctionStmt(stmt *FunctionStmt) (StmtReturn, error)
	visitIfStmt(stmt *IfStmt) (StmtReturn, error)
	visitPrintStmt(stmt *PrintStmt) (StmtReturn, error)
	visitReturnStmt(stmt *ReturnStmt) (StmtReturn, error)
	visitVarStmt(stmt *VarStmt) (StmtReturn, error)
	visitWhileStmt(stmt *WhileStmt) (StmtReturn, error)
}

// ================================================================================
// ### BLOCK
// ================================================================================

type BlockStmt struct {
	stmts []Stmt
}

func (stmt *BlockStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitBlockStmt(stmt)
}

// ================================================================================
// ### CLASS
// ================================================================================

type ClassStmt struct {
	name    *Token
	methods []*FunctionStmt
}

func (stmt *ClassStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitClassStmt(stmt)
}

// ================================================================================
// ### EXPRESSION
// ================================================================================

type ExpressionStmt struct {
	expr Expr
}

func (stmt *ExpressionStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitExpressionStmt(stmt)
}

// ================================================================================
// ### FUNCTION
// ================================================================================

type FunctionStmt struct {
	name   *Token
	params []*Token
	body   []Stmt
}

func (stmt *FunctionStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitFunctionStmt(stmt)
}

// ================================================================================
// ### IF
// ================================================================================

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (stmt *IfStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitIfStmt(stmt)
}

// ================================================================================
// ### PRINT
// ================================================================================

type PrintStmt struct {
	expr Expr
}

func (stmt *PrintStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitPrintStmt(stmt)
}

// ================================================================================
// ### RETURN
// ================================================================================

type ReturnStmt struct {
	keyword *Token
	value   Expr
}

func (stmt *ReturnStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitReturnStmt(stmt)
}

// ================================================================================
// ### VAR
// ================================================================================

type VarStmt struct {
	name        *Token
	initializer Expr
}

func (stmt *VarStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitVarStmt(stmt)
}

// ================================================================================
// ### WHILE
// ================================================================================

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (stmt *WhileStmt) Accept(v StmtVisitor) (StmtReturn, error) {
	return v.visitWhileStmt(stmt)
}
