package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/reader"
)

type ExecutionGraph struct {
	nodeMap map[string]*Node
}

type Node struct {
	Alias      string
	ReaderName string
	Type       string
	Profiles   []string
	Filters    []reader.Filter
	Attrs      []string

	parents  []*Node
	children []*Node

	items []*reader.Item
}

func (eg *ExecutionGraph) GetRoots() ([]*Node, error) {
	roots := []*Node{}
	for _, node := range eg.nodeMap {
		if len(node.parents) == 0 {
			roots = append(roots, node)
		}
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("GetRoots: no roots found")
	}
	return roots, nil
}

func (eg *ExecutionGraph) GetChildren(parent *Node) []*Node {
	return parent.children
}

func (eg *ExecutionGraph) GetParents(child *Node) []*Node {
	return child.parents
}
