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
			dacr := AwsConfigReader{}

			content, err := dacr.ReadConfigAsString(tempFile.Name())

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
	dacr := AwsConfigReader{}

	content, err := dacr.ReadConfigAsString("path")

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

func (mcr MockConfigReader) ReadConfigAsString(path string) (string, error) {
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
			expectedProfiles: []string{"default"},
		},
		{
			name: "test custom profile only",
			configContent: `
			[profile p1]
			a
			b
			c
			d
			`,
			expectedProfiles: []string{"p1"},
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
			expectedProfiles: []string{"default", "p1"},
		},
		{
			name: "test default and custom two",
			configContent: `
			[default]
			a
			b
			[p1]
			c
			[profile abc]
			d
			`,
			expectedProfiles: []string{"default", "p1", "abc"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			dapp := DefaultProfileProvider{configReader: MockConfigReader{content: test.configContent}}

			profiles, err := dapp.ProvideProfiles("path")

			if !cmp.Equal(test.expectedProfiles, profiles) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedProfiles, profiles)
			}
			if err != nil {
				t.Errorf("%s error occured: %v", test.name, err)
			}
		})
	}
}
