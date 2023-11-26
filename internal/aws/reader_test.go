package aws

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/google/go-cmp/cmp"
)

type mockProfileProviderOutput struct {
	profiles []string
	err      error
}

type mockProfileProvider struct {
	output mockProfileProviderOutput
}

func (mpp mockProfileProvider) ProvideProfiles(paths ...string) ([]string, error) {
	return mpp.output.profiles, mpp.output.err
}

type mockAwsFactoryOutput struct {
	client client.AwsClientInterface
	err    error
}

type mockAwsFactory struct {
	output mockAwsFactoryOutput
}

func (maf mockAwsFactory) NewClient(profile string) (client.AwsClientInterface, error) {
	return maf.output.client, maf.output.err
}

type mockAwsClient struct{}

func (mac mockAwsClient) GetItem(id string, typeName string) (*reader.ItemData, error) {
	return &reader.ItemData{}, nil
}

func (ac mockAwsClient) ListItems(typeName string) ([]*reader.ItemData, error) {
	return []*reader.ItemData{}, nil
}

type mockDefFactoryOutput struct {
	definition *definition.Definition
	err        error
}

type mockDefFactory struct {
	output mockDefFactoryOutput
}

func (mdf mockDefFactory) FromFile(jsonPath string) (*definition.Definition, error) {
	return mdf.output.definition, mdf.output.err
}

func TestNewReader(t *testing.T) {
	tests := []struct {
		name                                 string
		profileProviderProvideProfilesOutput []string
		profileProviderProvideProfilesError  error
		awsFactoryNewClientOutput            client.AwsClientInterface
		awsFactoryNewClientError             error
		defFactoryFromFileOutput             *definition.Definition
		defFactoryFromFileError              error
		wantErr                              bool
	}{
		{
			name:                                 "test no errors",
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  nil,
			awsFactoryNewClientOutput:            mockAwsClient{},
			awsFactoryNewClientError:             nil,
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              nil,
			wantErr:                              false,
		},
		{
			name:                                 "test profile provider error causes error",
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  errors.New("some error"),
			awsFactoryNewClientOutput:            mockAwsClient{},
			awsFactoryNewClientError:             nil,
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              nil,
			wantErr:                              true,
		},
		{
			name:                                 "test definition factory error causes error",
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  nil,
			awsFactoryNewClientOutput:            mockAwsClient{},
			awsFactoryNewClientError:             nil,
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              errors.New("some error"),
			wantErr:                              true,
		},
		{
			name:                                 "test aws factory error does not cause errors", // it is for lazy use
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  nil,
			awsFactoryNewClientOutput:            mockAwsClient{},
			awsFactoryNewClientError:             errors.New("some error"),
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              nil,
			wantErr:                              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileProvider := mockProfileProvider{mockProfileProviderOutput{tt.profileProviderProvideProfilesOutput, tt.profileProviderProvideProfilesError}}
			awsFactory := mockAwsFactory{mockAwsFactoryOutput{tt.awsFactoryNewClientOutput, tt.awsFactoryNewClientError}}
			defFactory := mockDefFactory{mockDefFactoryOutput{tt.defFactoryFromFileOutput, tt.defFactoryFromFileError}}

			got, err := NewReader(profileProvider, awsFactory, defFactory, []string{})

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewReader() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			expectedReader := &AwsReader{awsFactory, tt.defFactoryFromFileOutput, tt.profileProviderProvideProfilesOutput, nil}

			if !reflect.DeepEqual(got, expectedReader) {
				t.Errorf("NewReader() = %v, want %v", got, expectedReader)
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
			factory:                mockAwsFactory{mockAwsFactoryOutput{nil, nil}},
			profilesToInit:         []string{"p1", "p3"},
			expectedClientProfiles: []string{"p1", "p3"},
			returnsErr:             false,
		},
		{
			name:                   "test init all profiles",
			allProfiles:            []string{"p1", "p2", "p3", "p4"},
			factory:                mockAwsFactory{mockAwsFactoryOutput{nil, nil}},
			profilesToInit:         nil,
			expectedClientProfiles: []string{"p1", "p2", "p3", "p4"},
			returnsErr:             false,
		},
		{
			name:                   "test init not existing profiles are filtered out",
			allProfiles:            []string{"p1", "p2"},
			factory:                mockAwsFactory{mockAwsFactoryOutput{nil, nil}},
			profilesToInit:         []string{"p1", "p2", "p3", "p4"},
			expectedClientProfiles: []string{"p1", "p2"},
			returnsErr:             false,
		},
		{
			name:                   "test init errs when factory errs",
			allProfiles:            []string{"p1", "p2", "p3", "p4"},
			factory:                mockAwsFactory{mockAwsFactoryOutput{nil, errors.New("some error")}},
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
				expectedClients[k] = nil
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

func TestAwsReader_GetItems(t *testing.T) {
	type fields struct {
		awsFactory     client.Factory
		definition     *definition.Definition
		configProfiles []string
		clients        map[string]client.AwsClientInterface
	}
	type args struct {
		resourceType  string
		attrs         []string
		filter        reader.Filter
		parentContext reader.ParentContext
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []reader.ItemData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{
				awsFactory:     tt.fields.awsFactory,
				definition:     tt.fields.definition,
				configProfiles: tt.fields.configProfiles,
				clients:        tt.fields.clients,
			}
			got, err := ar.GetItems(tt.args.resourceType, tt.args.attrs, tt.args.filter, tt.args.parentContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("AwsReader.GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AwsReader.GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
