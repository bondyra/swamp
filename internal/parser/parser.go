package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type AST struct {
	Profiles   Collection `( "in" @@)?`
	Entitities []Entity   `@@*`
}

type Collection struct {
	All      bool     `  @"*"`
	Elements []string `| @Ident ( "," @Ident )*`
}

type Entity interface{ value() }

type ItemEntity struct {
	Value Item `"item" @@`
}

func (f ItemEntity) value() {}

type LinkEntity struct {
	Value Link `"item" @@`
}

func (f LinkEntity) value() {}

type Item struct {
	Type      []string   `@Ident ( "." @Ident )*`
	Alias     string     `@Ident`
	Modifiers []Modifier `@@*`
}

type Modifier interface{ value() }

type SetModifier struct {
	Value Collection `"attr" @@`
}

func (f SetModifier) value() {}

type AddModifier struct {
	Value Collection `"add" @@`
}

func (f AddModifier) value() {}

type SubModifier struct {
	Value Collection `"sub" @@`
}

func (f SubModifier) value() {}

type SearchModifier struct {
	Value string `"where" @GoCode`
}

func (f SearchModifier) value() {}

type Link struct {
	From string `@Ident`
	To   string `@Ident`
}

func ParseString(input string) (*AST, error) {
	return parseString(input)
}

func parseString(input string) (*AST, error) {
	var lexer = lexer.MustSimple([]lexer.SimpleRule{
		{`GoCode`, `{[^}]+}`},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_]*`},
		{"Punct", `[-[!@#$%^&*()+_=\|:;<,>.?/]|]`},
		{"whitespace", `\s+`},
	})
	var parser = participle.MustBuild[AST](
		participle.Lexer(lexer),
		participle.Unquote("GoCode"),
		participle.Union[Modifier](SetModifier{}, SubModifier{}, AddModifier{}, SearchModifier{}),
		participle.Union[Entity](ItemEntity{}, LinkEntity{}),
	)
	return parser.ParseString("", input)
}
