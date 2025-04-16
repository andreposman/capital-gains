package main

import (
	"github.com/andreposman/capital-gains/internal/infra/cli"
	"github.com/andreposman/capital-gains/pkg/helpers"
)

func main() {
	helpers.Greeting()
	cli.Handle()
}
