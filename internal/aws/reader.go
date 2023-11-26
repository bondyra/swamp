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
		configProfiles: Intersect(profilesLists...),
	}, nil
}

type AwsReader struct {
	awsFactory     AwsFactory
	configProfiles []string
	connections    map[string]AwsConnection
}

func (ar *AwsReader) Init(selectedProfiles []string) error {
	return ar.initPool(selectedProfiles)
}

func (ar *AwsReader) InitAll() error {
	return ar.initPool(nil)
}

func (ar *AwsReader) initPool(selectedProfiles []string) error {
	existingProfiles := Intersect(ar.configProfiles, selectedProfiles)
	for _, profile := range existingProfiles {
		ar.connections[profile] = AwsConnection{profile: profile}
		err := ar.connections[profile].Init(ar.awsFactory)
		if err != nil {
			return err
		}
	}
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
