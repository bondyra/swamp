package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AST
	}{
		{
			name:     "test empty input",
			input:    "",
			expected: AST{},
		},
		{
			name:     "test profiles only",
			input:    "in prof1 , prof2",
			expected: AST{Profiles: Collection{All: false, Elements: []string{"prof1", "prof2"}}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("unexpected error when parsing %v: %v", test.input, err)
			}
			if !cmp.Equal(test.expected, *actual) {
				t.Errorf("%s expected %v, got %v", test.name, test.expected, *actual)
			}
		})
	}
}
