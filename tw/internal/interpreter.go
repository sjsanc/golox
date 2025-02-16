package internal

import "fmt"

type interpreter struct {
}

func (i *interpreter) interpret(expr expr) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error?")
		}
	}()

	value := i.evaluate(expr)
	fmt.Println(value)
}

func (i *interpreter) visitLiteralExpr(expr literalExpr) interface{} {
	return expr.value
}

func (i *interpreter) visitGroupingExpr(expr groupingExpr) interface{} {
	return i.evaluate(expr.expr)
}

func (i *interpreter) visitUnaryExpr(expr unaryExpr) interface{} {
	right := i.evaluate(expr.right)

	switch expr.operator.ttype {
	case MINUS:
		checkNumberOperand(expr.operator, right)
		return -(right.(int))
	case BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *interpreter) visitBinaryExpr(expr binaryExpr) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.ttype {
	case MINUS:
		checkNumberOperands(expr.operator, left, right)
		return left.(int) - right.(int)
	case SLASH:
		checkNumberOperands(expr.operator, left, right)
		return left.(int) / right.(int)
	case STAR:
		checkNumberOperands(expr.operator, left, right)
		return left.(int) * right.(int)
	case PLUS:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l + r
			}
		}

		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}

		return nil
	case GREATER:
		checkNumberOperands(expr.operator, left, right)
		return left.(int) > right.(int)
	case GREATER_EQUAL:
		checkNumberOperands(expr.operator, left, right)
		return left.(int) >= right.(int)
	case LESS:
		checkNumberOperands(expr.operator, left, right)
		return left.(int) < right.(int)
	case LESS_EQUAL:
		checkNumberOperands(expr.operator, left, right)
		return left.(int) <= right.(int)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i *interpreter) evaluate(expr expr) interface{} {
	return expr.accept(i)
}

func checkNumberOperand(operator *token, operand interface{}) {
	if _, ok := operand.(int); !ok {
		panic(RuntimeError{operator, "Operand must be a number"})
	} else {
		return
	}
}

func checkNumberOperands(operator *token, left, right interface{}) {
	if _, ok := left.(int); !ok {
		panic(RuntimeError{operator, "Left operand must be a number"})
	}
	if _, ok := right.(int); !ok {
		panic(RuntimeError{operator, "Right operand must be a number"})
	}
}

func (i *interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}

func (i *interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
