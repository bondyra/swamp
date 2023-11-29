package aws

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/reader"
	"golang.org/x/exp/slices"
)

func newMockProfileProvider(profiles []string, err error) profile.ProfileProvider {
	return func() ([]string, error) {
		return profiles, err
	}
}

type mockPoolFactory struct{}

func (maf mockPoolFactory) NewPool(profiles ...string) client.Pool {
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
			profileProvider := newMockProfileProvider(tt.profileProviderProvideProfilesOutput, tt.profileProviderProvideProfilesError)
			poolFactory := mockPoolFactory{}
			defFactory := mockDefFactory{mockDefFactoryOutput{tt.defFactoryFromFileOutput, tt.defFactoryFromFileError}}

			got, err := NewReader(profileProvider, poolFactory, defFactory)

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

const (
	ID_1_1 = iota
	ID_1_2
	ID_1_3
	ID_2_1
	ID_2_2
	ID_3_1
	ID_3_ERR
)

var (
	baseItemData []reader.ItemData = []reader.ItemData{
		{Identifier: fmt.Sprint(ID_1_1)},
		{Identifier: fmt.Sprint(ID_1_2)},
		{Identifier: fmt.Sprint(ID_1_3)},
		{Identifier: fmt.Sprint(ID_2_1)},
		{Identifier: fmt.Sprint(ID_2_2)},
		{Identifier: fmt.Sprint(ID_3_1)},
		{Identifier: fmt.Sprint(ID_3_ERR)},
	}
	itemData []reader.ItemData = []reader.ItemData{
		{Identifier: fmt.Sprint(ID_1_1), Properties: &map[string]string{"a": "a_1_1", "b": "11", "c": "true"}},
		{Identifier: fmt.Sprint(ID_1_2), Properties: &map[string]string{"a": "a_1_2", "b": "12", "c": "false"}},
		{Identifier: fmt.Sprint(ID_1_3), Properties: &map[string]string{"a": "a_1_3", "b": "13", "c": "true"}},
		{Identifier: fmt.Sprint(ID_2_1), Properties: &map[string]string{"a": "a_2_1", "b": "21", "c": "true"}},
		{Identifier: fmt.Sprint(ID_2_2), Properties: &map[string]string{"a": "a_2_2", "b": "22", "c": "false"}},
		{Identifier: fmt.Sprint(ID_3_1), Properties: &map[string]string{"a": "a_3_1", "b": "31", "c": "true"}},
		{Identifier: fmt.Sprint(ID_3_ERR), Properties: &map[string]string{"a": "a_3_2", "b": "32", "c": "false"}},
	}
)

const (
	PROFILE_1                     = "p_1"
	PROFILE_2                     = "p_2"
	TYPE_1                        = "type_1"
	TYPE_1_ID_FIELD               = "type_1_id_field"
	TYPE_2                        = "type_2"
	TYPE_2_ID_FIELD               = "type_2_id_field"
	TYPE_3                        = "type_3"
	TYPE_3_ID_FIELD               = "type_3_id_field"
	TYPE_THAT_CAUSES_ERR          = "type_err"
	TYPE_THAT_CAUSES_ERR_ID_FIELD = "type_err_id_field"
	TYPE_NOT_FOUND_IN_DEFINITION  = "type_def_err"
)

var testDefinition definition.Definition = definition.Definition{
	TypeDefinitions: []definition.TypeDefinition{
		{Type: TYPE_1, IdentifierField: TYPE_1_ID_FIELD, Alias: "Alias1", Attrs: []definition.Attr{{Field: "a"}, {Field: "b"}, {Field: "c"}}},
		{Type: TYPE_2, IdentifierField: TYPE_2_ID_FIELD, Alias: "Alias2", Attrs: []definition.Attr{{Field: "a"}, {Field: "b"}, {Field: "c"}}},
		{Type: TYPE_3, IdentifierField: TYPE_3_ID_FIELD, Alias: "Alias2", Attrs: []definition.Attr{{Field: "a"}, {Field: "b"}, {Field: "c"}}},
		{Type: TYPE_THAT_CAUSES_ERR, IdentifierField: TYPE_THAT_CAUSES_ERR_ID_FIELD, Alias: "Alias2"},
	},
}

type mockPool struct {
}

func (mp mockPool) ListResources(profile string, typeName string) ([]*reader.Item, error) {
	if typeName == TYPE_THAT_CAUSES_ERR {
		return nil, errors.New("some error")
	}
	switch typeName {
	case TYPE_1:
		return []*reader.Item{{Profile: profile, Data: &baseItemData[ID_1_1]}, {Profile: profile, Data: &baseItemData[ID_1_2]}, {Profile: profile, Data: &baseItemData[ID_1_3]}}, nil
	case TYPE_2:
		return []*reader.Item{{Profile: profile, Data: &baseItemData[ID_2_1]}, {Profile: profile, Data: &baseItemData[ID_2_2]}}, nil
	case TYPE_3:
		return []*reader.Item{{Profile: profile, Data: &baseItemData[ID_3_1]}, {Profile: profile, Data: &baseItemData[ID_3_ERR]}}, nil
	}
	panic("test setup error, unexpected type in mock " + typeName)
}

func (mp mockPool) GetResource(profile string, id string, typeName string) (*reader.Item, error) {
	if typeName == TYPE_THAT_CAUSES_ERR {
		return nil, errors.New("some error")
	}
	switch typeName {
	case TYPE_1, TYPE_2, TYPE_3:
		intId, err := strconv.Atoi(id)
		if err != nil {
			panic("test setup error, non number id in mock " + id)
		}
		switch intId {
		case ID_3_ERR:
			return nil, errors.New("some error")
		case ID_1_1, ID_1_2, ID_1_3, ID_2_1, ID_2_2, ID_3_1:
			return &reader.Item{Profile: profile, Data: &itemData[intId]}, nil
		default:
			panic("test setup error, unexpected id in mock " + id)
		}
	default:
		panic("test setup error, unexpected type in mock " + typeName)
	}
}

func TestAwsReader_GetItems(t *testing.T) {
	keyFunc := func(a *reader.Item) string { return fmt.Sprintf("%v+%v", a.Profile, a.Data.Identifier) }
	sortFunc := func(a *reader.Item, b *reader.Item) int { return strings.Compare(keyFunc(a), keyFunc(b)) }
	type args struct {
		itemType string
		profiles []string
		attrs    []string
		filters  []reader.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []*reader.Item
		wantErr bool
	}{
		{
			name:    "test no profiles returns nothing",
			args:    args{itemType: TYPE_1, profiles: []string{}, attrs: nil, filters: nil},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name: "test one profile no id filter",
			args: args{itemType: TYPE_1, profiles: []string{PROFILE_1}, attrs: nil, filters: nil},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_2]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_3]},
			},
			wantErr: false,
		},
		{
			name: "test multiple profiles no id filter",
			args: args{itemType: TYPE_1, profiles: []string{PROFILE_1, PROFILE_2}, attrs: nil, filters: nil},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_2]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_3]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_2]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_3]},
			},
			wantErr: false,
		},
		{
			name: "test multiple profiles with one id filter",
			args: args{
				itemType: TYPE_1,
				profiles: []string{PROFILE_1, PROFILE_2},
				attrs:    nil,
				filters:  []reader.Filter{{Attr: TYPE_1_ID_FIELD, Op: reader.OpEquals, Value: fmt.Sprint(ID_1_1)}},
			},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_1]},
			},
			wantErr: false,
		},
		{
			name: "test multiple profiles with multiple id filters",
			args: args{
				itemType: TYPE_1,
				profiles: []string{PROFILE_1, PROFILE_2},
				attrs:    nil,
				filters: []reader.Filter{
					{Attr: TYPE_1_ID_FIELD, Op: reader.OpLike, Value: fmt.Sprintf("[%v%v%v]", ID_1_1, ID_1_2, ID_1_3)},
					{Attr: TYPE_1_ID_FIELD, Op: reader.OpNotEquals, Value: fmt.Sprint(ID_1_1)},
					{Attr: TYPE_1_ID_FIELD, Op: reader.OpNotLike, Value: fmt.Sprint(ID_1_3)},
				},
			},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_2]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_2]},
			},
			wantErr: false,
		},
		{
			name: "test attr filters",
			args: args{
				itemType: TYPE_1,
				profiles: []string{PROFILE_1},
				attrs:    nil,
				filters: []reader.Filter{
					{Attr: "a", Op: reader.OpLike, Value: "a_1_*"},
					{Attr: "a", Op: reader.OpNotLike, Value: "a_1_2"},
					{Attr: "b", Op: reader.OpNotEquals, Value: "13"},
				},
			},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_1]},
			},
			wantErr: false,
		},
		{
			name:    "test error when type is not found in definition",
			args:    args{itemType: TYPE_NOT_FOUND_IN_DEFINITION, profiles: []string{PROFILE_1}, attrs: nil, filters: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test error on list resources error",
			args:    args{itemType: TYPE_THAT_CAUSES_ERR, profiles: []string{PROFILE_1}, attrs: nil, filters: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test error on get resource error",
			args:    args{itemType: TYPE_3, profiles: []string{PROFILE_1}, attrs: nil, filters: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test error on not supported filter op",
			args: args{
				itemType: TYPE_1,
				profiles: []string{PROFILE_1, PROFILE_2},
				attrs:    nil,
				filters:  []reader.Filter{{Attr: TYPE_1_ID_FIELD, Op: -1, Value: fmt.Sprint(ID_1_1)}},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test error on not supported filter",
			args: args{
				itemType: TYPE_1,
				profiles: []string{PROFILE_1, PROFILE_2},
				attrs:    nil,
				filters:  []reader.Filter{{Attr: "unknown", Op: reader.OpEquals, Value: fmt.Sprint(ID_1_1)}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{
				pool:         mockPool{},
				def:          &testDefinition,
				knownTypes:   nil,
				knownAliases: nil,
			}
			got, err := ar.GetItems(tt.args.itemType, tt.args.profiles, tt.args.attrs, tt.args.filters)
			slices.SortFunc(got, sortFunc)
			slices.SortFunc(tt.want, sortFunc)
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
