package aws

import (
	"errors"
	"fmt"
	"testing"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/google/go-cmp/cmp"
)

type MockProfileProvider struct {
	pathToOutput map[string]output
}

type output struct {
	profiles []string
	err      error
}

func (mpp MockProfileProvider) ProvideProfiles(path string) ([]string, error) {
	if output, ok := mpp.pathToOutput[path]; ok {
		return output.profiles, output.err
	} else {
		panic(fmt.Sprintf("path: %v does not exist in setup %v", path, mpp.pathToOutput))
	}
}

type MockProfileFactory struct {
	mpp MockProfileProvider
}

func (mpf MockProfileFactory) NewProfileProvider() ProfileProvider {
	return mpf.mpp
}

type MockAwsFactory struct{}

func (maf MockAwsFactory) NewClient(profile string) (AwsClientInterface, error) {
	return MockAwsClient{}, nil
}

type MockAwsClient struct{}

func (mac MockAwsClient) GetItem(id string, typeName string) (map[string]string, error) {
	return make(map[string]string, 0), nil
}
func (ac MockAwsClient) ListItems(typeName string) ([]map[string]string, error) {
	return make([]map[string]string, 0), nil
}

func TestNewReader(t *testing.T) {
	tests := []struct {
		name             string
		pathToOutput     map[string]output
		expectedProfiles []string
		returnsErr       bool
	}{
		{
			name: "test no profiles",
			pathToOutput: map[string]output{
				"awsCredentialsPath": {[]string{}, nil},
				"awsConfigPath":      {[]string{}, nil},
			},
			expectedProfiles: []string{},
			returnsErr:       false,
		},
		{
			name: "test one path with profiles 1",
			pathToOutput: map[string]output{
				"awsCredentialsPath": {[]string{"p1"}, nil},
				"awsConfigPath":      {[]string{}, nil},
			},
			expectedProfiles: []string{"p1"},
			returnsErr:       false,
		},
		{
			name: "test one path with profiles 2",
			pathToOutput: map[string]output{
				"awsCredentialsPath": {[]string{}, nil},
				"awsConfigPath":      {[]string{"p1"}, nil},
			},
			expectedProfiles: []string{"p1"},
			returnsErr:       false,
		},
		{
			name: "test both paths with profiles",
			pathToOutput: map[string]output{
				"awsCredentialsPath": {[]string{"p1"}, nil},
				"awsConfigPath":      {[]string{"p2"}, nil},
			},
			expectedProfiles: []string{"p1", "p2"},
			returnsErr:       false,
		},
		{
			name: "test both paths with overlapping profiles",
			pathToOutput: map[string]output{
				"awsCredentialsPath": {[]string{"p1", "p2"}, nil},
				"awsConfigPath":      {[]string{"p2"}, nil},
			},
			expectedProfiles: []string{"p1", "p2"},
			returnsErr:       false,
		},
		{
			name: "test return error when any provider returns error",
			pathToOutput: map[string]output{
				"awsCredentialsPath": {[]string{"anything"}, errors.New("some error")},
				"awsConfigPath":      {[]string{"p2"}, nil},
			},
			expectedProfiles: nil,
			returnsErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profileFactory := MockProfileFactory{mpp: MockProfileProvider{pathToOutput: test.pathToOutput}}
			awsFactory := MockAwsFactory{}

			reader, err := NewReader(profileFactory, awsFactory, maps.Keys(test.pathToOutput))

			if test.returnsErr {
				if reader != nil {
					t.Errorf("%s expected no reader to return", test.name)
				}
				if err == nil {
					t.Errorf("expected:\nerror\ngot:\n%v", err)
				}
			} else {
				slices.Sort(reader.configProfiles)
				slices.Sort(test.expectedProfiles)
				if !cmp.Equal(reader.configProfiles, test.expectedProfiles) {
					t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedProfiles, reader.configProfiles)
				}
				if err != nil {
					t.Errorf("%s error occured: %v", test.name, err)
				}
			}
		})
	}
}
