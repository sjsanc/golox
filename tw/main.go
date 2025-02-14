package main

import (
	"log"
	"os"

	"github.com/sjsanc/golox/internal"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		log.Println("Usage: golox [script]")
		os.Exit(64)
	}

	internal.Program.RunPrompt()
}
