package main

import (
	"embed"
	"fmt"
	"github.com/raitonbl/ant/cmd"
	"github.com/thatisuday/commando"
	"os"
)

var (
	//go:embed docs/version
	//go:embed docs/feat
	resources embed.FS
)

var version string

func main() {

	binary, err := resources.ReadFile("docs/version")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	version = string(binary)

	registry := commando.
		SetExecutableName("ant").
		SetVersion(version).
		SetDescription("manipulates cli specification language")

	registry.
		Register(nil).
		SetAction(index)

	cmd.AddLintCommand(registry)

	registry.Parse(nil)
}

func index(_ map[string]commando.ArgValue, _ map[string]commando.FlagValue) {

	binary, err := resources.ReadFile("docs/feat")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	feat := string(binary)

	fmt.Println(fmt.Sprintf("Ant: %s\nFeatures:\n %s", version, feat))
}
