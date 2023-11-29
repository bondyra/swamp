package profile

import (
	"fmt"
	"os"
	"regexp"
)

type ProfileProvider func() ([]string, error)
type ConfigReader func(string) ([]string, error)

func DefaultReadConfig(path string) ([]string, error) {
	results := make([]string, 0)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}
	data := string(bytes)
	regex := regexp.MustCompile(`([^[\s]+)\]`)
	for _, match := range regex.FindAllStringSubmatch(data, -1) {
		results = append(results, match[1])
	}
	return results, nil
}

func NewConfigFileProfileProvider(readConfig ConfigReader, configPaths ...string) ProfileProvider {
	return func() ([]string, error) {
		resultsMap := make(map[string]bool)
		for _, configPath := range configPaths {
			profiles, err := readConfig(configPath)
			if err != nil {
				return nil, fmt.Errorf("config file profile provider error: %w", err)
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
}
