package profile

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
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Provider interface {
	ProvideProfiles(string) ([]string, error)
}

type DefaultProvider struct {
	configReader ConfigReader
}

func (dp DefaultProvider) ProvideProfiles(path string) ([]string, error) {
	results := []string{}
	data, err := dp.configReader.ReadConfigAsString(path)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(`([^[\s]+)\]`)
	for _, match := range regex.FindAllStringSubmatch(data, -1) {
		results = append(results, match[1])
	}
	return results, nil
}
