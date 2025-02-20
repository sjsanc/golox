package main

import (
	"log"
	"os"

	"github.com/sjsanc/golox/tw"
)

func main() {
	if len(os.Args) > 1 {
		log.Println("Usage: golox [script]")
		os.Exit(64)
	}

	program := tw.NewProgram()
	program.RunFile("test.lox")
}
