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
			name:     "test profiles list",
			input:    "in prof1 , prof2",
			expected: AST{Profiles: Collection{All: false, Elements: []string{"prof1", "prof2"}}},
		},
		{
			name:     "test profiles all",
			input:    "in *",
			expected: AST{Profiles: Collection{All: true}},
		},
		{
			name:  "test alias only",
			input: "item ns.res alias",
			expected: AST{
				Profiles: Collection{},
				Entitities: []Entity{
					ItemEntity{
						Value: Item{
							Type:  []string{"ns", "res"},
							Alias: "alias",
						},
					},
				},
			},
		},
		{
			name:  "test alias attrs",
			input: "item ns.res alias attr a1,a2 sub a2 add a3 where {some-go-code#123\"} add a4 sub a4",
			expected: AST{
				Entitities: []Entity{
					ItemEntity{
						Value: Item{
							Type:  []string{"ns", "res"},
							Alias: "alias",
							Modifiers: []Modifier{
								SetModifier{Value: Collection{Elements: []string{"a1", "a2"}}},
								SubModifier{Value: Collection{Elements: []string{"a2"}}},
								AddModifier{Value: Collection{Elements: []string{"a3"}}},
								SearchModifier{Value: "some-go-code#123\""},
								AddModifier{Value: Collection{Elements: []string{"a4"}}},
								SubModifier{Value: Collection{Elements: []string{"a4"}}},
							},
						},
					},
				},
			},
		},

		// {
		// 	name:  "test profiles and ",
		// 	input: "in prof1 , prof2 item ns.res alias attr a1,a2 sub a2 add a3 where {some-go-code#123\"}",
		// 	expected: AST{
		// 		Profiles: Collection{All: false, Elements: []string{"prof1", "prof2"}},
		// 		Entitities: []Entity{
		// 			Entity{
		// 				Item: Item{
		// 					Type:      []string{"ns", "res"},
		// 					Alias:     "alias1",
		// 					Modifiers: []Modifier{
		// 						Modifier{Set: []string ["a1", "a2"]}
		// 						Modifier{}
		// 						Modifier{}
		// 					},
		// 				},
		// 				Link: Link{},
		// 			},
		// 		},
		// 	},
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("unexpected error when parsing %v: %v", test.input, err)
			}
			if !cmp.Equal(test.expected, *actual) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expected, *actual)
			}
		})
	}
}
