package definition

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestDefaultReader_ReadDefinition(t *testing.T) {
	tests := []struct {
		name            string
		dr              Reader
		inputFileExists bool
		inputContent    string
		want            *Definition
		wantErr         bool
	}{
		{
			name:            "test empty file",
			dr:              &DefaultReader{},
			inputFileExists: true,
			inputContent:    "{}",
			want:            &Definition{},
			wantErr:         false,
		},
		{
			name:            "test non existent file",
			dr:              &DefaultReader{},
			inputFileExists: false,
			inputContent:    "{}",
			want:            nil,
			wantErr:         true,
		},
		{
			name:            "test invalid file",
			dr:              &DefaultReader{},
			inputFileExists: true,
			inputContent:    "{invalid json",
			want:            nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputPath := "path"
			if tt.inputFileExists {
				fileNamePrefix := strings.Replace(tt.name, " ", "-", -1)
				tempFile, _ := os.CreateTemp(".", fileNamePrefix)
				tempFile.Write([]byte(tt.inputContent))
				tempFile.Close()
				inputPath = tempFile.Name()
			}
			dr := &DefaultReader{}

			got, err := dr.ReadDefinition(inputPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultReader.ReadDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultReader.ReadDefinition() = %v, want %v", got, tt.want)
			}
			if tt.inputFileExists {
				os.Remove(inputPath)
			}
		})
	}
}

func TestDefaultReader_ReadDefinitionFull(t *testing.T) {
	dr := &DefaultReader{}
	expected := Definition{
		TypeDefinitions: []TypeDefinition{
			{Type: "Type1", IdentifierField: "Identifier1", Alias: "Alias1", Parents: []ParentDefinition{}, Attrs: []Attr{{Field: "Attribute1_1"}, {Field: "Attribute1_2"}}},
			{
				Type: "Type2", IdentifierField: "Identifier2", Alias: "Alias2",
				Parents: []ParentDefinition{{Type: "ParentType2_1", LinkType: "LinkType2_1", Links: []Link{{ParentField: "Link2_1ParentField1", Field: "Link2_1Field1"}}}},
				Attrs:   []Attr{{Field: "Attribute2_1"}, {Field: "Attribute2_2"}},
			},
		},
	}
	_, thisFileName, _, _ := runtime.Caller(0)

	fmt.Println(path.Join(path.Dir(thisFileName), "testing/test_definition.json"))
	content, err := dr.ReadDefinition(path.Join(path.Dir(thisFileName), "testing/test_definition.json"))

	if err != nil {
		t.Errorf("DefaultReader.ReadDefinition() error = %v", err)
		return
	}
	if !reflect.DeepEqual(*content, expected) {
		t.Errorf("DefaultReader.ReadDefinition() = %v, want %v", *content, expected)
	}
}
