package bashdoc

import (
	"fmt"
	"io"
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
		return functions, fmt.Errorf("loading functions from source: %w", err)
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
