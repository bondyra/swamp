package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/topology"
)

type Engine struct {
	topologyLoader topology.TopologyLoader
	planner        ExecutionPlanner
	runner         ExecutionRunner
	plotter        Plotter
}

func NewEngine(tl topology.TopologyLoader, p ExecutionPlanner, r ExecutionRunner, pl Plotter) *Engine {
	return &Engine{
		topologyLoader: tl,
		planner:        p,
		runner:         r,
		plotter:        pl,
	}
}

func (e *Engine) Run(ast *language.AST, readers map[string]reader.Reader) error {
	topology, err := e.topologyLoader()
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	ep, err := e.planner(ast, topology, readers)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	err = e.runner(ep)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	err = e.plotter(ep)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	return nil
}
