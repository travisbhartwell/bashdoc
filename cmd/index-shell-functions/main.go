package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/travisbhartwell/bashdoc"
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
	if len(arguments) < 3 {
		return fmt.Errorf("no file specified")
	}

	script := arguments[1]
	fmt.Printf("Analyzing script: %s\n", script)

	reader, err := os.Open(script)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	defer reader.Close()

	csvFile := arguments[2]
	fmt.Printf("Outputing index to %s\n", csvFile)

	writer, err := os.Create(csvFile)
	if err != nil {
		return fmt.Errorf("error opening file for write: %w", err)
	}
	defer writer.Close()

	return indexLines(reader, writer, script)
}

const OUTSIDE_OF_FUNCTION_DECLARATION = "OUTSIDE_OF_FUNCTION_DECLARATION"

func functionNameForLine(functions []bashdoc.Function, pos syntax.Pos) string {
	for _, f := range functions {
		if f.IsWithinDeclaration(pos) {
			return f.Name
		}
	}

	return OUTSIDE_OF_FUNCTION_DECLARATION
}

func indexLines(reader io.ReadSeeker, writer io.Writer, fileName string) error {
	functions, err := bashdoc.LoadFunctionsFromSource(reader)
	if err != nil {
		return fmt.Errorf("error loading functions from file: %w", err)
	}

	for _, f := range functions {
		fmt.Printf(
			"Found function %s at %d, %d\n",
			f.Name,
			f.Start.Line(),
			f.Start.Col(),
		)
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error rewinding reader: %w", err)
	}

	lines, err := bashdoc.LoadLinesWithCode(reader)
	if err != nil {
		return fmt.Errorf("error loading lines with code: %w", err)
	}

	for _, line := range lines {
		fmt.Printf("Found line at %d, %d\n", line.Line(), line.Col())
	}

	records := [][]string{
		{"filename", "functionname", "linenumber"},
	}

	for _, line := range lines {
		fnName := functionNameForLine(functions, line)
		records = append(records, []string{
			fileName, fnName, fmt.Sprintf("%d", line.Line()),
		})
	}

	fmt.Printf("Records: %+v", records)

	csvWriter := csv.NewWriter(writer)

	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("error writing record to csv: %w", err)
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("error writing to csv file: %w", err)
	}

	return nil
}
