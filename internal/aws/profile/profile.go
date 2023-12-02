package profile

import (
	"fmt"
	"os"
	"regexp"
)

type ConfigReader func(string) ([]string, error)

func FromConfigFiles(configPaths ...string) ([]string, error) {
	return provideProfilesFromConfig(defaultReadConfig, configPaths...)
}

func provideProfilesFromConfig(readConfig ConfigReader, configPaths ...string) ([]string, error) {
	resultsMap := make(map[string]bool)
	for _, configPath := range configPaths {
		profiles, err := readConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("provideProfilesFromConfig: %w", err)
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

func defaultReadConfig(path string) ([]string, error) {
	results := make([]string, 0)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("defaultReadConfig: %w", err)
	}
	data := string(bytes)
	regex := regexp.MustCompile(`([^[\s]+)\]`)
	for _, match := range regex.FindAllStringSubmatch(data, -1) {
		results = append(results, match[1])
	}
	return results, nil
}
