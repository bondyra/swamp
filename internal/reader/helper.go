package reader

import (
	"slices"

	"github.com/bondyra/swamp/internal/common"
)

type InlineFilter func(*Item) bool

func CreateInlineIdFilter(conditions []Condition) InlineFilter {
	matchingIds := make(map[string]bool, 0)
	notMatchingIds := make(map[string]bool, 0)
	for _, c := range conditions {
		if c.Attr == "id" {
			switch c.Op {
			case common.EqualsTo:
				matchingIds[c.Value] = true
			case common.NotEqualsTo:
				notMatchingIds[c.Value] = true
			}
		}
	}
	return func(it *Item) bool {
		result := true
		if len(matchingIds) > 0 {
			result = result && matchingIds[it.Data.Identifier]
		}
		if len(notMatchingIds) > 0 {
			result = result && !notMatchingIds[it.Data.Identifier]
		}
		return result
	}
}

func CreateInlineFilter(conditions []Condition) InlineFilter {
	eq := make(map[string][]string, 0)
	neq := make(map[string][]string, 0)
	for _, c := range conditions {
		switch c.Op {
		case common.EqualsTo:
			eq[c.Attr] = append(eq[c.Attr], c.Value)
		case common.NotEqualsTo:
			neq[c.Attr] = append(neq[c.Attr], c.Value)
		}
	}
	return func(it *Item) bool {
		result := true
		for attr, val := range eq {
			if v, ok := (*it.Data.Properties)[attr]; ok {
				result = result && slices.Contains(val, v)
			}
		}
		for attr, val := range neq {
			if v, ok := (*it.Data.Properties)[attr]; ok {
				result = result && !slices.Contains(val, v)
			}
		}
		return result
	}
}

type InlineTransformer func(*Properties) *Properties

func CreateInlineTransformer(attrs []string) InlineTransformer {
	return func(p *Properties) *Properties {
		if len(attrs) == 0 {
			return p
		}
		result := make(Properties)
		for _, attr := range attrs {
			if val, ok := (*p)[attr]; ok {
				result[attr] = val
			}
		}
		return &result
	}
}
