package aws

import (
	"os"
	"regexp"
)

type ConfigReader interface {
	ReadConfigAsString(string) (string, error)
}

type AwsConfigReader struct {
}

func (dacr AwsConfigReader) ReadConfigAsString(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type ProfileProvider interface {
	ReadProfiles() ([]string, error)
}

type DefaultProfileProvider struct {
	configReader ConfigReader
}

func (dapp DefaultProfileProvider) ProvideProfiles(path string) ([]string, error) {
	results := []string{}
	data, err := dapp.configReader.ReadConfigAsString(path)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(`([^[\s]+)\]`)
	for _, match := range regex.FindAllStringSubmatch(data, -1) {
		results = append(results, match[1])
	}
	return results, nil
}
