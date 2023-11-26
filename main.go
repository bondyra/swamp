package main

import (
	"os"

	"github.com/bondyra/wtf/internal/command"
)

func main() {
	command.Execute(os.Args[1:])
}
