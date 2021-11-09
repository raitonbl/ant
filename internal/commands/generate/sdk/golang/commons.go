package golang

import (
	"github.com/raitonbl/ant/internal/project"
	"strings"
)

func makeGoMode(cliObject *project.CliObject, cfg *Configuration, isCli bool) string {
	version := cfg.Sdk.version
	moduleName :=getModuleName(cliObject,cfg)

	if version == "" {
		version = "1.16"
	}

	filename := "app.mod.txt"

	if !isCli {
		panic("Not implemented")
	}

	binary, _ := resources.ReadFile(filename)

	content := string(binary)

	content = strings.ReplaceAll(content, "${version}", version)
	content = strings.ReplaceAll(content, "${module}", moduleName)

	return content
}

func getModuleName(cliObject *project.CliObject, cfg *Configuration) string {

	moduleName := cfg.Module

	if moduleName == "" {
		moduleName = *cliObject.Name
	}

	return moduleName
}
