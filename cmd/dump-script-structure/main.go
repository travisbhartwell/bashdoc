package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"mvdan.cc/sh/v3/syntax"
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run(arguments []string) error {

	if len(arguments) == 1 {
		return fmt.Errorf("no file specified")
	}

	script := arguments[1]
	log.Printf("Analyzing script: %s", script)

	reader, err := os.Open(script)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	defer reader.Close()

	if err := dumpStructure(reader); err != nil {
		return fmt.Errorf("error dumping structure: %w", err)
	}

	return nil
}

func dumpStructure(reader io.Reader) error {
	parser, err := syntax.NewParser(syntax.KeepComments(true)).Parse(reader, "")

	if err != nil {
		return fmt.Errorf("error parsing script source: %w", err)
	}

	syntax.Walk(parser, func(node syntax.Node) bool {
		if node != nil {
			pos := node.Pos()
			fmt.Printf("Node %T at (%d, %d): %+v\n", node, pos.Line(), pos.Col(), node)
		}

		return true
	})

	return nil
}
