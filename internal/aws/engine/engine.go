package engine

import (
	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
)

type DAG struct {
}

func Run(reader reader.Reader, ast *language.AST) (map[string]any, error) {
	dag := translate(ast)
	return executeQuery(reader, dag)
}

func translate(ast *language.AST) *DAG {
	return &DAG{}
}

func executeQuery(reader reader.Reader, d *DAG, opts ...any) (map[string]any, error) {
	// var parentContext any
	// for n := range d.Traverse() {
	// 	items, err := reader.GetItems(..., parentContext)
	// 	if err != nil {
	// 		// todo some additional debug context, e.g. layer number
	// 		return nil, fmt.Errorf("executeQuery: %w", err)
	// 	}
	// 	parentContext := d.Context()
	// }
	return map[string]any{"a": "1", "b": 1, "c": true}, nil
}
