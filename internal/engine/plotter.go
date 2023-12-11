package engine

import (
	"fmt"

	"github.com/bondyra/swamp/internal/common"
	"github.com/bondyra/swamp/internal/reader"
	"gopkg.in/yaml.v3"
)

type Plotter func(*ExecutionResult, common.Verbosity) error

func PrintYamlPlotter(er *ExecutionResult, v common.Verbosity) error {
	output, err := yaml.Marshal(resultToMap(er, v))
	if err != nil {
		return fmt.Errorf("PrintYamlPlotter: %w", err)
	}
	fmt.Print(string(output) + "\n")
	return nil
}

func resultToMap(er *ExecutionResult, v common.Verbosity) map[string]any {
	rootGroups := er.GetRootGroups()
	return map[string]any{"results": lst(groupToMap, rootGroups, v)}
}

func groupToMap(ng *ResultGroup, v common.Verbosity) any {
	result := map[string]any{
		"type": ng.Type.String(),
	}
	result["items"] = lst(itemToMap, ng.Items, v)

	result["link"] = lst(func(c reader.Condition, v common.Verbosity) any { return c.String() }, ng.LinkConditions, v)
	if ng.LinkError != nil {
		result["linkError"] = ng.LinkError.Error()
	}
	return result
}

func itemToMap(ni *ResultItem, v common.Verbosity) any {
	result := map[string]any{
		"properties": ni.Item.Data.Properties,
	}
	if ni.QueryError != nil {
		result["queryError"] = ni.QueryError.Error()
	}
	if v == common.DebugVerbosity {
		result["id"] = ni.Id
		result["groupId"] = ni.GroupId
	}
	if len(ni.LinkErrors) > 0 && v == common.DebugVerbosity {
		result["linkErrors"] = lst(func(err error, v common.Verbosity) any { return err.Error() }, ni.LinkErrors, v)
	}
	return result
}

func lst[T any](f func(T, common.Verbosity) any, groups []T, v common.Verbosity) any {
	result := make([]any, len(groups))
	for i := range groups {
		result[i] = f(groups[i], v)
	}
	return result
}
