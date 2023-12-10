package engine

import (
	"encoding/json"
	"fmt"

	"github.com/bondyra/swamp/internal/reader"
)

type Plotter func(*ExecutionPlan) error

func PrintJsonPlotter(ep *ExecutionPlan) error {
	output, err := json.MarshalIndent(executionToMap(ep), "", " ")
	if err != nil {
		return fmt.Errorf("PrintJsonPlotter: %w", err)
	}
	fmt.Print(string(output) + "\n")
	return nil
}

func executionToMap(ep *ExecutionPlan) map[string]any {
	rootMap := make([]any, 0)
	result := map[string]any{ep.root.Type.String(): rootMap}
	nodeToMap(ep, ep.root, &rootMap)
	return result
}

func nodeToMap(ep *ExecutionPlan, node *ExecutionNode, result *[]any) {
	nodeMap := make([]any, 0)
	result = &nodeMap
	for _, item := range node.Items {
		nodeMap = append(nodeMap, itemToMap(item))
	}
	children, err := ep.Children(node)
	if err != nil {
		panic(err)
	}
	for _, child := range children {
		childMap := make(map[string]any, 0)

		nodeToMap(ep, child, &nodeMap)
	}

}

func itemToMap(item *reader.Item) any {
	return *item.Data.Properties
}
