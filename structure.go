package bashdoc

import (
	"fmt"
	"io"
	"slices"
	"sort"

	"mvdan.cc/sh/v3/syntax"
)

type Function struct {
	Name       string
	DeclaredAt syntax.Pos
}

func SortedFunctions(functions []Function) {
	sort.Slice(functions, func(i, j int) bool {
		return functions[i].DeclaredAt.After(functions[j].DeclaredAt)
	})
}

func LoadFunctionsFromSource(reader io.Reader) ([]Function, error) {
	var functions []Function

	parser, err := syntax.NewParser(syntax.KeepComments(true)).Parse(reader, "")
	if err != nil {
		return functions, fmt.Errorf("error loading functions from source: %w", err)
	}

	syntax.Walk(parser, func(node syntax.Node) bool {
		if node != nil {
			if x, ok := node.(*syntax.FuncDecl); ok {
				f := Function{Name: x.Name.Value, DeclaredAt: x.Position}
				functions = append(functions, f)
			}
		}

		return true
	})

	SortedFunctions(functions)

	return functions, nil
}

type Comment struct {
	syntax.Comment
}

type CommentsByLine struct {
	Comments map[uint]Comment
}

func (c *CommentsByLine) LinesWithComments() []uint {
	var lines []uint
	for line := range c.Comments {
		lines = append(lines, line)
	}
	slices.Sort(lines)
	return lines
}

func LoadCommentsFromSource(reader io.Reader) (*CommentsByLine, error) {
	var commentsByLine map[uint]Comment = make(map[uint]Comment)

	parser, err := syntax.NewParser(syntax.KeepComments(true)).Parse(reader, "")
	if err != nil {
		return nil, fmt.Errorf("error loading comment from source: %w", err)
	}

	syntax.Walk(parser, func(node syntax.Node) bool {
		if node != nil {
			if x, ok := node.(*syntax.Comment); ok {
				commentsByLine[x.Hash.Line()] = Comment{*x}
			}
		}

		return true
	})

	return &CommentsByLine{Comments: commentsByLine}, nil
}
