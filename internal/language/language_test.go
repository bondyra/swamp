package language

import (
	"reflect"
	"testing"

	"github.com/bondyra/swamp/internal/common"
)

type collectionOption func(*AST)
type rootOption func(*AST)
type itemOption func(*Item)
type seqOption func(*AST)
type linkedItemOption func(*LinkedItem)

func ast(collectionOpt collectionOption, rootOpt rootOption, seqOptions ...seqOption) AST {
	a := AST{}
	if collectionOpt != nil {
		collectionOpt(&a)
	}
	if rootOpt != nil {
		rootOpt(&a)
	}
	for _, opt := range seqOptions {
		opt(&a)
	}
	return a
}

func emptyProfiles() collectionOption {
	return func(a *AST) {
		a.Profiles = Collection{All: false, Elements: nil}
	}
}

func profiles(p ...string) collectionOption {
	return func(a *AST) {
		a.Profiles = Collection{All: false, Elements: p}
	}
}

func allProfiles() collectionOption {
	return func(a *AST) {
		a.Profiles = Collection{All: true}
	}
}

func root(opts ...itemOption) rootOption {
	return func(a *AST) {
		a.Query.Root = Item{}
		for _, opt := range opts {
			opt(&a.Query.Root)
		}
	}
}

func typ(t ...string) itemOption {
	return func(i *Item) {
		i.Type = t
	}
}

func attrs(attrs ...string) itemOption {
	return func(i *Item) {
		i.Modifiers = append(i.Modifiers, AttrModifier{Value: Collection{All: false, Elements: attrs}})
	}
}

func allAttrs() itemOption {
	return func(i *Item) {
		i.Modifiers = append(i.Modifiers, AttrModifier{Value: Collection{All: true}})
	}
}

func search(attr string, op common.Operator, value string) itemOption {
	return func(i *Item) {
		i.Modifiers = append(i.Modifiers, SearchModifier{Value: SearchExpression{Attr: attr, Op: op, Value: value}})
	}
}

func sequenceItem(opts ...linkedItemOption) seqOption {
	return func(a *AST) {
		it := LinkedItem{}
		for _, opt := range opts {
			opt(&it)
		}
		a.Query.Sequence = append(a.Query.Sequence, it)
	}
}

func fullOutput(opts ...itemOption) linkedItemOption {
	return func(it *LinkedItem) {
		it.Link.FullOutput = true
		for _, opt := range opts {
			opt(&it.Item)
		}
	}
}

func shortOutput(opts ...itemOption) linkedItemOption {
	return func(it *LinkedItem) {
		it.Link.ShortOutput = true
		for _, opt := range opts {
			opt(&it.Item)
		}
	}
}

func item(opts ...itemOption) linkedItemOption {
	return func(it *LinkedItem) {
		for _, opt := range opts {
			opt(&it.Item)
		}
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AST
	}{
		{
			name:     "test empty input",
			input:    "",
			expected: ast(nil, nil),
		},
		{
			name:     "test profiles list",
			input:    "in prof1 , prof2",
			expected: ast(profiles("prof1", "prof2"), nil),
		},
		{
			name:     "test profiles all",
			input:    "in *",
			expected: ast(allProfiles(), nil),
		},
		{
			name:     "test root single type",
			input:    "t1",
			expected: ast(emptyProfiles(), root(typ("t1"))),
		},
		{
			name:     "test root namespaced type",
			input:    "n1.t1",
			expected: ast(emptyProfiles(), root(typ("n1", "t1"))),
		},
		{
			name:     "test root attrs",
			input:    "n1.t1:a1,a2",
			expected: ast(emptyProfiles(), root(typ("n1", "t1"), attrs("a1", "a2"))),
		},
		{
			name:     "test root all attrs",
			input:    "n1.t1:*",
			expected: ast(emptyProfiles(), root(typ("n1", "t1"), allAttrs())),
		},
		{
			name:  "test root search",
			input: "n1.t1?a1 eq '1'",
			expected: ast(
				emptyProfiles(),
				root(
					typ("n1", "t1"), search("a1", common.EqualsTo, "'1'"),
				),
			),
		},
		{
			name:  "test root multiple modifiers",
			input: "n1.t1?a1 eq '1.1':a2,a3?a2 ne 'abc?:-':a4",
			expected: ast(
				profiles(),
				root(
					typ("n1", "t1"), search("a1", common.EqualsTo, "'1.1'"),
					attrs("a2", "a3"), search("a2", common.NotEqualsTo, "'abc?:-'"),
					attrs("a4"),
				),
			),
		},
		{
			name:  "test two items",
			input: "n1.t1-t2",
			expected: ast(
				emptyProfiles(),
				root(typ("n1", "t1")),
				sequenceItem(shortOutput(), item(typ("t2"))),
			),
		},
		{
			name:  "test full",
			input: "in prof1,prof2 n1.t1?a1 eq '2':a2,a3?a2 ne 'abc?:-':a4 - t2?a1 eq '1':a2,a3:a4=n2.t3",
			expected: ast(
				profiles("prof1", "prof2"),
				root(
					typ("n1", "t1"), search("a1", common.EqualsTo, "'2'"),
					attrs("a2", "a3"), search("a2", common.NotEqualsTo, "'abc?:-'"),
					attrs("a4"),
				),
				sequenceItem(
					shortOutput(),
					item(
						typ("t2"), search("a1", common.EqualsTo, "'1'"),
						attrs("a2", "a3"), attrs("a4"),
					),
				),
				sequenceItem(
					fullOutput(),
					item(typ("n2", "t3")),
				),
			),
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

func TestPreprocess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "test case 1",
			input:    "input string",
			expected: "input string",
		},
		{
			name:     "test case 2",
			input:    "	  	input 	 string      ",
			expected: "input string",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := preprocess(test.input)
			if actual != test.expected {
				t.Errorf("%s expected: %s, got: %s", test.name, test.expected, actual)
			}
		})
	}
}
