package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/common"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/topology"
	"github.com/dominikbraun/graph"
	"github.com/google/uuid"
)

type ExecutionNode struct {
	id         string
	Type       topology.NamespacedType
	Profiles   []string
	Attrs      []string
	Conditions []reader.Condition
	Reader     reader.Reader

	tasks  []*ExecutionTask
	isRoot bool
}

type ExecutionTask struct {
	parentConditions []reader.Condition

	results []*reader.Item
}

func executionNodeHash(en *ExecutionNode) string {
	return en.id
}

type ExecutionPlan struct {
	g    graph.Graph[string, *ExecutionNode]
	root *ExecutionNode
}

func NewExecutionPlan() *ExecutionPlan {
	return &ExecutionPlan{
		g: graph.New(executionNodeHash, graph.Acyclic(), graph.Directed()),
	}
}

func (ep *ExecutionPlan) SetRoot(r *ExecutionNode) error {
	if ep.root != nil {
		return fmt.Errorf("SetRoot: root node is already set")
	}
	r.id = "ROOT"
	r.isRoot = true
	r.tasks = make([]*ExecutionTask, 0)
	ep.root = r
	return ep.g.AddVertex(r)
}

func (ep *ExecutionPlan) AddNode(en *ExecutionNode) error {
	if ep.root == nil {
		return fmt.Errorf("AddNode: root node is not set")
	}
	en.id = uuid.New().String()
	en.isRoot = false
	en.tasks = make([]*ExecutionTask, 0)
	return ep.g.AddVertex(en)
}

func (ep *ExecutionPlan) AddEdge(from, to *ExecutionNode, attributes map[string]string) error {
	return ep.g.AddEdge(from.id, to.id, graph.EdgeAttributes(attributes))
}

func (ep *ExecutionPlan) Children(node *ExecutionNode) ([]*ExecutionNode, error) {
	results := []*ExecutionNode{}
	edges, err := ep.g.Edges()
	if err != nil {
		return nil, fmt.Errorf("GetNextNodes: %w", err)
	}
	for _, e := range edges {
		if e.Source == node.id {
			target, err := ep.g.Vertex(e.Target)
			if err != nil {
				return nil, fmt.Errorf("GetNextNodes: %w", err)
			}
			results = append(results, target)
		}
	}
	return results, nil
}

func (ep *ExecutionPlan) GetRoot() *ExecutionNode {
	return ep.root
}

func (ep *ExecutionPlan) generateTasks(node *ExecutionNode) error {
	parentEdges, err := ep.getParentEdges(node)
	if err != nil {
		return fmt.Errorf("refreshNode: %w", err)
	}
	executionTasks := make([]*ExecutionTask, 0)
	for _, e := range parentEdges {
		parentNode, err := ep.g.Vertex(e.Source)
		if err != nil {
			return fmt.Errorf("refreshNode: %w", err)
		}
		attributes := e.Properties.Attributes
		parentAttr := attributes["from"]
		childAttr := attributes["to"]

		if len(parentNode.tasks) == 0 && !parentNode.isRoot {
			// stop condition - at least one parent in empty and it is not a root, so no point in querying the child
			// todo add logging
			return nil
		}
		for _, pt := range parentNode.tasks {
			parentAttrValue, ok := (*pi.Data.Properties)[parentAttr]
			if !ok {
				return fmt.Errorf("refreshNode: parent item does not have attribute \"%s\"", parentAttr)
			}
			linkConditions = append(linkConditions, reader.Condition{
				Attr: childAttr, Op: common.EqualsTo, Value: parentAttrValue,
			})
		}
	}
	node.Conditions = append(node.Conditions, linkConditions...)
	return nil
}

func (ep *ExecutionPlan) getParentEdges(target *ExecutionNode) ([]graph.Edge[string], error) {
	edges, err := ep.g.Edges()
	if err != nil {
		return nil, fmt.Errorf("getIncomingEdges: %w", err)
	}
	results := make([]graph.Edge[string], 0)
	for _, e := range edges {
		if e.Target == target.id {
			results = append(results, e)
		}
	}
	return results, nil
}
