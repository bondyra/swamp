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
	ProvideProfiles(...string) ([]string, error)
}

type DefaultProvider struct {
	configReader ConfigReader
}

func (dp DefaultProvider) ProvideProfiles(configPaths ...string) ([]string, error) {
	resultsMap := make(map[string]bool)
	for _, configPath := range configPaths {
		profiles, err := dp.provideProfiles(configPath)
		if err != nil {
			return nil, err
		}
		for _, p := range profiles {
			resultsMap[p] = true
		}
	}
	results := make([]string, 0, len(resultsMap))
	for r := range resultsMap {
		results = append(results, r)
	}
	return results, nil
}

func (dp DefaultProvider) provideProfiles(path string) ([]string, error) {
	results := make([]string, 0)
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
