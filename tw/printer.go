package tw

import (
	"fmt"
	"strings"
)

type Printer struct {
}

func (p *Printer) PrintExpr(expr Expr) string {
	v, _ := expr.Accept(p)
	return fmt.Sprintf("%v", v)
}

func (p *Printer) PrintStmt(stmt Stmt) string {
	v, _ := stmt.Accept(p)
	return fmt.Sprintf("%v", v)
}

// ================================================================================
// ### STATEMENTS
// ================================================================================

func (p *Printer) visitBlockStmt(stmt *BlockStmt) (StmtReturn, error) {
	sb := new(strings.Builder)
	sb.WriteString("(block ")
	for _, s := range stmt.stmts {
		v, _ := s.Accept(p)
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString(")")
	return StmtReturn{value: sb.String()}, nil
}

func (p *Printer) visitClassStmt(stmt *ClassStmt) (StmtReturn, error) {
	sb := new(strings.Builder)
	sb.WriteString("(class ")
	sb.WriteString(stmt.name.lexeme)
	for _, m := range stmt.methods {
		v, _ := m.Accept(p)
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString(")")
	return StmtReturn{value: sb.String()}, nil
}

func (p *Printer) visitExpressionStmt(stmt *ExpressionStmt) (StmtReturn, error) {
	return StmtReturn{value: p.parenthesize(";", stmt.expr)}, nil
}

func (p *Printer) visitFunctionStmt(stmt *FunctionStmt) (StmtReturn, error) {
	sb := new(strings.Builder)
	sb.WriteString("(fun " + stmt.name.lexeme + "(")
	for _, param := range stmt.params {
		if param != stmt.params[0] {
			sb.WriteString(" ")
		}
		sb.WriteString(param.lexeme)
	}
	sb.WriteString(") ")
	for _, s := range stmt.body {
		v, _ := s.Accept(p)
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString(")")
	return StmtReturn{value: sb.String()}, nil
}

func (p *Printer) visitIfStmt(stmt *IfStmt) (StmtReturn, error) {
	if stmt.elseBranch != nil {
		return StmtReturn{value: p.parenthesize("if-else", stmt.condition, stmt.thenBranch, stmt.elseBranch)}, nil
	}
	return StmtReturn{value: p.parenthesize("if", stmt.condition, stmt.thenBranch)}, nil
}

func (p *Printer) visitPrintStmt(stmt *PrintStmt) (StmtReturn, error) {
	return StmtReturn{value: p.parenthesize("print", stmt.expr)}, nil
}

func (p *Printer) visitReturnStmt(stmt *ReturnStmt) (StmtReturn, error) {
	if stmt.value != nil {
		return StmtReturn{value: "(return)"}, nil
	}
	return StmtReturn{value: p.parenthesize("return", stmt.value)}, nil
}

func (p *Printer) visitVarStmt(stmt *VarStmt) (StmtReturn, error) {
	if stmt.initializer != nil {
		return StmtReturn{value: p.parenthesize("var", stmt.name, "=", stmt.initializer)}, nil
	}
	return StmtReturn{value: p.parenthesize("var", stmt.name)}, nil
}

func (p *Printer) visitWhileStmt(stmt *WhileStmt) (StmtReturn, error) {
	return StmtReturn{value: p.parenthesize("while", stmt.condition, stmt.body)}, nil
}

// ================================================================================
// ### EXPRESSIONS
// ================================================================================

func (p *Printer) visitAssignExpr(expr *AssignExpr) (interface{}, error) {
	return p.parenthesize("=", expr.name.lexeme, expr.value), nil
}

func (p *Printer) visitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right), nil
}

func (p *Printer) visitCallExpr(expr *CallExpr) (interface{}, error) {
	return p.parenthesize("call", expr.callee, expr.paren, expr.args), nil
}

func (p *Printer) visitGetExpr(expr *GetExpr) (interface{}, error) {
	return p.parenthesize(".", expr.object, expr.name), nil
}

func (p *Printer) visitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	return p.parenthesize("group", expr.expr), nil
}

func (p *Printer) visitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	if expr.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.value), nil
}

func (p *Printer) visitLogicalExpr(expr *LogicalExpr) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right), nil
}

func (p *Printer) visitSetExpr(expr *SetExpr) (interface{}, error) {
	return p.parenthesize("=", expr.object, expr.name, expr.value), nil
}

func (p *Printer) visitThisExpr(expr *ThisExpr) (interface{}, error) {
	return "this", nil
}

func (p *Printer) visitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.right), nil
}

func (p *Printer) visitVariableExpr(expr *VariableExpr) (interface{}, error) {
	return expr.name.lexeme, nil
}

// ================================================================================
// ### HELPERS
// ================================================================================

func (p *Printer) parenthesize(name string, parts ...interface{}) string {
	sb := new(strings.Builder)
	sb.WriteString("(" + name)
	p.transform(sb, parts)
	sb.WriteString(")")
	return sb.String()
}

func (p *Printer) transform(sb *strings.Builder, parts []interface{}) {
	for _, part := range parts {
		sb.WriteString(" ")
		switch v := part.(type) {
		case Expr:
			res, _ := v.Accept(p)
			sb.WriteString(fmt.Sprintf("%v", res))
		case Stmt:
			res, _ := v.Accept(p)
			sb.WriteString(fmt.Sprintf("%v", res))
		case Token:
			sb.WriteString(v.lexeme)
		case []interface{}:
			p.transform(sb, v)
		default:
			sb.WriteString(fmt.Sprintf("%v", v))
		}
	}
}
