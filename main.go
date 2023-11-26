package main

import (
	"log"
	"os"
	"strings"

	"github.com/bondyra/swamp/internal/parser"
)

func main() {
	query := strings.Join(os.Args[1:], " ")
	ast, err := parser.ParseString(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", ast)
}
