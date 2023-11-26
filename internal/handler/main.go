package handler

import (
	"strings"

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

func (d ConfigHandler) Execute(c config.Config) {
	println("config -" + strings.Join(d.Args, "+"))
}

func (q QueryHandler) Execute(c config.Config) {
	println("query -" + q.Query)
}
