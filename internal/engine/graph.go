package engine

import (
	"github.com/bondyra/swamp/internal/reader"
)

type ExecutionGraph struct {
	Roots []*Node
}

type Node struct {
	Alias    string
	Type     string
	Profiles []string
	Filters  []reader.Filter
	Attrs    []string

	Reader reader.Reader
	Items  []*reader.Item

	Parents  []*Node
	Children []*Node
}
