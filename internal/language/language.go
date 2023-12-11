package language

import (
	"regexp"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/bondyra/swamp/internal/common"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type AST struct {
	Profiles Collection `parser:"( \"in\" @@)?"`
	Query    Query      `parser:"@@"`
}

type Collection struct {
	All      bool     `parser:"  @\"*\""`
	Elements []string `parser:"| @Ident ( \",\" @Ident )*"`
}

type Query struct {
	Root     Item         `parser:"@@*"`
	Sequence []LinkedItem `parser:"@@*"`
}

type Item struct {
	Type      []string   `parser:"@Ident ( \".\" @Ident )?"`
	Modifiers []Modifier `parser:"@@*"`
}

type LinkedItem struct {
	Link Link `parser:"@@"`
	Item Item `parser:"@@"`
}

type Modifier interface{ value() }

type AttrModifier struct {
	Value Collection `parser:"\":\" @@"`
}

func (f AttrModifier) value() {}

type SearchModifier struct {
	Value SearchExpression `parser:"\"?\" @@"`
}

type SearchExpression struct {
	Attr  string          `parser:"@Ident"`
	Op    common.Operator `parser:"@Ident"`
	Value string          `parser:"@String"`
}

func (f SearchModifier) value() {}

type Link struct {
	FullOutput  bool `parser:"( @\"=\""`
	ShortOutput bool `parser:"| @\"-\" )"`
}

type Parser interface {
	ParseString(string) (*AST, error)
}

func ParseString(input string) (*AST, error) {
	return parseString(input)
}

func parseString(input string) (*AST, error) {
	var lexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: `Number`, Pattern: `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{Name: `String`, Pattern: `'[^']*'|"[^"]*"`},
		{Name: "Punct", Pattern: `[=\-[@#$%^&*()+_\|:;,.?/]|]`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	var parser = participle.MustBuild[AST](
		participle.Lexer(lexer),
		participle.Union[Modifier](AttrModifier{}, SearchModifier{}),
	)
	preprocess(input)
	return parser.ParseString("", input)
}

func preprocess(input string) string {
	output := input
	output = replace(output, `^\s+`, "")
	output = replace(output, `\s+$`, "")
	output = replace(output, `\s`, " ")
	output = replace(output, `  +`, " ")
	return output
}

func replace(input string, regex string, new string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(input, new)
}
