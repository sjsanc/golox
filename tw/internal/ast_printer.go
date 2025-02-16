package internal

import (
	"fmt"
	"strings"
)

type astPrinter struct {
}

func (a astPrinter) print(e expr) string {
	return e.accept(a).(string)
}

func (a astPrinter) parenthesize(name string, exprs ...expr) string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(name)
	for _, e := range exprs {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%v", e.accept(a).(string)))
	}
	sb.WriteString(")")
	return sb.String()
}

func (a astPrinter) visitBinaryExpr(e binaryExpr) interface{} {
	return a.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (a astPrinter) visitGroupingExpr(e groupingExpr) interface{} {
	return a.parenthesize("group", e.expr)
}

func (a astPrinter) visitLiteralExpr(e literalExpr) interface{} {
	if e.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.value)
}

func (a astPrinter) visitUnaryExpr(e unaryExpr) interface{} {
	return a.parenthesize(e.operator.lexeme, e.right)
}
