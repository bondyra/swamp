package aws

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	MockReadFile = func(path string) (string, error) {
		return fmt.Sprintf("content %v", path), nil
	}
	MockErrorReadFile = func(path string) (string, error) {
		return "", errors.New("test error")
	}
)

func TestGetPath(t *testing.T) {
	dacr := DefaultAwsConfigReader{path: "path"}

	actualPath := dacr.GetPath()

	if !cmp.Equal(dacr.path, actualPath) {
		t.Errorf("expected:\n%v\ngot:\n%v", dacr.path, actualPath)
	}
}

func TestReadConfigAsString(t *testing.T) {
	tests := []struct {
		name         string
		inputContent string
		isError      bool
	}{
		{
			name:         "test simple string",
			inputContent: "simple string",
		},
		{
			name:         "test multiline string",
			inputContent: "lorem\nipsum\ndolor sit\n123 amet@&$*@jdsadas\nj9dwjwq\n\ndsa\t321&*\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempFile, _ := os.CreateTemp(".", strings.Replace(test.name, " ", "-", -1))
			tempFile.Write([]byte(test.inputContent))
			dacr := DefaultAwsConfigReader{path: tempFile.Name()}

			content, err := dacr.ReadConfigAsString()

			if !cmp.Equal(test.inputContent, content) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.inputContent, content)
			}
			if err != nil {
				t.Errorf("%s error occured: %v", test.name, err)
			}

			tempFile.Close()
			os.Remove(tempFile.Name())
		})
	}
}
func TestReadConfigAsStringForNonExistentFile(t *testing.T) {
	dacr := DefaultAwsConfigReader{path: "path"}

	content, err := dacr.ReadConfigAsString()

	if content != "" {
		t.Errorf("expected: NOTHING\ngot:\n%v", content)
	}
	if err == nil {
		t.Errorf("expected:\nnot found error\ngot:\n%v", err)
	}
}

type MockConfigReader struct {
	content string
}

func (mcr MockConfigReader) GetPath() string {
	return "mock"
}

func (mcr MockConfigReader) ReadConfigAsString() (string, error) {
	return mcr.content, nil
}

func TestProvideProfiles(t *testing.T) {
	tests := []struct {
		name             string
		configContent    string
		configRegex      regexp.Regexp
		expectedProfiles []string
	}{
		{
			name:             "test nothing",
			configContent:    "",
			configRegex:      *awsConfigRegex,
			expectedProfiles: []string{},
		},
		{
			name: "test default profile only",
			configContent: `
			[default]
			a
			b
			c
			d
			`,
			configRegex:      *awsConfigRegex,
			expectedProfiles: []string{"default"},
		},
		{
			name: "test default and custom one",
			configContent: `
			[default]
			a
			b
			[profile p1]
			c
			d
			`,
			configRegex:      *awsConfigRegex,
			expectedProfiles: []string{"default", "p1"},
		}, // TODO
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			dapp := DefaultAwsProfileProvider{awsConfigReader: MockConfigReader{content: test.configContent}, configRegex: test.configRegex, configDefaultRegex: *awsDefaultRegex}

			profiles, err := dapp.ProvideProfiles()

			if !cmp.Equal(test.expectedProfiles, profiles) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedProfiles, profiles)
			}
			if err != nil {
				t.Errorf("%s error occured: %v", test.name, err)
			}
		})
	}
}
