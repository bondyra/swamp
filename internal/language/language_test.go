package language

import (
	"reflect"
	"testing"
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
			name:  "test attrs",
			input: "item ns.res alias attr a1,a2 sub a1,a2 add a3 add a4 sub a4",
			expected: AST{
				Entitities: []Entity{
					ItemEntity{
						Value: Item{
							Type:  []string{"ns", "res"},
							Alias: "alias",
							Modifiers: []Modifier{
								SetModifier{Value: Collection{Elements: []string{"a1", "a2"}}},
								SubModifier{Value: Collection{Elements: []string{"a1", "a2"}}},
								AddModifier{Value: Collection{Elements: []string{"a3"}}},
								AddModifier{Value: Collection{Elements: []string{"a4"}}},
								SubModifier{Value: Collection{Elements: []string{"a4"}}},
							},
						},
					},
				},
			},
		},
		{
			name:  "test filter",
			input: "item ns.res alias where a1 eq 'abc'",
			expected: AST{
				Entitities: []Entity{
					ItemEntity{
						Value: Item{
							Type:  []string{"ns", "res"},
							Alias: "alias",
							Modifiers: []Modifier{
								SearchModifier{Value: SearchExpression{Attr: "a1", Op: "eq", Value: SearchValue{String: "'abc'"}}},
							},
						},
					},
				},
			},
		},
		{
			name:  "test filter neg",
			input: "item ns.res alias where a1 ne 'abc'",
			expected: AST{
				Entitities: []Entity{
					ItemEntity{
						Value: Item{
							Type:  []string{"ns", "res"},
							Alias: "alias",
							Modifiers: []Modifier{
								SearchModifier{Value: SearchExpression{Attr: "a1", Op: "ne", Value: SearchValue{String: "'abc'"}}},
							},
						},
					},
				},
			},
		},
		{
			name:  "test filter integer",
			input: "item ns.res alias where a1 eq 1",
			expected: AST{
				Entitities: []Entity{
					ItemEntity{
						Value: Item{
							Type:  []string{"ns", "res"},
							Alias: "alias",
							Modifiers: []Modifier{
								SearchModifier{Value: SearchExpression{Attr: "a1", Op: "eq", Value: SearchValue{Number: 1}}},
							},
						},
					},
				},
			},
		},
		{
			name:  "test filter flag",
			input: "item ns.res alias where a1 eq true",
			expected: AST{
				Entitities: []Entity{
					ItemEntity{
						Value: Item{
							Type:  []string{"ns", "res"},
							Alias: "alias",
							Modifiers: []Modifier{
								SearchModifier{Value: SearchExpression{Attr: "a1", Op: "eq", Value: SearchValue{Boolean: true}}},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("unexpected error when parsing %v: %v", test.input, err)
			}
			if !reflect.DeepEqual(test.expected, *actual) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expected, *actual)
			}
		})
	}
}
