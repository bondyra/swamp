package aws

import (
	"fmt"

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
		awsFactory:     awsFactory,
		def:            definition,
		configProfiles: profiles,
	}, nil
}

type AwsReader struct {
	awsFactory     client.PoolFactory
	def            *definition.Definition
	knownTypes     []string
	configProfiles []string
	clients        map[string]client.AwsClientInterface
}

func (ar AwsReader) Name() string {
	return "aws"
}

func (ar AwsReader) KnownTypes() []string {
	if ar.knownTypes == nil {
		ar.knownTypes = make([]string, 0)
		for _, typeDefinition := range ar.def.TypeDefinitions {
			ar.knownTypes = append(ar.knownTypes, typeDefinition.Type)
		}
	}
	return ar.knownTypes
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
	return slices.Contains(ar.KnownTypes(), itemType)
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

func (ar AwsReader) GetItems(itemType string, profiles []string, attrs []string, filter reader.Filter, parentContext reader.ParentContext) ([]*reader.Item, error) {
	typeDefinition, err := ar.typeDefinition(itemType)
	if err != nil {
		return nil, err
	}
	identifiers, err := ar.listIdentifiers(typeDefinition, filter)
	if err != nil {
		return nil, err
	}
	results := make([]*reader.Item, 0)
	for _, id := range identifiers {
		item, err := ar.getItem(id, attrs, filter)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (ar AwsReader) listIdentifiers(typeDefinition *definition.TypeDefinition, filter reader.Filter) ([]string, error) {

	return nil, nil
}

func (ar AwsReader) getItem(id string, attrs []string, filter reader.Filter) (*reader.Item, error) {
	return nil, nil
}
