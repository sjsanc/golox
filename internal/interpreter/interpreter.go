package interpreter

import (
	"fmt"

	"github.com/sjsanc/golox/internal/environment"
	"github.com/sjsanc/golox/internal/errors"
	"github.com/sjsanc/golox/internal/expr"
	"github.com/sjsanc/golox/internal/stmt"
	"github.com/sjsanc/golox/internal/token"
)

type Interpreter struct {
	environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: environment.NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(stmts []stmt.Stmt) {
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *Interpreter) VisitLiteralExpr(expr expr.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr expr.Grouping) interface{} {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitUnaryExpr(expr expr.Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		checkNumOperand(expr.Operator, right)
		return -(right.(int))
	case token.BANG:
		return !isTruthy(right)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr expr.Variable) interface{} {
	value, err := i.environment.Get(expr.Name)
	if err != nil {
		panic(err)
	}
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr expr.Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) - right.(int)
	case token.SLASH:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) / right.(int)
	case token.STAR:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) * right.(int)
	case token.PLUS:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l + r
			}
		}
	case token.GREATER:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) > right.(int)
	case token.GREATER_EQUAL:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) >= right.(int)
	case token.LESS:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) < right.(int)
	case token.LESS_EQUAL:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) <= right.(int)
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt stmt.Expression) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt stmt.Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Printf("%v\n", value)
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt stmt.Var) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) evaluate(expr expr.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt stmt.Stmt) {
	stmt.Accept(i)
}

func checkNumOperand(operator *token.Token, operand interface{}) {
	if _, ok := operand.(int); !ok {
		panic(errors.RuntimeErr{Token: operator, Message: "Operand must be a number"})
	}
}

func checkNumOperands(operator *token.Token, left, right interface{}) {
	if _, ok := left.(int); !ok {
		panic(errors.RuntimeErr{Token: operator, Message: "Left operand must be a number"})
	}
	if _, ok := right.(int); !ok {
		panic(errors.RuntimeErr{Token: operator, Message: "Right operand must be a number"})
	}
}

func isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
