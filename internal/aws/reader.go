package aws

import (
	"fmt"
	"regexp"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/reader"
	"golang.org/x/exp/slices"
)

func NewReader(profileProvider profile.Provider, awsFactory client.PoolFactory, defFactory definition.Factory, configPaths []string) (*AwsReader, error) {
	profiles, err := profileProvider.ProvideProfiles(configPaths...)
	if err != nil {
		return nil, err
	}
	definition, err := defFactory.FromFile("definition.json")
	if err != nil {
		return nil, err
	}
	return &AwsReader{
		pool: awsFactory.NewPool(profiles),
		def:  definition,
	}, nil
}

type AwsReader struct {
	pool         client.Pool
	def          *definition.Definition
	knownTypes   []string
	knownAliases []string
}

func (ar AwsReader) Name() string {
	return "aws"
}

func (ar AwsReader) getKnownTypes() []string {
	if ar.knownTypes == nil {
		ar.knownTypes = make([]string, len(ar.def.TypeDefinitions))
		for i := range ar.def.TypeDefinitions {
			ar.knownTypes[i] = ar.def.TypeDefinitions[i].Type
		}
	}
	return ar.knownTypes
}

func (ar AwsReader) getKnownAliases() []string {
	if ar.knownAliases == nil {
		ar.knownAliases = make([]string, len(ar.def.TypeDefinitions))
		for i := range ar.def.TypeDefinitions {
			ar.knownAliases[i] = ar.def.TypeDefinitions[i].Alias
		}
	}
	return ar.knownAliases
}

func (ar AwsReader) typeDefinition(itemType string) (*definition.TypeDefinition, error) {
	for _, td := range ar.def.TypeDefinitions {
		if td.Type == itemType {
			return &td, nil
		}
	}
	return nil, fmt.Errorf("%v type is not supported", itemType)
}

func (ar AwsReader) IsTypeSupported(itemType string) bool {
	return slices.Contains(ar.getKnownTypes(), itemType) || slices.Contains(ar.getKnownAliases(), itemType)
}

func (ar AwsReader) IsLinkSupported(itemType string, parentType string) bool {
	td, err := ar.typeDefinition(itemType)
	if err != nil {
		return false
	}
	for _, l := range (*td).Parents {
		if l.Type == parentType {
			return true
		}
	}
	return false
}

func (ar AwsReader) AreAttrsSupported(itemType string, attrs []string) bool {
	td, err := ar.typeDefinition(itemType)
	if err != nil {
		return false
	}
	supported := common.Map((*td).Attrs, func(a definition.Attr) string { return a.Field })
	return len(common.Difference(attrs, supported)) == 0
}

func (ar AwsReader) IsFilterSupported(itemType string, filter reader.Filter) bool {
	td, err := ar.typeDefinition(itemType)
	if err != nil {
		return false
	}
	fields := common.Map((*td).Attrs, func(a definition.Attr) string { return a.Field })
	return slices.Contains(fields, filter.Attr)
}

func (ar AwsReader) GetItems(itemType string, profiles []string, attrs []string, filters []reader.Filter) ([]*reader.Item, error) {
	typeDefinition, err := ar.typeDefinition(itemType)
	if err != nil {
		return nil, err
	}
	baseItems, err := ar.listBaseItemsForEachProfile(profiles, itemType)
	if err != nil {
		return nil, err
	}
	baseItems, err = ar.maybeFilterIds(baseItems, filters, typeDefinition.IdentifierField)
	if err != nil {
		return nil, err
	}
	items, err := ar.getItems(baseItems, itemType)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (ar AwsReader) listBaseItemsForEachProfile(profiles []string, itemType string) ([]*reader.Item, error) {
	results := []*reader.Item{}
	for _, p := range profiles {
		r, err := ar.pool.ListResources(p, itemType)
		if err != nil {
			return nil, err
		}
		results = append(results, r...)
	}
	return results, nil
}

func (ar AwsReader) getItems(baseItems []*reader.Item, itemType string) ([]*reader.Item, error) {
	r := make(chan *reader.Item, len(baseItems))
	e := make(chan error)
	results := make([]*reader.Item, len(baseItems))
	for _, baseItem := range baseItems {
		go ar.getItem(baseItem, itemType, r, e)
	}
	for i := 0; i < len(baseItems); i++ {
		select {
		case result := <-r:
			results[i] = result
		case err := <-e:
			return nil, err
		}
	}
	return results, nil
}

func (ar AwsReader) getItem(baseItem *reader.Item, itemType string, r chan *reader.Item, e chan error) {
	result, err := ar.pool.GetResource(baseItem.Profile, baseItem.Data.Identifier, itemType)
	if err != nil {
		e <- err
	} else {
		r <- result
	}
}

func (ar AwsReader) maybeFilterIds(baseItems []*reader.Item, filters []reader.Filter, idField string) ([]*reader.Item, error) {
	if len(filters) == 0 {
		return baseItems, nil
	}
	idFilters := common.Filter(filters, func(f reader.Filter) bool { return f.Attr == idField })

	for _, idFilter := range idFilters {
		var filterFunc func(*reader.Item) bool
		regex := regexp.MustCompile(idFilter.Value)
		switch idFilter.Op {
		case reader.OpEquals:
			filterFunc = func(r *reader.Item) bool { return r.Data.Identifier == idFilter.Value }
		case reader.OpNotEquals:
			filterFunc = func(r *reader.Item) bool { return r.Data.Identifier != idFilter.Value }
		case reader.OpLike:
			filterFunc = func(r *reader.Item) bool { return regex.MatchString(r.Data.Identifier) }
		case reader.OpNotLike:
			filterFunc = func(r *reader.Item) bool { return !regex.MatchString(r.Data.Identifier) }
		default:
			return nil, fmt.Errorf("cannot filter identifier for %v", idFilter.Op)
		}
		baseItems = common.Filter(baseItems, filterFunc)
	}
	return baseItems, nil
}
