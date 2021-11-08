package cmd

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/raitonbl/ant/internal"
	"github.com/raitonbl/ant/internal/commands/generate"
	"github.com/thatisuday/commando"
	"os"
)

func AddGenerateProjectCommand(registry *commando.CommandRegistry) *commando.Command {
	return registry.Register("generate").
		SetShortDescription("generates a project from a specification file").
		SetDescription("allows the generate a project from a CLI specification file").
		AddFlag("properties", "the uri for a property file containing properties that will be used in project generation", commando.String, "").
		AddArgument("type", "the type of project, it might be application for CLI project or tests for integration tests project", "application").
		AddArgument("language", "the type of project, it might be application for CLI project or tests for integration tests project", "golang").
		AddArgument("file", "the CLI specification file URI", "index.json").
		AddArgument("directory", "the project directory which will contain the sources for the generated project", "./project").
		SetAction(doGenerateCLIProject)
}

func doGenerateCLIProject(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	uri := args["file"].Value
	propertiesURI := flags["properties"].Value
	projectDirectory := args["directory"].Value
	projectType := parseType(args["type"].Value)
	projectLanguage := parseLanguage(args["language"].Value)

	if projectType == nil {
		fmt.Println(fmt.Sprintf("Unknown project type: %s", args["type"].Value))
		os.Exit(1)
	}

	if projectLanguage == nil {
		fmt.Println(fmt.Sprintf("Unknown project language: %s", args["language"].Value))
		os.Exit(1)
	}

	var props *properties.Properties = nil

	if propertiesURI != "" {

		config, err := internal.GetProperties(propertiesURI.(string))

		if err != nil {
			exit(err.(*internal.Problem))
		}

		props = config
	}

	factory := internal.ContextFactory{}
	ctx, err := factory.SetFilename(uri).SetProjectLanguage(*projectLanguage).SetProjectType(*projectType).
		SetProjectDestination(projectDirectory).SetProperties(props).
		GetGenerateContext()

	if err != nil {
		exit(err.(*internal.Problem))
	}

	path, err := generate.Generate(ctx)

	if err != nil {
		exit(err.(*internal.Problem))
	}

	fmt.Println(fmt.Sprintf("Project has been generated in %s", path))
}

func parseLanguage(value string) *internal.LanguageType {
	var addr *internal.LanguageType = nil

	if value == string(internal.GoLang) {
		value := internal.GoLang
		addr = &value
	} else if value == string(internal.Python3) {
		value := internal.Python3
		addr = &value
	}

	return addr
}

func parseType(value string) *internal.ProjectType {
	var addr *internal.ProjectType = nil

	if value == string(internal.ApplicationType) {
		value := internal.ApplicationType
		addr = &value
	} else if value == string(internal.TestsType) {
		value := internal.TestsType
		addr = &value
	}

	return addr
}
