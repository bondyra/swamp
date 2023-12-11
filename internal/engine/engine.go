package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/common"
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

func (e *Engine) Run(ast *language.AST, readers map[string]reader.Reader, v common.Verbosity) error {
	topology, err := e.topologyLoader()
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	plan, err := e.planner(ast, topology, readers)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	result := e.runner(plan)
	err = e.plotter(result, v)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}
	return nil
}
