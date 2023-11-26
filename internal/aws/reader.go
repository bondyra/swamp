package aws

import "github.com/bondyra/wtf/internal/reader"

func NewReader(awsCredentialsPath string, awsConfigPath string) (*AwsReader, error) {
	provider := DefaultProfileProvider{configReader: AwsConfigReader{}}
	credentialsProfiles, err := provider.ProvideProfiles(awsCredentialsPath)
	if err != nil {
		return nil, err
	}
	configProfiles, err := provider.ProvideProfiles(awsConfigPath)
	if err != nil {
		return nil, err
	}
	return &AwsReader{
		configProfileNames: RemoveDuplicates(append(configProfiles, credentialsProfiles...)),
		pool:               AwsConnectionPool{},
	}, nil
}

type AwsReader struct {
	configProfileNames []string
	pool               reader.ConnectionPool
}

func (acp AwsConnectionPool) Init(profileNames []string) error { return nil }

func (ar *AwsReader) Init(queryProfiles []string, all bool) {
	if all {
		ar.pool.Init(ar.configProfileNames)
	} else {
		ar.pool.Init(Intersect(ar.configProfileNames, queryProfiles))
	}
}
