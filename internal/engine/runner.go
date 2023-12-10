package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/reader"
)

type ExecutionRunner func(*ExecutionPlan) error

func ParallelExecutionRunner(ep *ExecutionPlan) error {
	layers, err := generateLayers(ep)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	for _, layer := range layers {
		err := executeLayer(ep, layer)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}
	}
	return nil
}

func generateLayers(ep *ExecutionPlan) ([][]*ExecutionNode, error) {
	if ep.root == nil {
		return nil, fmt.Errorf("generateLayers: root node is not set")
	}
	layers := [][]*ExecutionNode{}
	rootLayer := []*ExecutionNode{ep.root}
	lastLayer := rootLayer
	for len(lastLayer) > 0 {
		layers = append(layers, lastLayer)
		var newLayer []*ExecutionNode
		for _, node := range lastLayer {
			nextNodes, err := ep.Children(node)
			if err != nil {
				return nil, fmt.Errorf("generateLayers: %w", err)
			}
			newLayer = append(newLayer, nextNodes...)
		}
		lastLayer = newLayer
	}
	return layers, nil
}

func executeLayer(ep *ExecutionPlan, layer []*ExecutionNode) error {
	it := make(chan []*reader.Item, len(layer))
	e := make(chan error)
	for _, node := range layer {
		go executeNode(ep, node, it, e)
	}
	for i := 0; i < len(layer); i++ {
		select {
		case items := <-it:
			layer[i].Items = items
		case err := <-e:
			return fmt.Errorf("executeLayer: %w", err)
		}
	}
	return nil
}

func executeNode(ep *ExecutionPlan, node *ExecutionNode, it chan []*reader.Item, e chan error) {
	r := node.Reader
	err := ep.refreshNode(node)
	if err != nil {
		e <- fmt.Errorf("executeNode: %w", err)
	}
	if !node.IsExecutable() {
		it <- nil
		return
	}
	items, err := r.GetItems(node.Type.Type, node.Profiles, node.Attrs, node.Conditions)
	if err != nil {
		e <- fmt.Errorf("executeNode: %w", err)
	}
	it <- items
}
