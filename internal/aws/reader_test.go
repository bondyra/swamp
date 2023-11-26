package aws

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/reader"
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

type mockPoolFactory struct{}

func (maf mockPoolFactory) NewPool(profiles []string) client.Pool {
	return nil
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
		defFactoryFromFileOutput             *definition.Definition
		defFactoryFromFileError              error
		wantErr                              bool
	}{
		{
			name:                                 "test no errors",
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  nil,
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              nil,
			wantErr:                              false,
		},
		{
			name:                                 "test profile provider error causes error",
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  errors.New("some error"),
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              nil,
			wantErr:                              true,
		},
		{
			name:                                 "test definition factory error causes error",
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  nil,
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              errors.New("some error"),
			wantErr:                              true,
		},
		{
			name:                                 "test aws factory error does not cause errors", // it is for lazy use
			profileProviderProvideProfilesOutput: []string{"p1", "p2"},
			profileProviderProvideProfilesError:  nil,
			defFactoryFromFileOutput:             &definition.Definition{},
			defFactoryFromFileError:              nil,
			wantErr:                              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileProvider := mockProfileProvider{mockProfileProviderOutput{tt.profileProviderProvideProfilesOutput, tt.profileProviderProvideProfilesError}}
			poolFactory := mockPoolFactory{}
			defFactory := mockDefFactory{mockDefFactoryOutput{tt.defFactoryFromFileOutput, tt.defFactoryFromFileError}}

			got, err := NewReader(profileProvider, poolFactory, defFactory, []string{})

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewReader() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			expectedReader := &AwsReader{nil, tt.defFactoryFromFileOutput, nil, nil}

			if !reflect.DeepEqual(got, expectedReader) {
				t.Errorf("NewReader() = %v, want %v", got, expectedReader)
			}
		})
	}
}

func TestAwsReader_Name(t *testing.T) {
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
			if got := ar.Name(); got != tt.want {
				t.Errorf("AwsReader.GetReaderName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsReader_IsTypeSupported(t *testing.T) {
	tests := []struct {
		name     string
		def      *definition.Definition
		itemType string
		want     bool
	}{
		{
			name: "test true when type matches",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
				{Type: "t2", IdentifierField: "id2", Alias: "a2"},
				{Type: "t3", IdentifierField: "id3", Alias: "a3"},
			}},
			itemType: "t3",
			want:     true,
		},
		{
			name: "test true when alias matches",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
				{Type: "t2", IdentifierField: "id2", Alias: "a2"},
				{Type: "t3", IdentifierField: "id3", Alias: "a3"},
			}},
			itemType: "a1",
			want:     true,
		},
		{
			name: "test false when no type matches",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t11", IdentifierField: "id1", Alias: "aa1"},
				{Type: "t22", IdentifierField: "id2", Alias: "aa2"},
				{Type: "t33", IdentifierField: "id3", Alias: "aa3"},
			}},
			itemType: "t1",
			want:     false,
		},
		{
			name: "test false when no alias matches",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t11", IdentifierField: "id1", Alias: "aa1"},
				{Type: "t22", IdentifierField: "id2", Alias: "aa2"},
				{Type: "t33", IdentifierField: "id3", Alias: "aa3"},
			}},
			itemType: "a1",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{def: tt.def}
			if got := ar.IsTypeSupported(tt.itemType); got != tt.want {
				t.Errorf("AwsReader.IsTypeSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsReader_IsLinkSupported(t *testing.T) {
	tests := []struct {
		name       string
		def        *definition.Definition
		itemType   string
		parentType string
		want       bool
	}{
		{
			name: "test false when there are no parents defined",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
			}},
			itemType:   "t1",
			parentType: "p1",
			want:       false,
		},
		{
			name: "test false when parent does not match",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
				{
					Type: "t2", IdentifierField: "id1", Alias: "a1",
					Parents: []definition.ParentDefinition{
						{Type: "p1", LinkType: "inline", Links: []definition.Link{{ParentField: "pf", Field: "f"}}},
					},
				},
			}},
			itemType:   "t2",
			parentType: "p11",
			want:       false,
		},
		{
			name: "test false when type does not exist",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Parents: []definition.ParentDefinition{
						{Type: "p1", LinkType: "inline", Links: []definition.Link{{ParentField: "pf", Field: "f"}}},
					},
				},
			}},
			itemType:   "t2",
			parentType: "p1",
			want:       false,
		},
		{
			name: "test true",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Parents: []definition.ParentDefinition{
						{Type: "p1", LinkType: "inline", Links: []definition.Link{{ParentField: "pf", Field: "f"}}},
						{Type: "p2", LinkType: "inline", Links: []definition.Link{{ParentField: "pf", Field: "f"}}},
					},
				},
			}},
			itemType:   "t1",
			parentType: "p2",
			want:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{def: tt.def}
			if got := ar.IsLinkSupported(tt.itemType, tt.parentType); got != tt.want {
				t.Errorf("AwsReader.IsLinkSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsReader_AreAttrsSupported(t *testing.T) {
	tests := []struct {
		name     string
		def      *definition.Definition
		itemType string
		attrs    []string
		want     bool
	}{
		{
			name: "test false when there are no attrs",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
			}},
			itemType: "t1",
			attrs:    []string{"a1"},
			want:     false,
		},
		{
			name: "test false when type does not exist",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Attrs: []definition.Attr{
						{Field: "f1"},
						{Field: "f2"},
						{Field: "f3"},
					},
				},
			}},
			itemType: "t111",
			attrs:    []string{"f1", "f2"},
			want:     false,
		},
		{
			name: "test false on non existent attrs",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Attrs: []definition.Attr{
						{Field: "f1"},
						{Field: "f2"},
						{Field: "f3"},
					},
				},
			}},
			itemType: "t1",
			attrs:    []string{"f1", "f4"},
			want:     false,
		},
		{
			name: "test false when input is superset",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Attrs: []definition.Attr{
						{Field: "f1"},
						{Field: "f2"},
						{Field: "f3"},
					},
				},
			}},
			itemType: "t1",
			attrs:    []string{"f1", "f2", "f3", "f4"},
			want:     false,
		},
		{
			name: "test true",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "tt1", IdentifierField: "id1", Alias: "a1",
					Attrs: []definition.Attr{
						{Field: "f1"},
						{Field: "f2"},
						{Field: "f3"},
					},
				},
			}},
			itemType: "tt1",
			attrs:    []string{"f1", "f3"},
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{def: tt.def}
			if got := ar.AreAttrsSupported(tt.itemType, tt.attrs); got != tt.want {
				t.Errorf("AwsReader.AreAttrsSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsReader_IsFilterSupported(t *testing.T) {
	tests := []struct {
		name     string
		def      *definition.Definition
		itemType string
		filter   reader.Filter
		want     bool
	}{
		{
			name: "test false when type does not exist",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Attrs: []definition.Attr{
						{Field: "a1"},
						{Field: "a2"},
						{Field: "a3"},
					},
				},
			}},
			itemType: "t111",
			filter:   reader.Filter{Attr: "a1"},
			want:     false,
		},
		{
			name: "test false",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Attrs: []definition.Attr{
						{Field: "a1"},
						{Field: "a2"},
						{Field: "a3"},
					},
				},
			}},
			itemType: "t1",
			filter:   reader.Filter{Attr: "a4"},
			want:     false,
		},
		{
			name: "test false 2",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
			}},
			itemType: "t1",
			filter:   reader.Filter{Attr: "a1"},
			want:     false,
		},
		{
			name: "test true",
			def: &definition.Definition{TypeDefinitions: []definition.TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Attrs: []definition.Attr{
						{Field: "a1"},
						{Field: "a2"},
						{Field: "a3"},
					},
				},
			}},
			itemType: "t1",
			filter:   reader.Filter{Attr: "a3"},
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{def: tt.def}
			if got := ar.IsFilterSupported(tt.itemType, tt.filter); got != tt.want {
				t.Errorf("AwsReader.IsFilterSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}
