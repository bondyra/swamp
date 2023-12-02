package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
)

type DAG struct {
	profiles []string
	nodes    []Node
}

type Node struct {
	Type          string
	Alias         string
	Filters       []reader.Filter
	Attrs         []string
	ParentContext reader.ParentContext
}

func (d DAG) Traverse() []Node { return []Node{} }

func Run(ast *language.AST, readers []reader.Reader) (map[string]any, error) {
	dag, _ := translate(ast, readers)
	return executeQuery(dag, readers)
}

func translate(ast *language.AST, readers []reader.Reader) (*DAG, error) {
	var profiles []string
	var err error
	if ast.Profiles.All {
		profiles, err = getAllSupportedProfiles(readers)
		if err != nil {
			return nil, fmt.Errorf("translate: %w", err)
		}
	} else {
		profiles = ast.Profiles.Elements
	}
	// todo validate if profiles exist
	return &DAG{
		profiles: profiles,
		nodes:    []Node{},
	}, nil
}

func getAllSupportedProfiles(readers []reader.Reader) ([]string, error) {
	results := []string{}
	for _, r := range readers {
		results = append(results, r.GetSupportedProfiles()...)
	}
	// todo check if they intersect
	return results, nil
}

func executeQuery(d *DAG, readers []reader.Reader, opts ...any) (map[string]any, error) {
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
