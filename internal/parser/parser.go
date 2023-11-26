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

type Entity struct {
	Item Item `  ( "item" @@)`
	Link Link `| ( "link" @@)`
}

type Item struct {
	Type      []string   `@Ident ( "." @Ident )*`
	Alias     string     `@Ident`
	Modifiers []Modifier `(@@)*`
}

type Link struct {
	From string `@Ident`
	To   string `@Ident`
}

type Modifier struct {
	Set       Collection `"attr" @@`
	Add       Collection `| "add" @@`
	Sub       Collection `| "sub" @@`
	SearchMod string     `| "where" @GoCode`
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
	)
	return parser.ParseString("", input)
}
