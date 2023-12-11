package topology

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name            string
		inputFileExists bool
		inputContent    string
		want            *itemSchema
		wantErr         bool
	}{
		{
			name:            "test empty file",
			inputFileExists: true,
			inputContent:    "{}",
			want:            &itemSchema{},
			wantErr:         false,
		},
		{
			name:            "test empty file 2",
			inputFileExists: true,
			inputContent:    "{\"items\": []}",
			want:            &itemSchema{Items: []itemJson{}},
			wantErr:         false,
		},
		{
			name:            "test valid schema",
			inputFileExists: true,
			inputContent:    "{\"items\": [{\"type\": \"r1.t1\"}, {\"type\": \"r2.t2\", \"attrs\": [{\"field\": \"a1\"}]}], \"links\": []}",
			want:            &itemSchema{Items: []itemJson{{Type: NamespacedType{"r1", "t1"}}, {Type: NamespacedType{"r2", "t2"}, Attrs: []attrJson{{Field: "a1"}}}}},
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
			got, err := loadFromFile[itemSchema](inputPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFromFile() = %v, want %v", got, tt.want)
			}
			if tt.inputFileExists {
				os.Remove(inputPath)
			}
		})
	}
}
