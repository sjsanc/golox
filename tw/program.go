package tw

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Program struct {
	interpreter *Interpreter
}

func NewProgram() *Program {
	return &Program{
		interpreter: NewInterpreter(),
	}
}

func (p *Program) RunFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	err = p.run(string(file))

	if errors.Is(err, ErrCompiler) {
		log.Println(err)
		os.Exit(64)
	}

	if errors.Is(err, ErrRuntime) {
		log.Println(err)
		os.Exit(70)
	}

	return nil
}

func (p *Program) run(source string) error {
	var compileErr bool

	scanner := NewScanner(source)
	tokens, err := scanner.Scan()
	if err {
		compileErr = true
	}

	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err {
		compileErr = true
	}

	// fmt.Println(compileErr)
	// printer := &Printer{}
	// for _, stmt := range statements {
	// 	fmt.Println(printer.PrintStmt(stmt))
	// }

	if compileErr {
		return ErrCompiler
	}

	resolver := NewResolver(p.interpreter)
	err = resolver.Resolve(statements)
	if err {
		compileErr = true
	}

	if compileErr {
		return ErrCompiler
	}

	runtimeErr := p.interpreter.Interpret(statements)
	if runtimeErr != nil {
		fmt.Println(runtimeErr)
		return ErrRuntime
	}

	return nil
}
