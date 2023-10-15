package handler

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/bondyra/wtf/internal/config"
)

type Handler interface {
	Execute(c config.Config)
}

type ConfigHandler struct {
	Args []string
}

type QueryHandler struct {
	Query string
}

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

func (q QueryHandler) Execute(c config.Config) {
	ast, _ := fmt.Println(parse(q.Query))
	fmt.Println(ast)
}

func parse(input string) (*AST, error) {
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
	return parser.ParseString("", input, participle.Trace(os.Stderr))
}
