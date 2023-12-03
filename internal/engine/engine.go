package engine

import (
	"fmt"
	"strings"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
)

func Run(ast *language.AST, readers []reader.Reader) error {
	g, err := translate(ast, readers)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	err = execute(g, readers)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	printGraph(g)
	return nil
}

func translate(ast *language.AST, readers []reader.Reader) (*ExecutionGraph, error) {
	var err error
	gb := NewGraphBuilder()

	err = gb.ReadProfiles(ast, readers)
	if err != nil {
		return nil, fmt.Errorf("translate: %w", err)
	}
	err = gb.ReadNodes(ast)
	if err != nil {
		return nil, fmt.Errorf("translate: %w", err)
	}

	err = gb.LinkNodes(ast)
	if err != nil {
		return nil, fmt.Errorf("translate: %w", err)
	}
	return gb.Build(), nil
}

func execute(eg *ExecutionGraph, readers []reader.Reader, opts ...any) error {
	layer, err := eg.GetRoots()
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	for len(layer) > 0 {
		err := executeLayer(layer, readers, opts...)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}
		layer = getNextLayer(eg, layer)
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
			nodes[i].items = items
		case err := <-e:
			return fmt.Errorf("executeLayer: %w", err)
		}
	}
	return nil
}

func executeNode(node *Node, readers []reader.Reader, it chan []*reader.Item, e chan error) {
	for _, r := range readers {
		if r.Name() == node.ReaderName {
			parentItems := []*reader.Item{}
			for _, p := range node.parents {
				parentItems = append(parentItems, p.items...)
			}
			items, err := r.GetItems(node.Type, node.Profiles, node.Attrs, node.Filters, parentItems)
			if err != nil {
				e <- fmt.Errorf("executeNode: %w", err)
			}
			it <- items
		}
	}
	e <- fmt.Errorf("executeNode: reader %s not found", node.ReaderName)
}
func getNextLayer(eg *ExecutionGraph, parents []*Node) []*Node {
	childSet := map[*Node]bool{}
	for _, p := range parents {
		for _, c := range eg.GetChildren(p) {
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

type GraphBuilder struct {
	profiles []string
	nodeMap  map[string]*Node
}

func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{nodeMap: map[string]*Node{}}
}

func (gb *GraphBuilder) ReadProfiles(ast *language.AST, readers []reader.Reader) error {
	allSupportedProfiles, err := getAllSupportedProfiles(readers)
	if err != nil {
		return fmt.Errorf("ReadProfiles: %w", err)
	}
	if ast.Profiles.All {
		gb.profiles = allSupportedProfiles
	} else {
		unknownProfiles := common.Difference(ast.Profiles.Elements, allSupportedProfiles)
		if len(unknownProfiles) > 0 {
			return fmt.Errorf("ReadProfiles: unknown profiles: %s", strings.Join(unknownProfiles, ", "))
		}
		gb.profiles = ast.Profiles.Elements
	}
	return nil
}

func getAllSupportedProfiles(readers []reader.Reader) ([]string, error) {
	results := []string{}
	for _, r := range readers {
		results = append(results, r.GetSupportedProfiles()...)
	}
	// todo err if profiles from differrent readers intersect
	return results, nil
}

func (gb *GraphBuilder) ReadNodes(ast *language.AST) error {
	for _, e := range ast.Entities {
		switch e.(type) {
		case language.ItemEntity:
			ie := e.(language.ItemEntity).Value
			nodeType := ie.Type[len(ie.Type)-1]
			gb.addNode(ie.Alias, "aws", nodeType, gb.profiles, []reader.Filter{}, []string{})
		}
	}
	return nil
}

func (gb *GraphBuilder) addNode(alias string, readerName string, nodeType string, profiles []string, filters []reader.Filter, attrs []string) error {
	if _, exists := gb.nodeMap[alias]; exists {
		return fmt.Errorf("addNode: alias %s already defined", alias)
	}
	gb.nodeMap[alias] = &Node{Alias: alias, ReaderName: readerName, Type: nodeType, Profiles: profiles, Filters: filters, Attrs: attrs}
	return nil
}

func (gb *GraphBuilder) LinkNodes(ast *language.AST) error {
	for _, e := range ast.Entities {
		switch e.(type) {
		case language.LinkEntity:
			le := e.(language.LinkEntity).Value
			err := gb.linkNode(le.From, le.To)
			if err != nil {
				return fmt.Errorf("LinkNodes: %w", err)
			}
		}
	}
	return nil
}

func (gb *GraphBuilder) linkNode(fromAlias string, toAlias string) error {
	_, fromExists := gb.nodeMap[fromAlias]
	_, toExists := gb.nodeMap[toAlias]
	if !fromExists || !toExists {
		var errors []string = []string{}
		if !fromExists {
			errors = append(errors, fmt.Sprintf("alias %s is not defined", fromAlias))
		}
		if !toExists {
			errors = append(errors, fmt.Sprintf("alias %s is not defined", toAlias))
		}
		return fmt.Errorf("linkNode: %s", strings.Join(errors, ", "))
	}
	gb.nodeMap[fromAlias].children = append(gb.nodeMap[fromAlias].children, gb.nodeMap[toAlias])
	gb.nodeMap[toAlias].parents = append(gb.nodeMap[toAlias].parents, gb.nodeMap[fromAlias])
	return nil
}

func (gb *GraphBuilder) Build() *ExecutionGraph {
	return &ExecutionGraph{nodeMap: gb.nodeMap}
}
