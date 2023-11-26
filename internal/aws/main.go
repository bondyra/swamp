package aws

import (
	"os"
	"regexp"

	"github.com/bondyra/wtf/internal/reader"
)

type AwsReader struct {
	profileNames []string
}

type ProfileProvider interface {
	ReadProfiles() ([]string, error)
}

type AwsConfigReader interface {
	ReadConfigAsString(string) (string, error)
}

type DefaultAwsConfigReader struct {
}

func (dacr DefaultAwsConfigReader) ReadConfigAsString(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type AwsProfileProvider interface {
	ProvideProfiles() ([]string, error)
}

type DefaultAwsProfileProvider struct {
	awsConfigReader AwsConfigReader
}

func (dapp DefaultAwsProfileProvider) ProvideProfiles(path string) ([]string, error) {
	results := []string{}
	data, err := dapp.awsConfigReader.ReadConfigAsString(path)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(`([^[\s]+)\]`)
	for _, match := range regex.FindAllStringSubmatch(data, -1) {
		results = append(results, match[1])
	}
	return results, nil
}

func NewReader(awsCredentialsPath string, awsConfigPath string) (*AwsReader, error) {
	provider := DefaultAwsProfileProvider{awsConfigReader: DefaultAwsConfigReader{}}
	credentialsProfiles, err := provider.ProvideProfiles(awsCredentialsPath)
	if err != nil {
		return nil, err
	}
	configProfiles, err := provider.ProvideProfiles(awsConfigPath)
	if err != nil {
		return nil, err
	}
	return &AwsReader{profileNames: RemoveDuplicates(append(configProfiles, credentialsProfiles...))}, nil
}

func (r *AwsReader) QueryAllProfiles() ([]reader.Credentials, error) {
	return []reader.Credentials{}, nil
}

func (r *AwsReader) QueryProfiles(profiles []string) ([]reader.Credentials, error) {
	return []reader.Credentials{}, nil
}
