package golang

import (
	"embed"
	"fmt"
	"github.com/raitonbl/ant/internal"
	"github.com/raitonbl/ant/internal/project"
	"strings"
)

var (
	//go:embed app.mod.txt
	//go:embed app.main.txt
	//go:embed app.command.txt
	resources embed.FS
)

func GenerateProject(context internal.GenerateContext) (string, error) {
	var cfg Configuration
	cliObject, _ := context.GetDocument()
	err := context.BindConfiguration(&cfg)

	if err != nil {
		return "", internal.GetProblemFactory().GetProblem(err)
	}

	err = context.Write("go.mod", []byte(makeGoMode(cliObject, &cfg, true)))

	if err != nil {
		return "", err
	}

	err = context.Write("main.go", []byte(makeEntrypoint(cliObject, &cfg)))

	if err != nil {
		return "", err
	}

	return "", nil
}

func makeEntrypoint(cliObject *project.CliObject, cfg *Configuration) string {
	binary, _ := resources.ReadFile("app.main.txt")

	content := string(binary)

	content = strings.ReplaceAll(content, "${name}", *cliObject.Name)
	content = strings.ReplaceAll(content, "${version}", *cliObject.Version)
	content = strings.ReplaceAll(content, "${description}", *cliObject.Description)
	content = strings.ReplaceAll(content, "${module}", getModuleName(cliObject, cfg))

	additionalImports := ""
	additionalCommands := ""

	content = strings.ReplaceAll(content, "${cmd_section}", additionalCommands)
	content = strings.ReplaceAll(content, "${import_section}", additionalImports)

	return content
}

func makeCommand(context internal.GenerateContext, path []string, cfg *Configuration, cliObject *project.CliObject, commandObject *project.CommandObject) error {

	if commandObject.Subcommands != nil && len(commandObject.Subcommands) > 0 {
		for _, each := range commandObject.Subcommands {
			err := makeCommand(context, append(path, *commandObject.Name), cfg, cliObject, each)

			if err != nil {
				return err
			}
		}
	}

	binary, _ := resources.ReadFile("app.command.txt")

	content := string(binary)

	content = strings.ReplaceAll(content, "${description}", *commandObject.Description)
	content = strings.ReplaceAll(content, "${name}", strings.ToLower(*commandObject.Name))
	content = strings.ReplaceAll(content, "${name_in_capital}", strings.Title(*commandObject.Name))

	args := ""
	flags := ""

	for _, each := range commandObject.Parameters {
		parameter := toDefinitive(cliObject, &each)
		if parameter.In == nil || *parameter.In == project.Flags {
			flags += fmt.Sprintf("AddFlag(\"%s\", \"%s\", %s, \"%s\").", *parameter.Name, *parameter.Description, toCommandType(*parameter.Schema), *parameter.DefaultValue)
		} else {
			args += fmt.Sprintf("AddArgument(\"%s\", \"%s\", \"%s\").", *parameter.Name, *parameter.Description, *parameter.DefaultValue)
		}
	}

	content = strings.ReplaceAll(content, "${parameters}", fmt.Sprintf("%s%s", flags, args))

	err := context.WriteTo(path, *commandObject.Name+".go", []byte(content))

	if err != nil {
		return err
	}

	return nil
}

func toCommandType(_ project.Schema) string {
	return "command.String"
}

func toDefinitive(object *project.CliObject, each *project.ParameterObject) *project.ParameterObject {

	if each.RefersTo == nil {
		return each
	}

	return object.Components.Parameters[*each.RefersTo]
}

type Configuration struct {
	Sdk    Sdk    `properties:"sdk"`
	Module string `properties:"module"`
}

type Sdk struct {
	version string `properties:"version"`
}
