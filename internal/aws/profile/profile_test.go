package profile

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func multiline(args ...string) string {
	return strings.Join(args, "\n")
}

func TestDefaultReadConfig(t *testing.T) {
	tests := []struct {
		name        string
		createFile  bool
		fileContent string
		want        []string
		wantErr     bool
	}{
		{
			name:        "test return nothing on empty file",
			createFile:  true,
			fileContent: "",
			want:        []string{},
			wantErr:     false,
		},
		{
			name:        "test success",
			createFile:  true,
			fileContent: multiline("[default]", "[profile p1]", "to ignore", "[profile p2]", "[p3]", "[profile invalid"),
			want:        []string{"default", "p1", "p2", "p3"},
			wantErr:     false,
		},
		{
			name:       "test error when file does not exist",
			createFile: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tempFile *os.File
			var path string = "not_existent_path123"
			if tt.createFile {
				tempFile, _ = os.CreateTemp(".", strings.Replace(tt.name, " ", "-", -1))
				tempFile.Write([]byte(tt.fileContent))
				path = tempFile.Name()
			}
			got, err := defaultReadConfig(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultReadConfig() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultReadConfig() = %v, want %v", got, tt.want)
			}
			if tt.createFile {
				tempFile.Close()
				os.Remove(tempFile.Name())
			}
		})
	}
}

func newMockConfigReader(pathToProfiles map[string][]string) ConfigReader {
	return func(path string) ([]string, error) {
		return pathToProfiles[path], nil
	}
}

func newMockErrorConfigReader() ConfigReader {
	return func(path string) ([]string, error) {
		return nil, errors.New("")
	}
}

func TestProvideProfilesFromConfig(t *testing.T) {
	tests := []struct {
		name             string
		mockConfigReader ConfigReader
		configPaths      []string
		want             []string
		wantErr          bool
	}{
		{
			name:             "test return nothing on nil paths",
			mockConfigReader: newMockConfigReader(map[string][]string{}),
			configPaths:      nil,
			want:             []string{},
			wantErr:          false,
		},
		{
			name:             "test return nothing on no outputs",
			mockConfigReader: newMockConfigReader(map[string][]string{}),
			configPaths:      []string{"something"},
			want:             []string{},
			wantErr:          false,
		},
		{
			name: "test full",
			mockConfigReader: newMockConfigReader(map[string][]string{
				"1": {"p1", "p2"},
				"2": {"p2", "p3"},
				"3": {"p3", "p4"},
			}),
			configPaths: []string{"1", "2", "3", "4"},
			want:        []string{"p1", "p2", "p3", "p4"},
			wantErr:     false,
		},
		{
			name:             "test error when reader errs",
			configPaths:      []string{"1"},
			mockConfigReader: newMockErrorConfigReader(),
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provideProfilesFromConfig(tt.mockConfigReader, tt.configPaths...)

			slices.Sort(tt.want)
			slices.Sort(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigFileProfileProvider() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigFileProfileProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
