package cmd

import (
	"fmt"
	"github.com/raitonbl/ant/pkg/resources"
	"github.com/thatisuday/commando"
	"os"
)

func AddExportCommand(registry *commando.CommandRegistry) *commando.Command {
	return registry.Register("export").
		SetShortDescription("retrieves the schema used during linting").
		SetDescription("retrieves the JSON schema used during linting").
		AddArgument("object", "object which export is intended\nschema - JSON schema for CLI definition", "schema").
		AddArgument("file", "file which will contain the exported object", "schema.json").
		SetAction(export)
}

func export(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	objectType := args["object"].Value

	switch objectType {
	case "schema":
		doExportSchema(args, flags)
	default:
		fmt.Println(fmt.Sprintf("Fail: Unknown object %s", objectType))
		os.Exit(1)
	}
}

func doExportSchema(args map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
	path := args["file"].Value

	binary, err := resources.GetResource("schema.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if path == "" {
		fmt.Println(string(binary))
		return
	}

	err = os.WriteFile(path, binary, 0644)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

}
