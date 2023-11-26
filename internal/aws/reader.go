package aws

import "github.com/bondyra/swamp/internal/reader"

func NewReader(profileFactory ProfileFactory, awsFactory AwsFactory, configPaths []string) (*AwsReader, error) {
	provider := profileFactory.NewProfileProvider()
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
	awsFactory     AwsFactory
	configProfiles []string
	clients        map[string]AwsClientInterface
}

func (ar *AwsReader) Init(selectedProfiles []string) error {
	if selectedProfiles == nil {
		selectedProfiles = ar.configProfiles
	}
	existingProfiles := Intersect(ar.configProfiles, selectedProfiles)
	createdClients := make(map[string]AwsClientInterface, 0)
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
	return []string{}
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
