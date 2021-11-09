package golang

import (
	"fmt"
	"github.com/raitonbl/ant/internal"
	"github.com/raitonbl/ant/internal/project"
	"os"
	"path"
	"testing"
)

func TestEntrypoint(t *testing.T) {
	name := toPointer("ant-cli")
	version := toPointer("1.0.0")
	description := toPointer("<Description>")

	cfg := Configuration{Module: "github.com/raitonbl/ant"}
	cliObject := project.CliObject{Name: name, Version: version, Description: description}

	fmt.Println(makeEntrypoint(&cliObject, &cfg))
}

func TestMakeCommand(t *testing.T) {
	name := toPointer("ant-cli")
	version := toPointer("1.0.0")
	description := toPointer("<Description>")

	cliObject := project.CliObject{Name: name, Version: version, Description: description}
	cliObject.Subcommands = make([]project.CommandObject, 0)

	cliObject.Subcommands = append(cliObject.Subcommands, *buildCommand())

	factory := &internal.ContextFactory{}
	context, err := factory.SetFilename("testdata/index-001.json").SetProjectType(internal.ApplicationType).SetProjectLanguage(internal.GoLang).GetGenerateContext()

	if err != nil {
		panic(err)
	}

	err = makeCommand(context, []string{"cmd"}, nil, &cliObject, &cliObject.Subcommands[0])

	if err != nil {
		t.Fatal(err)
	}

	openFile(t, context, "cmd", "lint.go")
}

func buildCommand() *project.CommandObject {
	commandObject := project.CommandObject{}
	commandObject.Name = toPointer("lint")
	commandObject.Description = toPointer("<Description>")
	commandObject.Parameters = make([]project.ParameterObject, 0)

	commandObject.Parameters = append(commandObject.Parameters, *buildParameter("file", project.Arguments, project.String, "index.json"))
	commandObject.Parameters = append(commandObject.Parameters, *buildParameter("file", project.Flags, project.String, "index.json"))

	return &commandObject
}

func openFile(t *testing.T, ctx internal.GenerateContext, f ...string) {
	args := make([]string, 0)
	args = append(args, ctx.GetDirectory())
	args = append(args, f...)

	filename := path.Join(args...)

	binary, err := os.ReadFile(filename)

	if err != nil {
		t.Fatal(err)
	}

	content := fmt.Sprintf("%s\n------------\n%s", filename, string(binary))

	fmt.Println(content)
}

func buildParameter(name string, in project.In, schemaType project.SchemaType, defaultValue string) *project.ParameterObject {
	parameter := project.ParameterObject{}
	parameter.In = &in
	parameter.Name = toPointer(name)
	parameter.DefaultValue = toPointer(defaultValue)
	parameter.Description = toPointer("<Description>")
	parameter.Schema = &project.Schema{TypeOf: &schemaType}
	return &parameter
}

func toPointer(s string) *string {
	value := s
	return &value
}
