package aws

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/bondyra/wtf/internal/reader"
)

var (
	awsConfigRegex      = regexp.MustCompile(`\[profile ([^\s]+)\]`)
	awsCredentialsRegex = regexp.MustCompile(`\[(.+)\]`)
	awsDefaultRegex     = regexp.MustCompile(`\[default\]`)
)

type AwsReader struct {
	profileNames []string
}

type ProfileProvider interface {
	ReadProfiles() ([]string, error)
}

type AwsConfigReader interface {
	GetPath() string
	ReadConfigAsString() (string, error)
}

type DefaultAwsConfigReader struct {
	path string
}

func (dacr DefaultAwsConfigReader) GetPath() string {
	return dacr.path
}

func (dacr DefaultAwsConfigReader) ReadConfigAsString() (string, error) {
	if dacr.path == "" {
		return "", nil
	}
	data, err := os.ReadFile(dacr.path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type AwsProfileProvider interface {
	ProvideProfiles() ([]string, error)
}

type DefaultAwsProfileProvider struct {
	awsConfigReader    AwsConfigReader
	configRegex        regexp.Regexp
	configDefaultRegex regexp.Regexp
}

func (dapp DefaultAwsProfileProvider) ProvideProfiles() ([]string, error) {
	results := []string{}
	data, err := dapp.awsConfigReader.ReadConfigAsString()
	if err != nil {
		return nil, err
	}
	if dapp.configDefaultRegex.MatchString(data) {
		results = append(results, "default")
	}
	for _, match := range dapp.configRegex.FindAllStringSubmatch(data, -1) {
		results = append(results, match[1])
	}
	return results, nil
}

func NewReader(awsCredentialsPath string, awsConfigPath string) (*AwsReader, error) {
	credentialsProfiles, err := DefaultAwsProfileProvider{awsConfigReader: DefaultAwsConfigReader{path: awsCredentialsPath}, configRegex: *awsCredentialsRegex, configDefaultRegex: *awsDefaultRegex}.ProvideProfiles()
	if err != nil {
		return nil, err
	}
	configProfiles, err := DefaultAwsProfileProvider{awsConfigReader: DefaultAwsConfigReader{path: awsConfigPath}, configRegex: *awsConfigRegex, configDefaultRegex: *awsDefaultRegex}.ProvideProfiles()
	if err != nil {
		return nil, err
	}
	duplicatedProfiles := GetDuplicatedElements(append(credentialsProfiles, configProfiles...))
	if len(duplicatedProfiles) > 0 {
		return nil, errors.New(fmt.Sprintf("Cannot proceed, found duplicated profiles %v in AWS files", duplicatedProfiles))
	}
	return &AwsReader{profileNames: append(configProfiles, credentialsProfiles...)}, nil
}

func (r *AwsReader) QueryAllProfiles() ([]reader.Credentials, error) {
	return []reader.Credentials{}, nil
}

func (r *AwsReader) QueryProfiles(profiles []string) ([]reader.Credentials, error) {
	return []reader.Credentials{}, nil
}
