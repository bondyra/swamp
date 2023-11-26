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

func (ar AwsReader) GetItems(itemType string, profiles []string, attrs []string, filter *reader.Filter, parentContext *reader.ParentContext) ([]*reader.Item, error) {
	typeDefinition, err := ar.typeDefinition(itemType)
	if err != nil {
		return nil, err
	}
	baseItems, err := ar.listBaseItemsForEachProfile(profiles, itemType)
	if err != nil {
		return nil, err
	}
	if filter != nil && filter.Attr == typeDefinition.IdentifierField {
		var filterFunc func(*reader.Item) bool
		regex := regexp.MustCompile(filter.Value)
		switch filter.Op {
		case reader.OpEquals:
			filterFunc = func(r *reader.Item) bool { return r.Data.Identifier == filter.Value }
		case reader.OpNotEquals:
			filterFunc = func(r *reader.Item) bool { return r.Data.Identifier != filter.Value }
		case reader.OpLike:
			filterFunc = func(r *reader.Item) bool { return regex.MatchString(r.Data.Identifier) }
		case reader.OpNotLike:
			filterFunc = func(r *reader.Item) bool { return !regex.MatchString(r.Data.Identifier) }
		default:
			return nil, fmt.Errorf("cannot filter identifier for %v", filter.Op)
		}
		baseItems = common.Filter(baseItems, filterFunc)
	}
	items, err := ar.getItems(baseItems, itemType)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (ar AwsReader) listBaseItemsForEachProfile(profiles []string, itemType string) ([]*reader.Item, error) {
	results := make([]*reader.Item, 0)
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
	results := make([]*reader.Item, len(baseItems))
	for i, baseItem := range baseItems {
		result, err := ar.pool.GetResource(baseItem.Profile, baseItem.Data.Identifier, itemType)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}
