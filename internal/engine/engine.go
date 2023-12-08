package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/schema"
)

func Run(ast *language.AST, readers []reader.Reader, schemaLoader schema.SchemaLoader) error {
	_, err := schemaLoader()
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	g, err := translate(ast, readers)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	err = execute(g, readers)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	output, err := json.MarshalIndent(map[string]string{"result": "ok"}, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(output) + "\n")
	return nil
}

func translate(ast *language.AST, readers []reader.Reader) (*ExecutionGraph, error) {
	var err error
	gb := NewGraphBuilder()

	err = gb.ReadProfiles(ast, readers)
	if err != nil {
		return nil, fmt.Errorf("translate: %w", err)
	}
	err = gb.ReadNodes(ast, readers)
	if err != nil {
		return nil, fmt.Errorf("translate: %w", err)
	}

	err = gb.LinkNodes(ast)
	if err != nil {
		return nil, fmt.Errorf("translate: %w", err)
	}
	return gb.Build()
}

func execute(eg *ExecutionGraph, readers []reader.Reader, opts ...any) error {
	layers, err := traverse(eg)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	for _, layer := range layers {
		err := executeLayer(layer, readers, opts...)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}
	}
	return nil
}

func executeLayer(nodes []*Node, readers []reader.Reader, opts ...any) error {
	it := make(chan []*reader.Item, len(nodes))
	e := make(chan error)
	for _, node := range nodes {
		go executeNode(node, readers, it, e)
	}
	for i := 0; i < len(nodes); i++ {
		select {
		case items := <-it:
			nodes[i].Items = items
		case err := <-e:
			return fmt.Errorf("executeLayer: %w", err)
		}
	}
	return nil
}

func executeNode(node *Node, readers []reader.Reader, it chan []*reader.Item, e chan error) {
	r := node.Reader
	parentItems := []*reader.Item{}
	for _, p := range node.Parents {
		parentItems = append(parentItems, p.Items...)
	}
	items, err := r.GetItems(node.Type, node.Profiles, node.Attrs, node.Filters, parentItems)
	if err != nil {
		e <- fmt.Errorf("executeNode: %w", err)
	}
	it <- items
}

func traverse(eg *ExecutionGraph) ([][]*Node, error) {
	visited := map[*Node]bool{}
	layers := [][]*Node{}
	for i := 0; ; i++ {
		if i == 0 {
			layers = append(layers, eg.Roots)
		} else {
			layers = append(layers, getNextLayer(eg, layers[i-1]))
		}
		for _, j := range layers[i] {
			if _, exists := visited[j]; !exists {
				return nil, fmt.Errorf("cycle detected in execution graph: %s was visited at least twice", j.Alias)
			} else {
				visited[j] = true
			}
		}
		if len(layers[i]) == 0 {
			break
		}
	}
	return layers, nil
}

func getNextLayer(eg *ExecutionGraph, parents []*Node) []*Node {
	childSet := map[*Node]bool{}
	for _, p := range parents {
		for _, c := range p.Children {
			childSet[c] = true
		}
	}
	children := make([]*Node, len(childSet))
	i := 0
	for c := range childSet {
		children[i] = c
		i++
	}
	return children
}
