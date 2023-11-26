package aws

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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

func (mpf MockProfileFactory) NewProvider() profile.Provider {
	return mpf.mpp
}

type MockAwsFactory struct{}

func (maf MockAwsFactory) NewClient(profile string) (client.AwsClientInterface, error) {
	return MockAwsClient{}, nil
}

type MockErrorAwsFactory struct{}

func (meaf MockErrorAwsFactory) NewClient(profile string) (client.AwsClientInterface, error) {
	return nil, errors.New("some error")
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
				"path1": {[]string{}, nil},
				"path2": {[]string{}, nil},
			},
			expectedProfiles: []string{},
			returnsErr:       false,
		},
		{
			name: "test one path with profiles 1",
			pathToOutput: map[string]output{
				"path1": {[]string{"p1"}, nil},
				"path2": {[]string{}, nil},
			},
			expectedProfiles: []string{"p1"},
			returnsErr:       false,
		},
		{
			name: "test one path with profiles 2",
			pathToOutput: map[string]output{
				"path1": {[]string{}, nil},
				"path2": {[]string{"p1"}, nil},
			},
			expectedProfiles: []string{"p1"},
			returnsErr:       false,
		},
		{
			name: "test both paths with profiles",
			pathToOutput: map[string]output{
				"path1": {[]string{"p1"}, nil},
				"path2": {[]string{"p2"}, nil},
			},
			expectedProfiles: []string{"p1", "p2"},
			returnsErr:       false,
		},
		{
			name: "test both paths with overlapping profiles",
			pathToOutput: map[string]output{
				"path1": {[]string{"p1", "p2"}, nil},
				"path2": {[]string{"p2"}, nil},
			},
			expectedProfiles: []string{"p1", "p2"},
			returnsErr:       false,
		},
		{
			name: "test return error when any provider returns error",
			pathToOutput: map[string]output{
				"path1": {[]string{"anything"}, errors.New("some error")},
				"path2": {[]string{"p2"}, nil},
			},
			expectedProfiles: nil,
			returnsErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profileFactory := MockProfileFactory{mpp: MockProfileProvider{pathToOutput: test.pathToOutput}}
			awsFactory := MockAwsFactory{}

			reader, err := NewReader(profileFactory, awsFactory, definition.DefaultFactory{}, maps.Keys(test.pathToOutput))

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

func TestInit(t *testing.T) {

	tests := []struct {
		name                   string
		allProfiles            []string
		factory                client.Factory
		profilesToInit         []string
		expectedClientProfiles []string
		returnsErr             bool
	}{
		{
			name:                   "test init profiles",
			allProfiles:            []string{"p1", "p2", "p3", "p4"},
			factory:                MockAwsFactory{},
			profilesToInit:         []string{"p1", "p3"},
			expectedClientProfiles: []string{"p1", "p3"},
			returnsErr:             false,
		},
		{
			name:                   "test init all profiles",
			allProfiles:            []string{"p1", "p2", "p3", "p4"},
			factory:                MockAwsFactory{},
			profilesToInit:         nil,
			expectedClientProfiles: []string{"p1", "p2", "p3", "p4"},
			returnsErr:             false,
		},
		{
			name:                   "test init not existing profiles are filtered out",
			allProfiles:            []string{"p1", "p2"},
			factory:                MockAwsFactory{},
			profilesToInit:         []string{"p1", "p2", "p3", "p4"},
			expectedClientProfiles: []string{"p1", "p2"},
			returnsErr:             false,
		},
		{
			name:                   "test init errs when factory errs",
			allProfiles:            []string{"p1", "p2", "p3", "p4"},
			factory:                MockErrorAwsFactory{},
			profilesToInit:         nil,
			expectedClientProfiles: nil,
			returnsErr:             true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := AwsReader{test.factory, &definition.Definition{}, test.allProfiles, make(map[string]client.AwsClientInterface, 0)}
			expectedClients := make(map[string]client.AwsClientInterface, 0)
			for _, k := range test.expectedClientProfiles {
				expectedClients[k] = MockAwsClient{}
			}

			err := r.Init(test.profilesToInit)

			if test.returnsErr {
				if len(r.clients) > 0 {
					t.Errorf("%s expected no clients to be created", test.name)
				}
				if err == nil {
					t.Errorf("expected:\nerror\ngot:\n%v", err)
				}
			} else {
				if !cmp.Equal(r.clients, expectedClients) {
					t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, expectedClients, r.clients)
				}
				if err != nil {
					t.Errorf("%s error occured: %v", test.name, err)
				}
			}
		})
	}

}

func TestAwsReader_GetReaderName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test",
			want: "aws",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{}
			if got := ar.GetReaderName(); got != tt.want {
				t.Errorf("AwsReader.GetReaderName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsReader_GetProfileNames(t *testing.T) {
	type fields struct {
		configProfiles []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "test",
			fields: fields{configProfiles: []string{"p1", "p2"}},
			want:   []string{"p1", "p2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{
				configProfiles: tt.fields.configProfiles,
			}
			if got := ar.GetProfileNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AwsReader.GetProfileNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
