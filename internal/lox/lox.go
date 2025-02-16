package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sjsanc/golox/internal/interpreter"
	"github.com/sjsanc/golox/internal/parser"
	"github.com/sjsanc/golox/internal/scanner"
)

type Lox struct {
	interpreter      *interpreter.Interpreter
	HadCompilerError bool
	HadRuntimeError  bool
}

func NewLox() *Lox {
	return &Lox{
		interpreter: interpreter.NewInterpreter(),
	}
}

func (l *Lox) Run(source string) {
	s := scanner.NewScanner(source)
	tokens, err := s.ScanTokens()
	p := parser.NewParser(tokens)
	stmts, err := p.Parse()

	if err {
		fmt.Println("%w", err)
		return
	}

	l.interpreter.Interpret(stmts) // Keep the same interpreter instance
}

func (l *Lox) RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		l.Run(scanner.Text())
		l.HadCompilerError = false
	}
}
