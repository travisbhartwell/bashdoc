package main

import (
	"fmt"
	"github.com/travisbhartwell/bashdoc"
	"os"
)

func main() {
	reader, err := os.Open("test-script")
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		os.Exit(1)
	}

	functions, err := bashdoc.LoadFunctionsFromSource(reader)
	if err != nil {
		fmt.Printf("Error loading functions from file: %s", err)
		os.Exit(1)
	}

	for _, f := range functions {
		fmt.Printf("Found function %s at %d, %d\n", f.Name, f.DeclaredAt.Line(), f.DeclaredAt.Col())
	}
}
