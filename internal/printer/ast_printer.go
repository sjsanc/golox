package printer

import (
	"fmt"
	"strings"

	"github.com/sjsanc/golox/internal/expr"
)

type ASTPrinter struct{}

func (a ASTPrinter) Print(e expr.Expr) string {
	return e.Accept(a).(string)
}

func (a ASTPrinter) VisitBinaryExpr(e expr.Binary) interface{} {
	return a.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (a ASTPrinter) VisitGroupingExpr(e expr.Grouping) interface{} {
	return a.parenthesize("group", e.Expr)
}

func (a ASTPrinter) VisitLiteralExpr(e expr.Literal) interface{} {
	if e.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}

func (a ASTPrinter) VisitUnaryExpr(e expr.Unary) interface{} {
	return a.parenthesize(e.Operator.Lexeme, e.Right)
}

func (a ASTPrinter) VisitVariableExpr(e expr.Variable) interface{} {
	return e.Name.Lexeme
}

func (a ASTPrinter) parenthesize(name string, exprs ...expr.Expr) string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(name)
	for _, e := range exprs {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%v", e.Accept(a).(string)))
	}
	sb.WriteString(")")
	return sb.String()
}
