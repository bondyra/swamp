package definition

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

const testFilePath string = "testing/"

func TestDefaultFactory_NewDefinition(t *testing.T) {
	tests := []struct {
		name            string
		inputFileExists bool
		inputContent    string
		want            *Definition
		wantErr         bool
	}{
		{
			name:            "test empty file",
			inputFileExists: true,
			inputContent:    "{}",
			want:            &Definition{},
			wantErr:         false,
		},
		{
			name:            "test non existent file",
			inputFileExists: false,
			inputContent:    "{}",
			want:            nil,
			wantErr:         true,
		},
		{
			name:            "test invalid file",
			inputFileExists: true,
			inputContent:    "{invalid json",
			want:            nil,
			wantErr:         true,
		},
	}
	factory := DefaultFactory{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputPath string
			if tt.inputFileExists {
				fileNamePrefix := strings.Replace(tt.name, " ", "-", -1)
				tempFile, _ := os.CreateTemp(".", fileNamePrefix)
				tempFile.Write([]byte(tt.inputContent))
				tempFile.Close()
				inputPath = tempFile.Name()
			}

			got, err := factory.FromFile(inputPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultFactory.NewDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultFactory.NewDefinition() = %v, want %v", got, tt.want)
			}
			if tt.inputFileExists {
				os.Remove(inputPath)
			}
		})
	}
}

func TestDefaultFactory_NewDefinitionFull(t *testing.T) {
	expected := Definition{
		TypeDefinitions: []TypeDefinition{
			{Type: "Type1", IdentifierField: "Identifier1", Alias: "Alias1", Parents: []ParentDefinition{}, Attrs: []Attr{{Field: "Attribute1_1"}, {Field: "Attribute1_2"}}},
			{
				Type: "Type2", IdentifierField: "Identifier2", Alias: "Alias2",
				Parents: []ParentDefinition{{Type: "Type1", LinkType: "inline", Links: []Link{{ParentField: "Link2_1ParentField1", Field: "Link2_1Field1"}}}},
				Attrs:   []Attr{{Field: "Attribute2_1"}, {Field: "Attribute2_2"}},
			},
		},
	}
	factory := DefaultFactory{}

	got, err := factory.FromFile(testFilePath + "full_definition.json")

	if err != nil {
		t.Errorf("DefaultFactory.NewDefinition() error = %v", err)
		return
	}
	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("DefaultFactory.NewDefinition() = %v, want %v", *got, expected)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		tds     []TypeDefinition
		wantErr bool
	}{
		{
			name: "test full definition",
			tds: []TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1", Parents: []ParentDefinition{}, Attrs: []Attr{{Field: "a1_1"}, {Field: "a1_2"}}},
				{
					Type: "t2", IdentifierField: "id2", Alias: "a2",
					Parents: []ParentDefinition{{Type: "t1", LinkType: "inline", Links: []Link{{ParentField: "l2_pf1", Field: "l2_f1"}}}},
					Attrs:   []Attr{{Field: "a2_1"}, {Field: "a2_2"}},
				},
			},
			wantErr: false,
		},
		{
			name: "test error when types are duplicated",
			tds: []TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
				{Type: "t1", IdentifierField: "id2", Alias: "a2"},
			},
			wantErr: true,
		},
		{
			name: "test error when aliases are duplicated",
			tds: []TypeDefinition{
				{Type: "t1", IdentifierField: "id1", Alias: "a1"},
				{Type: "t2", IdentifierField: "id2", Alias: "a1"},
			},
			wantErr: true,
		},
		{
			name: "test error when parent types are duplicated",
			tds: []TypeDefinition{
				{
					Type: "t1", IdentifierField: "id1", Alias: "a1",
					Parents: []ParentDefinition{
						{Type: "t1", LinkType: "inline", Links: []Link{{ParentField: "pf", Field: "f"}}},
						{Type: "t1", LinkType: "inline", Links: []Link{{ParentField: "pf", Field: "f"}}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test error when link type is invalid",
			tds: []TypeDefinition{
				{
					Type: "t2", IdentifierField: "id2", Alias: "a2",
					Parents: []ParentDefinition{
						{Type: "t1", LinkType: "inline", Links: []Link{{ParentField: "pf1", Field: "lf2"}}},
						{Type: "t1", LinkType: "invalidLinkType", Links: []Link{{ParentField: "pf1", Field: "lf2"}}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test error when any type is empty",
			tds: []TypeDefinition{
				{Type: "", IdentifierField: "Identifier1", Alias: "Alias1"},
				{Type: "Type2", IdentifierField: "Identifier2", Alias: "Alias2"},
			},
			wantErr: true,
		},
		{
			name: "test error when any alias is empty",
			tds: []TypeDefinition{
				{Type: "Type1", IdentifierField: "Identifier1", Alias: "Alias1"},
				{Type: "Type2", IdentifierField: "Identifier2", Alias: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(Definition{tt.tds}); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
