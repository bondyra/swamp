package aws

import (
	"fmt"
	"regexp"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/reader"
	"golang.org/x/exp/slices"
)

func NewReader(profiles []string, createPool client.CreatePool, definition *definition.Definition) *AwsReader {
	return &AwsReader{
		profiles: profiles,
		pool:     createPool(profiles),
		def:      definition,
	}
}

type AwsReader struct {
	profiles     []string
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

func (ar AwsReader) GetSupportedProfiles() []string {
	return ar.profiles
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

func (ar AwsReader) GetItems(itemType string, profiles []string, attrs []string, filters []reader.Filter, parents []*reader.Item) ([]*reader.Item, error) {
	typeDefinition, err := ar.typeDefinition(itemType)
	if err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	filterFunc, err := ar.getFilterFunc(typeDefinition, filters, parents)
	if err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	baseItems, err := ar.listBaseItemsForEachProfile(profiles, itemType)
	if err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	baseItems = common.Filter(baseItems, filterFunc)
	items, err := ar.getItems(baseItems, itemType, filterFunc)
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

func (ar AwsReader) getItems(baseItems []*reader.Item, itemType string, filterFunc func(*reader.Item) bool) ([]*reader.Item, error) {
	r := make(chan *reader.Item, len(baseItems))
	e := make(chan error)
	results := make([]*reader.Item, len(baseItems))
	for _, baseItem := range baseItems {
		go ar.getItem(baseItem, itemType, filterFunc, r, e)
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

func (ar AwsReader) getItem(baseItem *reader.Item, itemType string, filterFunc func(*reader.Item) bool, r chan *reader.Item, e chan error) {
	result, err := ar.pool.GetResource(baseItem.Profile, baseItem.Data.Identifier, itemType)
	if err != nil {
		e <- err
	} else if filterFunc(result) {
		r <- result
	} else {
		r <- nil
	}
}

func (ar AwsReader) getFilterFunc(td *definition.TypeDefinition, filters []reader.Filter, parents []*reader.Item) (func(*reader.Item) bool, error) {
	// todo: add support for parent filters
	funcChain := []func(*reader.Item) bool{}
	getIdValue := func(r *reader.Item) *string { return &r.Data.Identifier }
	attrNames := ar.getAllAttrNames(td)
	for _, filter := range filters {
		var getAttrValue func(r *reader.Item) *string
		attr := filter.Attr
		value := filter.Value
		if attr == td.IdentifierField || attr == td.Alias {
			getAttrValue = getIdValue
		} else {
			_, supported := attrNames[attr]
			if !supported {
				return nil, fmt.Errorf("getFilterFunc: invalid attribute to filter : %v, supported ones are %v", attr, attrNames)
			}
			getAttrValue = func(r *reader.Item) *string {
				if r.Data.Properties != nil {
					attrValue := (*r.Data.Properties)[attr]
					return &attrValue
				}
				return nil
			}
		}
		var filterOp func(*string) bool
		filterFunc := func(r *reader.Item) bool {
			a := getAttrValue(r)
			if a != nil {
				return filterOp(a)
			}
			return true
		}
		regex := regexp.MustCompile(filter.Value)
		switch filter.Op {
		case reader.OpEquals:
			filterOp = func(v *string) bool { return *v == value }
		case reader.OpNotEquals:
			filterOp = func(v *string) bool { return *v != value }
		case reader.OpLike:
			filterOp = func(v *string) bool { return regex.MatchString(*v) }
		case reader.OpNotLike:
			filterOp = func(v *string) bool { return !regex.MatchString(*v) }
		default:
			return nil, fmt.Errorf("filter operation \"%v\" is not supported", filter.Op)
		}
		funcChain = append(funcChain, filterFunc)
	}
	return func(r *reader.Item) bool {
		result := true
		for _, f := range funcChain {
			result = result && f(r)
		}
		return result
	}, nil
}

func (ar AwsReader) getAllAttrNames(td *definition.TypeDefinition) map[string]bool {
	result := make(map[string]bool, len(td.Attrs))
	for _, a := range td.Attrs {
		result[a.Field] = true
	}
	return result
}
