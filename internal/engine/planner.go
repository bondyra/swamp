package engine

import (
	"fmt"
	"slices"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/topology"
)

type ExecutionPlanner func(ast *language.AST, t *topology.Topology, readers map[string]reader.Reader) (*ExecutionPlan, error)

func DefaultExecutionPlanner(ast *language.AST, t *topology.Topology, readers map[string]reader.Reader) (*ExecutionPlan, error) {
	var err error
	allProfiles, err := getAllProfiles(readers)
	if err != nil {
		return nil, fmt.Errorf("DefaultExecutionPlanner: %w", err)
	}
	selectedProfiles, err := selectProfiles(ast, allProfiles)
	if err != nil {
		return nil, fmt.Errorf("DefaultExecutionPlanner: %w", err)
	}
	plan, err := createPlan(ast, t, readers, selectedProfiles)
	if err != nil {
		return nil, fmt.Errorf("DefaultExecutionPlanner: %w", err)
	}
	return plan, nil
}

func getAllProfiles(readers map[string]reader.Reader) ([]string, error) {
	var results []string
	for _, reader := range readers {
		profiles := reader.GetSupportedProfiles()
		results = append(results, profiles...)
	}
	duplicates := common.DuplicatedElements(results)
	if len(duplicates) > 0 {
		return nil, fmt.Errorf("getAllProfiles: duplicate profiles found: %v", duplicates)
	}
	return results, nil
}

func selectProfiles(ast *language.AST, profiles []string) ([]string, error) {
	if ast.Profiles.All || len(ast.Profiles.Elements) == 0 { // not providing "in" block means all profiles
		return profiles, nil
	}
	var results []string
	for _, profile := range ast.Profiles.Elements {
		if !slices.Contains(profiles, profile) {
			return nil, fmt.Errorf("selectProfiles: profile \"%s\" is not supported", profile)
		}
		results = append(results, profile)
	}
	return results, nil
}

func createPlan(ast *language.AST, t *topology.Topology, readers map[string]reader.Reader, profiles []string) (*ExecutionPlan, error) {
	plan := NewExecutionPlan()
	root, err := createWaypoint(plan, ast.Query.Root, t, readers, profiles)
	if err != nil {
		return nil, fmt.Errorf("createPlan: %w", err)
	}
	err = plan.SetRoot(root)
	if err != nil {
		return nil, fmt.Errorf("createPlan: %w", err)
	}
	prevWaypoint := root
	for _, li := range ast.Query.Sequence {
		nextWaypoint, err := createWaypoint(plan, li.Item, t, readers, profiles)
		if err != nil {
			return nil, fmt.Errorf("createPlan: %w", err)
		}
		err = plan.AddNode(nextWaypoint)
		if err != nil {
			return nil, fmt.Errorf("createPlan: %w", err)
		}
		err = createPath(plan, t, profiles, readers, prevWaypoint, nextWaypoint)
		if err != nil {
			return nil, fmt.Errorf("createPlan: %w", err)
		}
		prevWaypoint = nextWaypoint
	}
	return plan, nil
}

func createWaypoint(plan *ExecutionPlan, it language.Item, t *topology.Topology, readers map[string]reader.Reader, profiles []string) (*ExecutionNode, error) {
	namespacedType, err := readNodeType(it.Type, t, readers)
	if err != nil {
		return nil, fmt.Errorf("createWaypoint: %w", err)
	}
	conditions, err := readConditions(namespacedType, t, it.Modifiers)
	if err != nil {
		return nil, fmt.Errorf("createWaypoint: %w", err)
	}
	attrs, err := readAttrs(namespacedType, t, it.Modifiers)
	if err != nil {
		return nil, fmt.Errorf("createWaypoint: %w", err)
	}
	r, err := getReader(namespacedType, readers)
	if err != nil {
		return nil, fmt.Errorf("createWaypoint: %w", err)
	}
	return &ExecutionNode{
		Type:       namespacedType,
		Profiles:   profiles,
		Attrs:      attrs,
		Conditions: conditions,
		Reader:     r,
		Items:      make([]*reader.Item, 0),
	}, nil
}

func readNodeType(typ []string, t *topology.Topology, readers map[string]reader.Reader) (topology.NamespacedType, error) {
	if len(typ) != 1 && len(typ) != 2 {
		return topology.NamespacedType{}, fmt.Errorf("readNodeType: invalid type: %v", typ)
	}
	var nt topology.NamespacedType
	var err error
	if len(typ) == 1 {
		nt, err = guessNamespacedType(typ[0], t)
		if err != nil {
			return topology.NamespacedType{}, fmt.Errorf("readNodeType: %w", err)
		}
	} else {
		nt = topology.NamespacedType{
			Namespace: typ[0],
			Type:      typ[1],
		}
	}
	return nt, nil
}

func guessNamespacedType(typ string, t *topology.Topology) (topology.NamespacedType, error) {
	namespaces, err := t.GetNamespacesForType(typ)
	if err != nil {
		return topology.NamespacedType{}, fmt.Errorf("guessNamespacedType: %w", err)
	}
	if len(namespaces) == 0 {
		return topology.NamespacedType{}, fmt.Errorf("guessNamespacedType: type \"%s\" not found", typ)
	}
	if len(namespaces) > 1 {
		return topology.NamespacedType{}, fmt.Errorf("guessNamespacedType: type \"%s\" is ambiguous: %v", typ, namespaces)
	}
	return topology.NamespacedType{
		Namespace: namespaces[0],
		Type:      typ,
	}, nil
}

func readConditions(typ topology.NamespacedType, t *topology.Topology, modifiers []language.Modifier) ([]reader.Condition, error) {
	results := make([]reader.Condition, 0)
	for _, m := range modifiers {
		switch x := m.(type) {
		case language.SearchModifier:
			s := x.Value
			if !t.SupportsAttr(typ, s.Attr) {
				return nil, fmt.Errorf("readConditions: type \"%s\" does not support attribute \"%s\"", typ, s.Attr)
			}
			results = append(results, reader.Condition{
				Attr:  s.Attr,
				Op:    s.Op,
				Value: s.Value,
			})
		}
	}
	return results, nil
}

func readAttrs(typ topology.NamespacedType, t *topology.Topology, modifiers []language.Modifier) ([]string, error) {
	attrs := make([]string, 0)
	for _, m := range modifiers {
		switch x := m.(type) {
		case language.AttrModifier:
			if x.Value.All {
				return t.GetAttrs(typ), nil
			} else {
				for _, attr := range x.Value.Elements {
					if !t.SupportsAttr(typ, attr) {
						return nil, fmt.Errorf("readAttrs: type \"%s\" does not support attribute \"%s\"", typ, attr)
					}
					attrs = append(attrs, attr)
				}
			}
		}
	}
	return attrs, nil
}

func getReader(typ topology.NamespacedType, readers map[string]reader.Reader) (reader.Reader, error) {
	r, ok := readers[typ.Namespace]
	if !ok {
		return nil, fmt.Errorf("getReader: no reader for namespace %s", typ.Namespace)
	}
	return r, nil
}

func createPath(plan *ExecutionPlan, t *topology.Topology, profiles []string, readers map[string]reader.Reader, a *ExecutionNode, b *ExecutionNode) error {
	path, err := t.ShortestPath(a.Type, b.Type)
	if err != nil {
		return fmt.Errorf("createPath: %w", err)
	}
	previousNode := a
	for i, currentType := range path {
		if i == 0 { // skip a
			continue
		}
		var currentNode *ExecutionNode
		if i != len(path)-1 {
			r, err := getReader(currentType, readers)
			if err != nil {
				return fmt.Errorf("createPath: %w", err)
			}
			currentNode = &ExecutionNode{
				Type:     currentType,
				Profiles: profiles,
				Reader:   r,
			}
			err = plan.AddNode(currentNode)
			if err != nil {
				return fmt.Errorf("createPath: %w", err)
			}
		} else {
			currentNode = b
		}
		topologyEdge, err := t.Edge(previousNode.Type, currentNode.Type)
		if err != nil {
			return fmt.Errorf("createPath: %w", err)
		}
		err = plan.AddEdge(previousNode, currentNode, topologyEdge.Properties.Attributes)
		if err != nil {
			return fmt.Errorf("createPath: %w", err)
		}
		previousNode = currentNode
	}
	return nil
}
