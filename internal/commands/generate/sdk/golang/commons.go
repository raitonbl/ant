package golang

import (
	"github.com/raitonbl/ant/internal/project"
	"strings"
)

func makeGoMode(cliObject *project.CliObject, cfg Configuration, isCli bool) []byte {
	name := cfg.Name
	version := cfg.Sdk.version

	if name == "" {
		name = *cliObject.Name
	}

	if version == "" {
		version = "1.16"
	}

	filename := "app.mod.txt"

	if !isCli {
		panic("Not implemented")
	}

	binary, _ := resources.ReadFile(filename)

	content := string(binary)

	content = strings.ReplaceAll(content, "${module}", name)
	content = strings.ReplaceAll(content, "${version}", version)

	return []byte(content)
}
