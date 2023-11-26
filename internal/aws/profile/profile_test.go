package profile

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
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

	_, err := dacr.ReadConfigAsString("path")

	if err == nil {
		t.Errorf("expected:\nnot found error\ngot:\n%v", err)
	}
}

type MockConfigReader struct {
	pathToContent map[string]string
}

func (mcr MockConfigReader) ReadConfigAsString(path string) (string, error) {
	return mcr.pathToContent[path], nil
}

type MockErrorConfigReader struct{}

func (mecr MockErrorConfigReader) ReadConfigAsString(path string) (string, error) {
	return "", errors.New("")
}

func multiline(args ...string) string {
	return strings.Join(args, "\n")
}

func TestProvideProfiles(t *testing.T) {
	tests := []struct {
		name             string
		pathToContent    map[string]string
		configRegex      regexp.Regexp
		expectedProfiles []string
	}{
		{
			name:             "test no paths",
			pathToContent:    map[string]string{},
			expectedProfiles: []string{},
		},
		{
			name: "test multiple paths",
			pathToContent: map[string]string{
				"path1":         multiline("a", "[default]", "b", "[profile p1]", "[profile in_both]", "c"),
				"path2":         multiline("[profile p2]", "a", "[default]", "b", "[profile in_both]", "c"),
				"path3":         multiline("[p3_1]", "a", "b", "[profile p3_2]"),
				"emptyFilePath": "",
			},
			expectedProfiles: []string{"default", "p1", "in_both", "p2", "p3_1", "p3_2"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dp := DefaultProvider{configReader: MockConfigReader{pathToContent: test.pathToContent}}
			paths := make([]string, 0, len(test.pathToContent))
			for k := range test.pathToContent {
				paths = append(paths, k)
			}

			profiles, err := dp.ProvideProfiles(paths...)

			slices.Sort(profiles)
			slices.Sort(test.expectedProfiles)
			if !cmp.Equal(test.expectedProfiles, profiles) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedProfiles, profiles)
			}
			if err != nil {
				t.Errorf("%s error occured: %v", test.name, err)
			}
		})
	}
}

func TestProvideProfilesError(t *testing.T) {
	dp := DefaultProvider{configReader: MockErrorConfigReader{}}

	_, err := dp.ProvideProfiles("path")

	if err == nil {
		t.Errorf("expected:\nsome error\ngot:\n%v", err)
	}
}
