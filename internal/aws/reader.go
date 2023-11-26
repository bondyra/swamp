package aws

import (
	"fmt"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/reader"
)

func NewReader(profileProvider profile.Provider, awsFactory client.Factory, defFactory definition.Factory, configPaths []string) (*AwsReader, error) {
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
		definition:     definition,
		configProfiles: profiles,
	}, nil
}

type AwsReader struct {
	awsFactory     client.Factory
	definition     *definition.Definition
	configProfiles []string
	clients        map[string]client.AwsClientInterface
}

func (ar *AwsReader) Init(selectedProfiles []string) error {
	if selectedProfiles == nil {
		selectedProfiles = ar.configProfiles
	}
	existingProfiles := common.Intersect(ar.configProfiles, selectedProfiles)
	createdClients := make(map[string]client.AwsClientInterface, 0)
	for _, profile := range existingProfiles {
		var err error
		createdClients[profile], err = ar.awsFactory.NewClient(profile)
		if err != nil {
			return err
		}
	}
	ar.clients = createdClients
	return nil
}

func (ar AwsReader) GetReaderName() string {
	return "aws"
}

func (ar AwsReader) GetProfileNames() []string {
	return ar.configProfiles
}

func (ar AwsReader) GetItemTypes() []string {
	return ar.definition.AllDefinedTypes()
}

func (ar AwsReader) GetDefaultItemAttributes(itemType string) []string {
	return ar.definition.GetAtributesForType(itemType, false)
}

func (ar AwsReader) GetAllItemAttributes(itemType string) []string {
	return ar.definition.GetAtributesForType(itemType, true)
}

func (ar AwsReader) GetItems(itemType string, attrs []string, filter reader.Filter, parentContext reader.ParentContext) ([]reader.ItemData, error) {
	if !ar.definition.SupportsType(itemType) {
		return nil, fmt.Errorf("unsupported type: %v", itemType)
	}
	if !ar.definition.SupportsAttrs(itemType, attrs) {
		return nil, fmt.Errorf("type %v does not support some of following attrs: %v", itemType, attrs)
	}
	if !ar.definition.SupportsFilter(itemType, filter) {
		return nil, fmt.Errorf("type %v does not support filter: %v", itemType, filter)
	}
	return nil, nil
}

func (ar AwsReader) listItems(itemType string) ([]*reader.ItemData, error) {
	return nil, nil
}
