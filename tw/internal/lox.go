package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	log.Printf("[line %d] Error %s: %s", line, where, message)
	l.hadError = true
}

func (l *Lox) errorToken(token *token, message string) {
	if token.ttype == EOF {
		l.report(token.line, "at end", message)
	} else {
		l.report(token.line, "at '"+token.lexeme+"'", message)
	}
}

func (l *Lox) runtimeError(error RuntimeError) {
	log.Printf("%s\n[line %d]", error.Message, error.Token.line)
	l.hadRuntimeError = true
}

func (l *Lox) run(source string) {
	scanner := newScanner(source)
	tokens := scanner.scanTokens()
	parser := newParser(tokens)
	expr := parser.parse()

	if l.hadError {
		return
	}

	interpreter := interpreter{}
	interpreter.interpret(expr)
}

func (l *Lox) RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print("> ")
		l.run(scanner.Text())
		l.hadError = false
	}
}

var Program = &Lox{}

func (l *Lox) runFile(path string) error {
	fd, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	l.run(string(fd))

	if l.hadError {
		os.Exit(65)
	}

	if l.hadRuntimeError {
		os.Exit(70)
	}

	return nil
}
