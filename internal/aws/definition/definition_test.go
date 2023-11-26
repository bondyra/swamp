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
			factory := DefaultFactory{inputPath}

			got, err := factory.NewDefinition()

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
				Parents: []ParentDefinition{{Type: "Type1", LinkType: "LinkType2_1", Links: []Link{{ParentField: "Link2_1ParentField1", Field: "Link2_1Field1"}}}},
				Attrs:   []Attr{{Field: "Attribute2_1"}, {Field: "Attribute2_2"}},
			},
		},
	}
	factory := DefaultFactory{testFilePath + "full_definition.json"}

	got, err := factory.NewDefinition()

	if err != nil {
		t.Errorf("DefaultFactory.NewDefinition() error = %v", err)
		return
	}
	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("DefaultFactory.NewDefinition() = %v, want %v", *got, expected)
	}
}

func TestDefinition_Validate(t *testing.T) {
	type fields struct {
		TypeDefinitions []TypeDefinition
	}
	tests := []struct {
		name     string
		testFile string
		wantErr  bool
	}{
		{
			name:     "test empty definition",
			testFile: "empty_definition.json",
			wantErr:  false,
		},
		{
			name:     "test full definition",
			testFile: "full_definition.json",
			wantErr:  false,
		},
		{
			name:     "test error when types are duplicated",
			testFile: "invalid_definition1.json",
			wantErr:  true,
		},
		{ // TODO: more test cases for this V
			name:     "test error when type is empty",
			testFile: "invalid_definition2.json",
			wantErr:  true,
		},
		{
			name:     "test error when link does not refer to defined type",
			testFile: "invalid_definition3.json",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, fileError := DefaultFactory{testFilePath + tt.testFile}.NewDefinition()

			if fileError != nil {
				t.Errorf("Definition.Validate() invalid test file: %v", tt.testFile)
			}

			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Definition.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
