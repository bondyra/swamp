package aws

import (
	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/reader"
)

func NewReader(profileFactory profile.Factory, awsFactory client.Factory, configPaths []string) (*AwsReader, error) {
	provider := profileFactory.NewProvider()
	profilesLists := [][]string{}
	for _, configPath := range configPaths {
		profiles, err := provider.ProvideProfiles(configPath)
		if err != nil {
			return nil, err
		}
		profilesLists = append(profilesLists, profiles)
	}
	return &AwsReader{
		awsFactory:     awsFactory,
		configProfiles: Sum(profilesLists...),
	}, nil
}

type AwsReader struct {
	awsFactory     client.Factory
	configProfiles []string
	clients        map[string]client.AwsClientInterface
}

func (ar *AwsReader) Init(selectedProfiles []string) error {
	if selectedProfiles == nil {
		selectedProfiles = ar.configProfiles
	}
	existingProfiles := Intersect(ar.configProfiles, selectedProfiles)
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

func (ar AwsReader) GetItemNames() []string {
	return []string{}
}

func (ar AwsReader) GetDefaultItemAttributes(itemType string) []string {
	return []string{}
}

func (ar AwsReader) GetAllItemAttributes(itemType string) []string {
	return []string{}
}

func (ar AwsReader) GetItems(resourceType string, attrs []string, filter reader.Filter, parentContext reader.ParentContext) ([]reader.ItemData, error) {
	return []reader.ItemData{}, nil
}
