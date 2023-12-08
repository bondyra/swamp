package engine

import (
	"fmt"
	"strings"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/schema"
)

type GraphBuilder struct {
	profiles []string
	nodeMap  map[string]*Node
	s        schema.Schema
}

func NewGraphBuilder(s schema.Schema) *GraphBuilder {
	return &GraphBuilder{nodeMap: map[string]*Node{}, s: s}
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
			r, err := gb.getReader(ie.Value.Type, readers)
			if err != nil {
				return fmt.Errorf("ReadNodes: %w", err)
			}
			nodeType := gb.getNodeType(ie.Value.Type)
			transforms, err := gb.getNodeTransforms(r.Name(), nodeType, ie.Value.Modifiers)
			if err != nil {
				return fmt.Errorf("ReadNodes: %w", err)
			}
			filters, err := gb.getNodeFilters(r.Name(), nodeType, ie.Value.Modifiers)
			if err != nil {
				return fmt.Errorf("ReadNodes: %w", err)
			}
			gb.addNode(ie.Value.Alias, nodeType, gb.profiles, filters, transforms, r)
		}
	}
	return nil
}

func (gb *GraphBuilder) getReader(typ []string, readers []reader.Reader) (reader.Reader, error) {
	if len(typ) == 0 {
		return nil, fmt.Errorf("getReader: type is empty")
	}
	if len(typ) == 1 {
		// guessing reader as no namespace part is given
		t, err := gb.guessReader(typ[0], readers)
		if err != nil {
			return nil, fmt.Errorf("getReader: %w", err)
		}
		return t, nil
	}
	if len(typ) > 2 {
		return nil, fmt.Errorf("getReader: type is invalid") // should not normally happen unless parser is wrong
	}

	for _, r := range readers {
		if r.Name() == typ[0] && gb.s.IsTypeSupported(typ[0], typ[1]) {
			return r, nil
		}
	}
	return nil, fmt.Errorf("getReader: no reader found for type \"%s\"", typ)
}

func (gb *GraphBuilder) guessReader(typ string, readers []reader.Reader) (reader.Reader, error) {
	matchingReaders := []reader.Reader{}
	for _, r := range readers {
		if gb.s.IsTypeSupported(r.Name(), typ) {
			matchingReaders = append(matchingReaders, r)
		}
	}
	if len(matchingReaders) == 0 {
		return nil, fmt.Errorf("guessReader: no reader found for type %s", typ)
	}
	if len(matchingReaders) > 1 {
		return nil, fmt.Errorf(
			"guessReader: type %s is defined in multiple readers: %s",
			typ, common.Map(matchingReaders, func(r reader.Reader) string { return r.Name() }),
		)
	}
	// todo add debug log
	return matchingReaders[0], nil
}

func (gb *GraphBuilder) getNodeType(typ []string) string {
	return typ[len(typ)-1]
}

func (gb *GraphBuilder) getNodeTransforms(reader string, nodeType string, modifiers []language.Modifier) ([]reader.Transform, error) {
	return nil, nil // todo
}

func (gb *GraphBuilder) getNodeFilters(reader string, nodeType string, modifiers []language.Modifier) ([]reader.Filter, error) {
	return nil, nil // todo
}

func (gb *GraphBuilder) addNode(alias string, nodeType string, profiles []string, filters []reader.Filter, transforms []reader.Transform, r reader.Reader) error {
	if _, exists := gb.nodeMap[alias]; exists {
		return fmt.Errorf("addNode: alias %s already defined", alias)
	}
	// todo validate filters and attrs
	gb.nodeMap[alias] = &Node{Alias: alias, Reader: r, Type: nodeType, Profiles: profiles, Filters: filters, Transforms: []reader.Transform{}}
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
	fromReader := gb.nodeMap[parentAlias].Reader.Name()
	fromType := gb.nodeMap[parentAlias].Type
	toReader := gb.nodeMap[childAlias].Reader.Name()
	toType := gb.nodeMap[childAlias].Type
	if !gb.s.IsLinkSupported(fromReader, fromType, toReader, toType) {
		return fmt.Errorf("linkNode: link from \"%s\" to \"%s\" is not supported", parentAlias, childAlias)
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
