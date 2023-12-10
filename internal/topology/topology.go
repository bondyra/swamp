package topology

import (
	"fmt"
	"slices"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/dominikbraun/graph"
)

type TopologyLoader func() (*Topology, error)

func ReaderTopologyLoader(readers []reader.Reader) TopologyLoader {
	return func() (*Topology, error) {
		itemSchemas := []*itemSchema{}
		linkSchemas := []*linkSchema{}
		for _, reader := range readers {
			itemSchemaPath := reader.GetItemSchemaPath()
			linkSchemaPath := reader.GetLinkSchemaPath()

			itemSchema, err := loadFromFile[itemSchema](itemSchemaPath)
			if err != nil {
				return nil, fmt.Errorf("ReaderTopologyLoader: %w", err)
			}
			linkSchema, err := loadFromFile[linkSchema](linkSchemaPath)
			if err != nil {
				return nil, fmt.Errorf("ReaderTopologyLoader: %w", err)
			}
			itemSchemas = append(itemSchemas, itemSchema)
			linkSchemas = append(linkSchemas, linkSchema)
		}
		return newDefaultTopology(itemSchemas, linkSchemas)
	}
}

type Topology struct {
	graph.Graph[NamespacedType, Node]

	typeToNamespaces map[string][]string
}

type Node struct {
	Type  NamespacedType
	Attrs []string
}

func nodeHash(n Node) NamespacedType {
	return n.Type
}

func newDefaultTopology(itemSchemas []*itemSchema, linkSchemas []*linkSchema) (*Topology, error) {
	t := Topology{Graph: graph.New(nodeHash), typeToNamespaces: make(map[string][]string, 0)}
	for _, schema := range itemSchemas {
		for _, item := range schema.Items {
			err := t.AddNode(
				Node{
					Type:  item.Type,
					Attrs: common.Map(item.Attrs, func(a attrJson) string { return a.Field }),
				},
			)
			if err != nil {
				return nil, fmt.Errorf("newTopology: %w", err)
			}
		}
	}
	for _, schema := range linkSchemas {
		for _, link := range schema.Links {
			err := t.AddEdge(link.From, link.To, linkMappings(link))
			if err != nil {
				return nil, fmt.Errorf("newTopology: %w", err)
			}
		}
	}
	return &t, nil
}

func linkMappings(link linkJson) func(*graph.EdgeProperties) {
	return graph.EdgeAttributes(map[string]string{
		"from": link.Mapping.From,
		"to":   link.Mapping.To,
	})
}

func (t *Topology) AddNode(node Node) error {
	t.typeToNamespaces[node.Type.Type] = append(t.typeToNamespaces[node.Type.Type], node.Type.Namespace)
	return t.Graph.AddVertex(node)
}

func (t *Topology) GetNamespacesForType(typ string) ([]string, error) {
	namespaces, ok := t.typeToNamespaces[typ]
	if !ok {
		return nil, fmt.Errorf("getNamespacesForType: type \"%s\" not found", typ)
	}
	return namespaces, nil
}

func (t *Topology) SupportsAttr(typ NamespacedType, attr string) bool {
	n, err := t.Graph.Vertex(typ)
	if err != nil {
		// todo add warn log
		return false
	}
	return slices.Contains(n.Attrs, attr)
}

func (t *Topology) GetAttrs(typ NamespacedType) []string {
	n, err := t.Graph.Vertex(typ)
	if err != nil {
		// todo add warn log
		return nil
	}
	return n.Attrs
}

func (t *Topology) ShortestPath(from NamespacedType, to NamespacedType) ([]NamespacedType, error) {
	path, err := graph.ShortestPath(t.Graph, from, to)
	if err != nil {
		return nil, fmt.Errorf("ShortestPath: %w", err)
	}
	return path, nil
}
