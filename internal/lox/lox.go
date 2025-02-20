package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sjsanc/golox/internal/interpreter"
	"github.com/sjsanc/golox/internal/parser"
	"github.com/sjsanc/golox/internal/resolver"
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
	tokens, _ := s.ScanTokens()
	// if err {
	// 	fmt.Println("Error scanning tokens: %w", err)
	// }

	p := parser.NewParser(tokens)
	stmts, _ := p.Parse()
	// if err {
	// 	fmt.Println("Error parsing: %w", err)
	// 	return
	// }

	r := resolver.NewResolver(l.interpreter)
	err := r.Resolve(stmts)
	if err != nil {
		fmt.Println(err)
		l.HadCompilerError = true
		return
	}

	l.interpreter.Interpret(stmts)
}

func (l *Lox) RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		l.Run(scanner.Text())
		l.HadCompilerError = false
	}
}
func (l *Lox) RunFile(path string) {
	file, err := os.ReadFile(path) // Read the entire file at once
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	l.Run(string(file)) // Pass full content to Run
}
