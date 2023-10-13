package main

import (
	"fmt"
	"io"
	"os"

	"github.com/travisbhartwell/bashdoc"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	reader, err := os.Open("test-script")
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	defer reader.Close()

	functions, err := bashdoc.LoadFunctionsFromSource(reader)
	if err != nil {
		return fmt.Errorf("error loading functions from file: %w", err)
	}

	for _, f := range functions {
		fmt.Printf(
			"Found function %s at %d, %d\n",
			f.Name,
			f.DeclaredAt.Line(),
			f.DeclaredAt.Col(),
		)
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error rewinding rearder: %w", err)
	}

	comments, err := bashdoc.LoadCommentsFromSource(reader)
	if err != nil {
		return fmt.Errorf("error loading comments from file: %s", err)
	}

	for line, comment := range comments {
		fmt.Printf(
			"On line %d, found comment at %d, %d: %s\n",
			line,
			comment.Pos().Line(),
			comment.Pos().Col(),
			comment.Text,
		)
	}

	return nil
}
