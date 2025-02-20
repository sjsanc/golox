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
	globals     *environment.Environment
	environment *environment.Environment
	locals      map[expr.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := environment.NewGlobalEnvironment()

	globals.Define("clock", &ClockBuiltin{})

	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[expr.Expr]int),
	}
}

func (i *Interpreter) Interpret(stmts []stmt.Stmt) error {
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (i *Interpreter) evaluate(expr expr.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(s stmt.Stmt) (stmt.ReturnValue, error) {
	result, err := s.Accept(i)
	if err != nil {
		return stmt.ReturnValue{}, err
	}
	return result, nil
}

func (i *Interpreter) executeBlock(stmts []stmt.Stmt, env *environment.Environment) (stmt.ReturnValue, error) {
	prev := i.environment
	i.environment = env
	defer func() { i.environment = prev }()

	for _, s := range stmts {
		result, err := i.execute(s)
		if err != nil {
			return stmt.ReturnValue{}, err
		}
		if result.IsReturn {
			return result, nil
		}
	}
	return stmt.ReturnValue{}, nil
}

func (i *Interpreter) Resolve(e expr.Expr, depth int) {
	i.locals[e] = depth
}

// ================================================================================
// ### EXPR VISITORS
// ================================================================================

func (i *Interpreter) VisitLiteralExpr(expr expr.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr expr.Logical) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitGroupingExpr(expr expr.Grouping) (interface{}, error) {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitUnaryExpr(expr expr.Unary) (interface{}, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.MINUS:
		checkNumOperand(expr.Operator, right)
		return -(right.(int)), nil
	case token.BANG:
		return !isTruthy(right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(expr expr.Variable) (interface{}, error) {
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) VisitBinaryExpr(expr expr.Binary) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.MINUS:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) - right.(int), nil
	case token.SLASH:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) / right.(int), nil
	case token.STAR:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) * right.(int), nil
	case token.PLUS:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l + r, nil
			}
		}
	case token.GREATER:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) > right.(int), nil
	case token.GREATER_EQUAL:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) >= right.(int), nil
	case token.LESS:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) < right.(int), nil
	case token.LESS_EQUAL:
		checkNumOperands(expr.Operator, left, right)
		return left.(int) <= right.(int), nil
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr expr.Call) (interface{}, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, arg := range expr.Args {
		argVal, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, argVal)
	}

	function, ok := callee.(Callable)
	if !ok {
		return nil, errors.RuntimeErr{Token: expr.Paren, Message: "Can only call functions and classes"}
	}

	if len(args) != function.Arity() {
		return nil, errors.RuntimeErr{Token: expr.Paren, Message: fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(args))}
	}

	return function.Call(i, args)
}

func (i *Interpreter) VisitAssignExpr(expr expr.Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	if distance, ok := i.locals[expr]; ok {
		err := i.environment.AssignAt(distance, expr.Name, value)
		if err != nil {
			return nil, err
		}
	} else {
		err := i.globals.Assign(expr.Name, value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

// ================================================================================
// ### STMT VISITORS
// ================================================================================

func (i *Interpreter) VisitBlockStmt(stmt stmt.Block) (stmt.ReturnValue, error) {
	return i.executeBlock(stmt.Statements, environment.NewEnvironment(i.environment))
}

func (i *Interpreter) VisitExpressionStmt(s stmt.Expression) (stmt.ReturnValue, error) {
	_, err := i.evaluate(s.Expression)
	if err != nil {
		return stmt.ReturnValue{}, err
	}
	return stmt.ReturnValue{}, nil
}

func (i *Interpreter) VisitFunctionStmt(s stmt.Function) (stmt.ReturnValue, error) {
	fn := NewFunction(s, i.environment)
	i.environment.Define(s.Name.Lexeme, fn)
	return stmt.ReturnValue{}, nil
}

func (i *Interpreter) VisitIfStmt(s stmt.If) (stmt.ReturnValue, error) {
	val, err := i.evaluate(s.Condition)
	if err != nil {
		return stmt.ReturnValue{}, err
	}

	if isTruthy(val) {
		return i.execute(s.ThenBranch)
	} else if s.ElseBranch != nil {
		return i.execute(s.ElseBranch)
	}

	return stmt.ReturnValue{}, nil
}

func (i *Interpreter) VisitPrintStmt(s stmt.Print) (stmt.ReturnValue, error) {
	val, err := i.evaluate(s.Expression)
	if err != nil {
		return stmt.ReturnValue{}, err
	}
	fmt.Printf("%v\n", val)
	return stmt.ReturnValue{}, nil
}

func (i *Interpreter) VisitReturnStmt(s stmt.Return) (stmt.ReturnValue, error) {
	var value interface{}
	if s.Value != nil {
		var err error
		value, err = i.evaluate(s.Value)
		if err != nil {
			return stmt.ReturnValue{}, err
		}
	}
	return stmt.ReturnValue{Value: value, IsReturn: true}, nil
}

func (i *Interpreter) VisitVarStmt(s stmt.Var) (stmt.ReturnValue, error) {
	var value interface{}
	if s.Initializer != nil {
		var err error
		value, err = i.evaluate(s.Initializer)
		if err != nil {
			return stmt.ReturnValue{}, err
		}
	}
	i.environment.Define(s.Name.Lexeme, value)
	return stmt.ReturnValue{}, nil
}

func (i *Interpreter) VisitWhileStmt(s stmt.While) (stmt.ReturnValue, error) {
	for {
		val, err := i.evaluate(s.Condition)
		if err != nil {
			return stmt.ReturnValue{}, err
		}

		if !isTruthy(val) {
			break
		}

		result, err := i.execute(s.Body)
		if err != nil {
			return stmt.ReturnValue{}, err
		}
		if result.IsReturn {
			return result, nil
		}
	}
	return stmt.ReturnValue{}, nil
}

// ================================================================================
// ### HELPERS
// ================================================================================

func (i *Interpreter) lookupVariable(name *token.Token, e expr.Expr) (interface{}, error) {
	distance := i.locals[e]
	if distance != 0 {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name)
	}
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
