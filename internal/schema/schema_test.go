package schema

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestFromFile(t *testing.T) {
	tests := []struct {
		name            string
		inputFileExists bool
		inputContent    string
		want            *Schema
		wantErr         bool
	}{
		{
			name:            "test empty file",
			inputFileExists: true,
			inputContent:    "{}",
			want:            &Schema{},
			wantErr:         false,
		},
		{
			name:            "test empty file 2",
			inputFileExists: true,
			inputContent:    "{\"items\": [], \"links\": []}",
			want:            &Schema{Items: []ItemSchema{}, Links: []LinkSchema{}},
			wantErr:         false,
		},
		{
			name:            "test valid schema",
			inputFileExists: true,
			inputContent:    "{\"items\": [{\"type\": \"r1.t1\"}, {\"type\": \"r2.t2\", \"attrs\": [{\"field\": \"a1\"}]}], \"links\": []}",
			want:            &Schema{Items: []ItemSchema{{Type: NamespacedType{"r1", "t1"}}, {Type: NamespacedType{"r2", "t2"}, Attrs: []Attr{{Field: "a1"}}}}, Links: []LinkSchema{}},
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
		{
			name:            "test invalid schema",
			inputFileExists: true,
			inputContent:    "{\"items\": [{\"type\": \"r1.t1\"}, {\"type\": \"r1.t1\"}]}",
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
			got, err := fromFile(inputPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("fromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromFile() = %v, want %v", got, tt.want)
			}
			if tt.inputFileExists {
				os.Remove(inputPath)
			}
		})
	}
}

func TestFromFile_Full(t *testing.T) {
	_, err := fromFile("schema.json")

	if err != nil {
		t.Errorf("full schema error = %v", err)
		return
	}
}

func TestValidateItemSchema(t *testing.T) {
	tests := []struct {
		name    string
		is      []ItemSchema
		wantErr bool
	}{
		{
			name:    "test empty",
			is:      []ItemSchema{},
			wantErr: false,
		},
		{
			name: "test full schema",
			is: []ItemSchema{
				{Type: NamespacedType{"r1", "t1"}, Attrs: []Attr{{Field: "a1_1"}, {Field: "a1_2"}}},
				{
					Type:  NamespacedType{"r1", "t2"},
					Attrs: []Attr{{Field: "a2_1"}, {Field: "a2_2"}},
				},
			},
			wantErr: false,
		},
		{
			name: "test error when types are duplicated",
			is: []ItemSchema{
				{Type: NamespacedType{"r1", "t1"}},
				{Type: NamespacedType{"r1", "t1"}},
			},
			wantErr: true,
		},
		{
			name: "test error when any type is empty",
			is: []ItemSchema{
				{Type: NamespacedType{"r1", ""}},
				{Type: NamespacedType{"r2", "t2"}},
			},
			wantErr: true,
		},
		{
			name: "test error when any reader is empty",
			is: []ItemSchema{
				{Type: NamespacedType{"", "t1"}},
				{Type: NamespacedType{"r2", "t2"}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(Schema{Items: tt.is}); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLinkSchema(t *testing.T) {
	tests := []struct {
		name    string
		ls      []LinkSchema
		wantErr bool
	}{
		{
			name:    "test empty",
			ls:      []LinkSchema{},
			wantErr: false,
		},
		{
			name: "test full schema",
			ls: []LinkSchema{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r2", "t2"},
					Mappings: []Mapping{{From: "a1", To: "a2"}},
				},
				{
					From: NamespacedType{"r2", "t2"}, To: NamespacedType{"r1", "t1"},
					Mappings: []Mapping{{From: "a2", To: "a1"}},
				},
			},
			wantErr: false,
		},
		{
			name: "test error when any type is empty",
			ls: []LinkSchema{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r2", ""},
					Mappings: []Mapping{{From: "a1", To: "a2"}},
				},
			},
			wantErr: true,
		},
		{
			name: "test error when any reader is empty",
			ls: []LinkSchema{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"", "t2"},
					Mappings: []Mapping{{From: "a1", To: "a2"}},
				},
			},
			wantErr: true,
		},
		{
			name: "test error when any mapping attr is empty",
			ls: []LinkSchema{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r2", "t2"},
					Mappings: []Mapping{{From: "", To: "a2"}},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(Schema{Links: tt.ls}); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestIsTypeSupported(t *testing.T) {
	s := &Schema{
		Items: []ItemSchema{
			{Type: NamespacedType{"r1", "t1"}},
			{Type: NamespacedType{"r2", "t2"}},
		},
	}

	tests := []struct {
		name     string
		reader   string
		typ      string
		expected bool
	}{
		{
			name:     "test supported type",
			reader:   "r1",
			typ:      "t1",
			expected: true,
		},
		{
			name:     "test unsupported type",
			reader:   "r3",
			typ:      "t3",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.IsTypeSupported(tt.reader, tt.typ)
			if got != tt.expected {
				t.Errorf("IsTypeSupported() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestIsLinkSupported(t *testing.T) {
	s := &Schema{
		Links: []LinkSchema{
			{
				From: NamespacedType{"r1", "t1"},
				To:   NamespacedType{"r2", "t2"},
			},
		},
	}

	tests := []struct {
		name       string
		fromReader string
		fromType   string
		toReader   string
		toType     string
		expected   bool
	}{
		{
			name:       "test supported link",
			fromReader: "r1",
			fromType:   "t1",
			toReader:   "r2",
			toType:     "t2",
			expected:   true,
		},
		{
			name:       "test unsupported link",
			fromReader: "r2",
			fromType:   "t2",
			toReader:   "r1",
			toType:     "t1",
			expected:   false,
		},
		{
			name:       "test non-existent link",
			fromReader: "r1",
			fromType:   "t1",
			toReader:   "r3",
			toType:     "t3",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.IsLinkSupported(tt.fromReader, tt.fromType, tt.toReader, tt.toType)
			if got != tt.expected {
				t.Errorf("IsLinkSupported() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestAreAttrsSupported(t *testing.T) {
	s := Schema{
		Items: []ItemSchema{
			{Type: NamespacedType{"r1", "t1"}, Attrs: []Attr{{Field: "a1"}, {Field: "a2"}}},
			{Type: NamespacedType{"r2", "t2"}, Attrs: []Attr{{Field: "a3"}, {Field: "a4"}}},
		},
	}

	tests := []struct {
		name     string
		reader   string
		typ      string
		attrs    []string
		expected bool
	}{
		{
			name:     "test supported attrs",
			reader:   "r1",
			typ:      "t1",
			attrs:    []string{"a1", "a2"},
			expected: true,
		},
		{
			name:     "test unsupported attrs",
			reader:   "r1",
			typ:      "t1",
			attrs:    []string{"a1", "a2", "a3"},
			expected: false,
		},
		{
			name:     "test empty attrs",
			reader:   "r2",
			typ:      "t2",
			attrs:    []string{},
			expected: true,
		},
		{
			name:     "test empty schema",
			reader:   "r3",
			typ:      "t3",
			attrs:    []string{"a1", "a2"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.AreAttrsSupported(tt.reader, tt.typ, tt.attrs)
			if got != tt.expected {
				t.Errorf("AreAttrsSupported() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
