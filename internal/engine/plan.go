package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/topology"
	"github.com/dominikbraun/graph"
	"github.com/google/uuid"
)

type Link struct {
	SourceNode *ExecutionNode
	TargetNode *ExecutionNode
	SourceAttr string
	TargetAttr string
}

type ExecutionNode struct {
	id         string
	Type       topology.NamespacedType
	Profiles   []string
	Attrs      []string
	Conditions []reader.Condition
	Reader     reader.Reader

	isRoot bool
}

func executionNodeHash(en *ExecutionNode) string {
	return en.id
}

type ExecutionPlan struct {
	g    graph.Graph[string, *ExecutionNode]
	root *ExecutionNode

	adjacencyMap map[string]map[string]graph.Edge[string]
}

func newExecutionPlan() *ExecutionPlan {
	return &ExecutionPlan{
		g: graph.New(executionNodeHash, graph.Acyclic(), graph.Directed()),
	}
}

func (ep *ExecutionPlan) setRoot(r *ExecutionNode) error {
	if ep.root != nil {
		return fmt.Errorf("SetRoot: root node is already set")
	}
	r.id = "ROOT"
	ep.root = r
	return ep.g.AddVertex(r)
}

func (ep *ExecutionPlan) addNode(en *ExecutionNode) error {
	if ep.root == nil {
		return fmt.Errorf("AddNode: root node is not set")
	}
	en.id = uuid.New().String()
	return ep.g.AddVertex(en)
}

func (ep *ExecutionPlan) addEdge(from, to *ExecutionNode, attributes map[string]string) error {
	return ep.g.AddEdge(from.id, to.id, graph.EdgeAttributes(attributes))
}

func (ep *ExecutionPlan) edgesFrom(node *ExecutionNode) []graph.Edge[string] {
	if ep.adjacencyMap == nil {
		mp, err := ep.g.AdjacencyMap()
		if err != nil {
			ep.adjacencyMap = make(map[string]map[string]graph.Edge[string])
		}
		ep.adjacencyMap = mp
	}
	edgeMap, ok := ep.adjacencyMap[node.id]
	if !ok {
		return []graph.Edge[string]{}
	}
	edges := make([]graph.Edge[string], 0, len(edgeMap))
	i := 0
	for _, edge := range edgeMap {
		edges[i] = edge
		i++
	}
	return edges
}

func (ep *ExecutionPlan) GetLinks(sourceNode *ExecutionNode) ([]Link, error) {
	links := []Link{}
	edges := ep.edgesFrom(sourceNode)
	for _, edge := range edges {
		targetNode, err := ep.g.Vertex(edge.Target)
		if err != nil {
			return nil, fmt.Errorf("GetLinks: %w", err)
		}
		sourceAttr, ok := edge.Properties.Attributes["from"]
		if !ok {
			return nil, fmt.Errorf("GetLinks: edge has no \"from\" attribute")
		}
		targetAttr, ok := edge.Properties.Attributes["to"]
		if !ok {
			return nil, fmt.Errorf("GetLinks: edge has no \"to\" attribute")
		}
		links = append(links, Link{
			SourceNode: sourceNode,
			TargetNode: targetNode,
			SourceAttr: sourceAttr,
			TargetAttr: targetAttr,
		})
	}
	return links, nil
}
