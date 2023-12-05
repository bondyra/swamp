package engine

import (
	"encoding/json"
	"fmt"
	"log"
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

func (gb *GraphBuilder) ReadNodes(ast *language.AST, readers []reader.Reader) error {
	for _, e := range ast.Entities {
		switch ie := e.(type) {
		case language.ItemEntity:
			r, err := getReaderForType(ie.Value.Type, readers)
			if err != nil {
				return fmt.Errorf("ReadNodes: %w", err)
			}
			nodeType := ie.Value.Type[len(ie.Value.Type)-1]
			gb.addNode(ie.Value.Alias, "aws", nodeType, gb.profiles, []reader.Filter{}, []string{}, r)
		}
	}
	return nil
}

func getReaderForType(typ []string, readers []reader.Reader) (reader.Reader, error) {
	if len(typ) == 0 {
		return nil, fmt.Errorf("getReaderForType: nodeType is empty")
	}
	if len(typ) == 1 {
		t, err := guessReaderForType(typ[0], readers)
		if err != nil {
			return nil, fmt.Errorf("getReaderForType: %w", err)
		}
		return t, nil
	}
	if len(typ) > 2 {
		return nil, fmt.Errorf("getReaderForType: type is too long") // should not normally happen unless parser is wrong
	}

	for _, r := range readers {
		if r.Name() == typ[0] && r.IsTypeSupported(typ[1]) {
			return r, nil
		}
	}
	return nil, fmt.Errorf("getReaderForType: no reader found for type %s", typ)
}

func guessReaderForType(typ string, readers []reader.Reader) (reader.Reader, error) {
	matchingReaders := []reader.Reader{}
	for _, r := range readers {
		if r.IsTypeSupported(typ) {
			matchingReaders = append(matchingReaders, r)
		}
	}
	if len(matchingReaders) == 0 {
		return nil, fmt.Errorf("guessReaderForType: no reader found for type %s", typ)
	}
	if len(matchingReaders) > 1 {
		return nil, fmt.Errorf(
			"guessReaderForType: type %s is defined in multiple readers: %s",
			typ, common.Map(matchingReaders, func(r reader.Reader) string { return r.Name() }),
		)
	}
	// todo add debug log
	return matchingReaders[0], nil
}

func (gb *GraphBuilder) addNode(alias string, readerName string, nodeType string, profiles []string, filters []reader.Filter, attrs []string, r reader.Reader) error {
	if _, exists := gb.nodeMap[alias]; exists {
		return fmt.Errorf("addNode: alias %s already defined", alias)
	}
	gb.nodeMap[alias] = &Node{Alias: alias, Reader: r, Type: nodeType, Profiles: profiles, Filters: filters, Attrs: attrs}
	return nil
}

func (gb *GraphBuilder) LinkNodes(ast *language.AST) error {
	for _, e := range ast.Entities {
		switch le := e.(type) {
		case language.LinkEntity:
			err := gb.linkNode(le.Value.From, le.Value.To)
			if err != nil {
				return fmt.Errorf("LinkNodes: %w", err)
			}
		}
	}
	return nil
}

func (gb *GraphBuilder) linkNode(parentAlias string, childAlias string) error {
	_, childDefined := gb.nodeMap[parentAlias]
	_, parentDefined := gb.nodeMap[childAlias]
	if !childDefined || !parentDefined {
		var errors []string = []string{}
		if !parentDefined {
			errors = append(errors, fmt.Sprintf("alias %s is not defined", parentAlias))
		}
		if !childDefined {
			errors = append(errors, fmt.Sprintf("alias %s is not defined", childAlias))
		}
		return fmt.Errorf("linkNode: %s", strings.Join(errors, ", "))
	}
	if !gb.nodeMap[childAlias].Reader.IsLinkSupported(gb.nodeMap[childAlias].Type, gb.nodeMap[parentAlias].Reader.Name(), gb.nodeMap[parentAlias].Type) {
		return fmt.Errorf("linkNode: link from %s to %s is not supported", parentAlias, childAlias)
	}
	gb.nodeMap[parentAlias].Children = append(gb.nodeMap[parentAlias].Children, gb.nodeMap[childAlias])
	gb.nodeMap[childAlias].Parents = append(gb.nodeMap[childAlias].Parents, gb.nodeMap[parentAlias])
	return nil
}

func (gb *GraphBuilder) Build() (*ExecutionGraph, error) {
	roots := []*Node{}
	for _, n := range gb.nodeMap {
		if len(n.Parents) == 0 {
			roots = append(roots, n)
		}
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("Build: no roots found")
	}

	return &ExecutionGraph{Roots: roots}, nil
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
