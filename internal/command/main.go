package command

import (
	"strings"

	"github.com/bondyra/wtf/internal/handler"
)

func Execute(args []string) {
	handler.QueryHandler{Query: strings.Join(args, " ")}.Execute()
}
