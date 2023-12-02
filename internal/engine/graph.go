package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/reader"
)

type ExecutionGraph struct {
	nodeMap map[string]*Node
	edges   map[string]map[string]bool // map[from]to
}

type Node struct {
	Alias      string
	ReaderName string
	Type       string
	Profiles   []string
	Filters    []reader.Filter
	Attrs      []string

	items []*reader.Item
}

func (eg *ExecutionGraph) GetRoots() ([]*Node, error) {
	roots := []*Node{}
	childAliases := map[string]bool{}
	for _, edge := range eg.edges {
		for to := range edge {
			childAliases[to] = true
		}
	}
	for v := range eg.nodeMap {
		if _, exists := childAliases[v]; !exists {
			roots = append(roots, eg.nodeMap[v])
		}
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("GetRoots: no roots found")
	}
	return roots, nil
}

func (eg *ExecutionGraph) GetChildren(parent *Node) []*Node {
	children := []*Node{}
	for _, child := range eg.nodeMap {
		if _, exists := eg.edges[parent.Alias][child.Alias]; exists {
			children = append(children, child)
		}
	}
	return children
}

func (eg *ExecutionGraph) GetParents(child *Node) []*Node {
	parents := []*Node{}
	for _, parent := range eg.nodeMap {
		if _, exists := eg.edges[parent.Alias][child.Alias]; exists {
			parents = append(parents, parent)
		}
	}
	return parents
}
