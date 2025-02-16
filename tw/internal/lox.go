package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Lox struct {
	hadError bool
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	log.Printf("[line %d] Error %s: %s", line, where, message)
	l.hadError = true
}

func (l *Lox) run(source string) {
	scanner := newScanner(source)
	tokens := scanner.scanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func (l *Lox) RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	e := binaryExpr{
		left: unaryExpr{
			token{MINUS, "-", nil, 1},
			literalExpr{123},
		},
		operator: token{STAR, "*", nil, 1},
		right: groupingExpr{
			literalExpr{45.67},
		},
	}

	fmt.Println(e.accept(astPrinter{}).(string))

	for scanner.Scan() {
		fmt.Print("> ")
		l.run(scanner.Text())
		l.hadError = false
	}
}

var Program = &Lox{}

// func (l *Lox) runFile(path string) error {
// 	fd, err := os.ReadFile(path)
// 	if err != nil {
// 		return err
// 	}

// 	l.run(string(fd))

// 	if l.HadError {
// 		os.Exit(65)
// 	}

// 	return nil
// }
