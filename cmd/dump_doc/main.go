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

func associated_comments(
	function bashdoc.Function,
	comments bashdoc.CommentsByLine,
) []bashdoc.Comment {
	var associated []bashdoc.Comment

	line := function.Start.Line() - 1

	for line > 0 {
		if comment, found := comments.Comments[line]; found {
			associated = append([]bashdoc.Comment{comment}, associated...)
		} else {
			break
		}

		line--
	}

	return associated
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

	fmt.Println("Functions found:")
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

	comments, err := bashdoc.LoadCommentsFromSource(reader)
	if err != nil {
		return fmt.Errorf("error loading comments from file: %s", err)
	}

	fmt.Println("Comments found:")
	for _, line := range comments.LinesWithComments() {
		comment := comments.Comments[line]
		fmt.Printf(
			"On line %d, found comment at %d, %d: %s\n",
			line,
			comment.Pos().Line(),
			comment.Pos().Col(),
			comment.Text,
		)
	}

	for _, f := range functions {
		associated := associated_comments(f, *comments)
		fmt.Printf(
			"Associated comments for function %s:\n",
			f.Name,
		)
		for _, comment := range associated {
			fmt.Printf(
				"Comment at %d, %d: %s\n",
				comment.Pos().Line(),
				comment.Pos().Col(),
				comment.Text,
			)
		}
	}

	return nil
}
