package engine

import (
	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/topology"
)

type ExecutionTask struct {
	Id             string
	Type           string
	Profiles       []string
	Attrs          []string
	BaseConditions []reader.Condition
	LinkConditions []reader.Condition
	Reader         reader.Reader

	executionNode      *ExecutionNode
	parentResultItemId string
}

type ExecutionTaskResult struct {
	Id    string
	Task  *ExecutionTask
	Items []*reader.Item
	Err   error
}

type ResultItem struct {
	Id         string
	Item       *reader.Item
	QueryError error
	LinkErrors []error

	GroupId string
}

type ResultGroup struct {
	Type           topology.NamespacedType
	Items          []*ResultItem
	LinkConditions []reader.Condition
	LinkError      error
}

// todo extract interface
type ExecutionResult struct {
	links        map[string]map[string]bool
	vertices     map[string]*ResultItem
	groups       map[string]*ResultGroup
	rootGroupIds []string
}

func newExecutionResult() *ExecutionResult {
	return &ExecutionResult{
		links:    make(map[string]map[string]bool),
		vertices: make(map[string]*ResultItem),
	}
}

func (er *ExecutionResult) LinkItemToGroup(sourceItemId, targetGroupId string) {
	if _, ok := er.links[sourceItemId]; !ok {
		er.links[sourceItemId] = make(map[string]bool)
	}
	er.links[sourceItemId][targetGroupId] = true
}

func (er *ExecutionResult) AddGroup(groupId string, group *ResultGroup) {
	er.groups[groupId] = group
}

func (er *ExecutionResult) AddRootGroup(groupId string, group *ResultGroup) {
	er.rootGroupIds = append(er.rootGroupIds, groupId)
	er.AddGroup(groupId, group)
}

func (er *ExecutionResult) GetRootGroups() []*ResultGroup {
	rootGroups := make([]*ResultGroup, 0, len(er.rootGroupIds))
	for i, root := range er.rootGroupIds {
		rootGroups[i] = er.groups[root]
	}
	return rootGroups
}

func (er *ExecutionResult) GetChildGroups(item *ResultItem) []*ResultGroup {
	children := make([]*ResultGroup, len(er.links[item.Id]))
	for groupId := range er.links[item.Id] {
		children = append(children, er.groups[groupId])
	}
	return children
}
