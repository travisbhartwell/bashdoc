package bashdoc

import (
	"fmt"
	"io"
	"slices"
	"sort"

	"mvdan.cc/sh/v3/syntax"
)

type Function struct {
	Name  string
	Start syntax.Pos
	End   syntax.Pos
}

func (f *Function) IsWithinDeclaration(pos syntax.Pos) bool {
	return pos == f.Start || pos == f.End || (pos.After(f.Start) && f.End.After(pos))
}

func SortedFunctions(functions []Function) {
	sort.Slice(functions, func(i, j int) bool {
		return functions[j].Start.After(functions[i].Start)
	})
}

func LoadLinesWithCode(reader io.Reader) ([]syntax.Pos, error) {
	linesSeen := make(map[uint]syntax.Pos)
	var lines []syntax.Pos

	parsedFile, err := syntax.NewParser(syntax.KeepComments(false)).Parse(reader, "")
	if err != nil {
		return lines, fmt.Errorf("error loading lines from source: %w", err)
	}

	syntax.Walk(parsedFile, func(node syntax.Node) bool {
		if node != nil {
			line := node.Pos().Line()

			_, ok := linesSeen[line]
			if !ok {
				linesSeen[line] = node.Pos()
			}
		}

		return true
	})

	keys := make([]uint, len(linesSeen))
	i := 0
	for k := range linesSeen {
		keys[i] = k
		i++
	}

	slices.Sort(keys)

	lines = make([]syntax.Pos, len(keys))
	for i, k := range keys {
		lines[i] = linesSeen[k]
	}

	return lines, nil
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
				f := Function{Name: x.Name.Value, Start: x.Pos(), End: x.End()}
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
	var commentsByLine = make(map[uint]Comment)

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
