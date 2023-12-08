package aws

import (
	"fmt"
	"slices"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/reader"
)

func NewReader(profiles []string, createPool client.CreatePool) (*AwsReader, error) {
	return &AwsReader{
		profiles: profiles,
		pool:     createPool(profiles),
	}, nil
}

type AwsReader struct {
	profiles []string
	pool     client.Pool
}

func (ar AwsReader) Name() string {
	return "aws"
}

func (ar AwsReader) GetSupportedProfiles() []string {
	return ar.profiles
}

func (ar AwsReader) GetItems(itemType string, profiles []string, ids []string, filters []reader.Filter, transforms []reader.Transform) ([]*reader.Item, error) {
	baseItems, err := ar.listBaseItemsForEachProfile(profiles, itemType)
	if err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	if len(ids) > 0 {
		baseItems = filterByIds(baseItems, ids)
	}
	items, err := ar.getItems(baseItems, itemType, filters, transforms)
	if err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	return common.Filter(items, func(r *reader.Item) bool { return r != nil }), nil
}

func (ar AwsReader) listBaseItemsForEachProfile(profiles []string, itemType string) ([]*reader.Item, error) {
	r := make(chan []*reader.Item, len(profiles))
	e := make(chan error)
	allResults := make([]*reader.Item, 0)
	for _, p := range profiles {
		go ar.listBaseItems(p, itemType, r, e)
	}
	for i := 0; i < len(profiles); i++ {
		select {
		case results := <-r:
			allResults = append(allResults, results...)
		case err := <-e:
			if err != nil {
				return nil, fmt.Errorf("listBaseItemsForEachProfile: %w", err)
			}
		}
	}
	return allResults, nil
}

func (ar AwsReader) listBaseItems(profile string, itemType string, r chan []*reader.Item, e chan error) {
	results, err := ar.pool.ListResources(profile, itemType)
	if err != nil {
		e <- err
	} else {
		r <- results
	}
}

func (ar AwsReader) getItems(baseItems []*reader.Item, itemType string, filters []reader.Filter, transforms []reader.Transform) ([]*reader.Item, error) {
	r := make(chan *reader.Item, len(baseItems))
	e := make(chan error)
	results := make([]*reader.Item, len(baseItems))
	for _, baseItem := range baseItems {
		go ar.getItem(baseItem, itemType, filters, transforms, r, e)
	}
	for i := 0; i < len(baseItems); i++ {
		select {
		case result := <-r:
			results[i] = result
		case err := <-e:
			return nil, fmt.Errorf("getItems: %w", err)
		}
	}
	return results, nil
}

func (ar AwsReader) getItem(baseItem *reader.Item, itemType string, filters []reader.Filter, transforms []reader.Transform, r chan *reader.Item, e chan error) {
	result, err := ar.pool.GetResource(baseItem.Profile, baseItem.Data.Identifier, itemType)
	if err != nil {
		e <- err
	} else if filterItem(result, filters) {
		result = transformItem(result, transforms)
		r <- result
	} else {
		r <- nil
	}
}

func filterByIds(items []*reader.Item, ids []string) []*reader.Item {
	result := make([]*reader.Item, 0)
	for _, item := range items {
		if slices.Contains(ids, item.Data.Identifier) {
			result = append(result, item)
		}
	}
	return result
}

func filterItem(item *reader.Item, filters []reader.Filter) bool {
	matches := true
	for _, filter := range filters {
		matches = matches && filter(item)
	}
	return matches
}

func transformItem(item *reader.Item, transforms []reader.Transform) *reader.Item {
	for _, transform := range transforms {
		item.Data.Properties = transform(item.Data.Properties)
	}
	return item
}
