package aws

import (
	"fmt"
	"path"
	"runtime"

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

func (ar AwsReader) GetNamespace() string {
	return "aws"
}

func (ar AwsReader) GetItemSchemaPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename) + "/item_schema.json"
}

func (ar AwsReader) GetLinkSchemaPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename) + "/link_schema.json"
}

func (ar AwsReader) GetSupportedProfiles() []string {
	return ar.profiles
}

func (ar AwsReader) GetItems(itemType string, profiles []string, attrs []string, conditions []reader.Condition) ([]*reader.Item, error) {
	baseItems, err := ar.listBaseItemsForEachProfile(profiles, itemType)
	if err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	// optimization: filter baseItems by id before getting full items
	idFilter := reader.CreateInlineIdFilter(conditions)
	baseItems = common.Filter(baseItems, idFilter)
	//
	transformer := reader.CreateInlineTransformer(attrs)
	filter := reader.CreateInlineFilter(conditions)
	items, err := ar.getItems(baseItems, itemType, transformer, filter)
	if err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	return items, nil
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

func (ar AwsReader) getItems(
	baseItems []*reader.Item, itemType string,
	transform reader.InlineTransformer, filter reader.InlineFilter,
) ([]*reader.Item, error) {
	r := make(chan *reader.Item, len(baseItems))
	e := make(chan error)
	results := make([]*reader.Item, len(baseItems))
	for _, baseItem := range baseItems {
		go ar.getItem(baseItem, itemType, transform, filter, r, e)
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

func (ar AwsReader) getItem(
	baseItem *reader.Item, itemType string,
	transform reader.InlineTransformer, filter reader.InlineFilter,
	r chan *reader.Item, e chan error,
) {
	result, err := ar.pool.GetResource(baseItem.Profile, baseItem.Data.Identifier, itemType)
	if err != nil {
		e <- err
	} else if filter(result) {
		result.Data.Properties = transform(result.Data.Properties)
		r <- result
	} else {
		r <- nil
	}
}
