package main

import (
	"github.com/raitonbl/ant/cmd"
	"github.com/thatisuday/commando"
)

func main() {

	registry := commando.SetExecutableName("cli").
		SetVersion("1.0.0").
		SetDescription("manipulates cli specification language")

	cmd.AddLintCommand(registry)

	registry.Parse(nil)
}
