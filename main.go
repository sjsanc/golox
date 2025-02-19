package main

import (
	"log"
	"os"

	"github.com/sjsanc/golox/internal/lox"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		log.Println("Usage: golox [script]")
		os.Exit(64)
	}

	Program := lox.NewLox() // Use constructor to maintain state
	Program.RunFile("test.lox")
}
