package main

import (
	"os"
	"strings"

	"github.com/bondyra/swamp/internal/cli"
)

func main() {
	query := strings.Join(os.Args[1:], " ")
	cli := cli.Cli{}
	cli.Run(query)
}
