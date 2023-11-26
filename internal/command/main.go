package command

import (
	"strings"

	"github.com/bondyra/wtf/internal/config"
	"github.com/bondyra/wtf/internal/handler"
)

func Execute(args []string) {
	switch args[0] {
	case "config":
		handler.ConfigHandler{Args: args[1:]}.Execute(config.Config{})
	default:
		handler.QueryHandler{Query: strings.Join(args[1:], " ")}.Execute(config.Config{})
	}
}
