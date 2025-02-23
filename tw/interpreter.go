package tw

import (
	"fmt"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewGlobalEnvironment()

	globals.Define("clock", &ClockBuiltin{})

	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[Expr]int),
	}
}

func (i *Interpreter) Interpret(stmts []Stmt) error {
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(s Stmt) (StmtReturn, error) {
	result, err := s.Accept(i)
	if err != nil {
		return StmtReturn{}, err
	}
	return result, nil
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) (StmtReturn, error) {
	prev := i.environment
	i.environment = env
	defer func() { i.environment = prev }()
	for _, s := range stmts {
		result, err := i.execute(s)
		if err != nil {
			return StmtReturn{}, err
		}
		if result.isReturn {
			return result, nil
		}
	}
	return StmtReturn{}, nil
}

// ================================================================================
// ### STMT VISITORS
// ================================================================================

func (i *Interpreter) visitBlockStmt(stmt *BlockStmt) (StmtReturn, error) {
	return i.executeBlock(stmt.stmts, NewEnvironment(i.environment))
}

func (i *Interpreter) visitClassStmt(stmt *ClassStmt) (StmtReturn, error) {
	var superclass interface{}
	if stmt.superclass != nil {
		val, err := i.evaluate(stmt.superclass)
		if err != nil {
			return StmtReturn{}, err
		}

		val, ok := val.(Class) // Type assertion
		if !ok {
			i.error(stmt.superclass.name, "Superclass must be a class") // Ensure 'Name' exists
			return StmtReturn{}, nil
		}
		superclass = val
	}

	i.environment.Define(stmt.name.lexeme, nil)

	if stmt.superclass != nil {
		environment := NewEnvironment(i.environment)
		environment.Define("super", superclass)
	}

	methods := make(map[string]*Function)
	for _, method := range stmt.methods {
		function := NewFunction(method, i.environment, method.name.lexeme == "init")
		methods[method.name.lexeme] = function
	}

	class := NewClass(stmt.name.lexeme, superclass.(*Class), methods)
	if superclass != nil {
		i.environment = i.environment.enclosing
	}
	i.environment.Assign(stmt.name, class)
	return StmtReturn{}, nil
}

func (i *Interpreter) visitExpressionStmt(stmt *ExpressionStmt) (StmtReturn, error) {
	_, err := i.evaluate(stmt.expr)
	if err != nil {
		return StmtReturn{}, err
	}
	return StmtReturn{}, nil
}

func (i *Interpreter) visitFunctionStmt(stmt *FunctionStmt) (StmtReturn, error) {
	function := NewFunction(stmt, i.environment, false)
	i.environment.Define(stmt.name.lexeme, function)
	return StmtReturn{}, nil
}

func (i *Interpreter) visitIfStmt(stmt *IfStmt) (StmtReturn, error) {
	condition, err := i.evaluate(stmt.condition)
	if err != nil {
		return StmtReturn{}, err
	}
	if isTruthy(condition) {
		return i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		return i.execute(stmt.elseBranch)
	}
	return StmtReturn{}, nil
}

func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) (StmtReturn, error) {
	value, err := i.evaluate(stmt.expr)
	if err != nil {
		return StmtReturn{}, err
	}
	fmt.Println(value)
	return StmtReturn{}, nil
}

func (i *Interpreter) visitReturnStmt(stmt *ReturnStmt) (StmtReturn, error) {
	var value interface{}
	if stmt.value != nil {
		v, err := i.evaluate(stmt.value)
		if err != nil {
			return StmtReturn{}, err
		}
		value = v
	}
	return StmtReturn{value, true}, nil
}

func (i *Interpreter) visitVarStmt(stmt *VarStmt) (StmtReturn, error) {
	var value interface{}
	if stmt.initializer != nil {
		v, err := i.evaluate(stmt.initializer)
		if err != nil {
			return StmtReturn{}, err
		}
		value = v
	}
	i.environment.Define(stmt.name.lexeme, value)
	return StmtReturn{}, nil
}

func (i *Interpreter) visitWhileStmt(stmt *WhileStmt) (StmtReturn, error) {
	for {
		condition, err := i.evaluate(stmt.condition)
		if err != nil {
			return StmtReturn{}, err
		}
		if !isTruthy(condition) {
			break
		}
		result, err := i.execute(stmt.body)
		if err != nil {
			return StmtReturn{}, err
		}
		if result.isReturn {
			return result, nil
		}
	}
	return StmtReturn{}, nil
}

// ================================================================================
// ### EXPR VISITORS
// ================================================================================

func (i *Interpreter) visitAssignExpr(expr *AssignExpr) (interface{}, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	distance, ok := i.locals[expr]
	if ok {
		i.environment.AssignAt(distance, expr.name, value)
	} else {
		i.globals.Assign(expr.name, value)
	}
	return value, nil
}

func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.ttype {
	case MINUS:
		err := i.checkNumOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(int) - right.(int), nil
	case SLASH:
		err := i.checkNumOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(int) / right.(int), nil
	case STAR:
		err := i.checkNumOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(int) * right.(int), nil
	case PLUS:
		if _, ok := left.(int); ok {
			if _, ok := right.(int); ok {
				return left.(int) + right.(int), nil
			}
		}
		if _, ok := left.(string); ok {
			if _, ok := right.(string); ok {
				return left.(string) + right.(string), nil
			}
		}
	case GREATER:
		err := i.checkNumOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(int) > right.(int), nil
	case GREATER_EQUAL:
		err := i.checkNumOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(int) >= right.(int), nil
	case LESS:
		err := i.checkNumOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(int) < right.(int), nil
	case LESS_EQUAL:
		err := i.checkNumOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(int) <= right.(int), nil
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	}
	return nil, nil
}

func (i *Interpreter) visitCallExpr(expr *CallExpr) (interface{}, error) {
	callee, err := i.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, 0)
	for _, arg := range expr.args {
		value, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}
	if f, ok := callee.(Callable); ok {
		if len(args) != f.Arity() {
			return nil, i.error(expr.paren, fmt.Sprintf("Expected %d arguments but got %d", f.Arity(), len(args)))
		}
		return f.Call(i, args)
	}
	return nil, i.error(expr.paren, "Can only call functions and classes")
}

func (i *Interpreter) visitGetExpr(expr *GetExpr) (interface{}, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}
	if instance, ok := object.(*Instance); ok {
		return instance.Get(expr.name)
	}
	return nil, i.error(expr.name, "Only instances have properties")
}

func (i *Interpreter) visitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	return i.evaluate(expr.expr)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return expr.value, nil
}

func (i *Interpreter) visitLogicalExpr(expr *LogicalExpr) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	if expr.operator.ttype == OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return i.evaluate(expr.right)
}

func (i *Interpreter) visitSetExpr(expr *SetExpr) (interface{}, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}
	if instance, ok := object.(*Instance); ok {
		value, err := i.evaluate(expr.value)
		if err != nil {
			return nil, err
		}
		instance.Set(expr.name, value)
		return value, nil
	}
	return nil, i.error(expr.name, "Only instances have fields")
}

func (i *Interpreter) visitSuperExpr(expr *SuperExpr) (interface{}, error) {
	distance := i.locals[expr]
	superclass, err := i.environment.GetAt(distance, "super")
	if err != nil {
		return nil, err
	}

	object, err := i.environment.GetAt(distance-1, "this")
	if err != nil {
		return nil, err
	}

	method := superclass.(*Class).FindMethod(expr.method.lexeme)

	if method == nil {
		i.error(expr.method, "undefined property"+expr.method.lexeme+"'.")
		return nil, nil
	}

	return method.Bind(object.(*Instance)), nil
}

func (i *Interpreter) visitThisExpr(expr *ThisExpr) (interface{}, error) {
	return i.lookupVariable(expr.keyword, expr)
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.ttype {
	case MINUS:
		err := i.checkNumOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(int), nil
	case BANG:
		return !isTruthy(right), nil
	}
	return nil, nil
}

func (i *Interpreter) visitVariableExpr(expr *VariableExpr) (interface{}, error) {
	return i.lookupVariable(expr.name, expr)
}

// ================================================================================
// ### HELPERS
// ================================================================================

func (i *Interpreter) lookupVariable(name *Token, expr Expr) (interface{}, error) {
	distance := i.locals[expr]
	if distance != 0 {
		return i.environment.GetAt(distance, name.lexeme)
	} else {
		return i.globals.Get(name)
	}
}

func (i *Interpreter) checkNumOperand(operator *Token, operand interface{}) error {
	if _, ok := operand.(int); !ok {
		return i.error(operator, "Operand must be a number")
	}
	return nil
}

func (i *Interpreter) checkNumOperands(operator *Token, left, right interface{}) error {
	if _, ok := left.(int); !ok {
		return i.error(operator, "Left operand must be a number")
	}
	if _, ok := right.(int); !ok {
		return i.error(operator, "Right operand must be a number")
	}
	return nil
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

func (i *Interpreter) error(token *Token, message string) error {
	return fmt.Errorf("[line %d] Error: %s", token.line, message)
}
