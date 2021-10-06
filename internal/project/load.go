package project

import (
	"encoding/json"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/project/structure"
	"gopkg.in/yaml.v3"
	"strings"
)

func Load(context internal.ProjectContext) (*structure.Cli, error) {

	if context == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	if context.GetProjectFile() == nil || context.GetProjectFile().GetName() == "" {
		return nil, internal.GetProblemFactory().GetConfigurationFileNotFound()
	}

	filename := context.GetProjectFile().GetName()

	if !(strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml")) {
		return nil, internal.GetProblemFactory().GetUnsupportedDescriptor()
	}

	binary := context.GetProjectFile().GetContent()

	if strings.HasSuffix(filename, ".json") {
		return parseJson(binary)
	}

	return parseYaml(binary)
}

func parseYaml(binary []byte) (*structure.Cli, error) {

	if binary == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	descriptor := structure.Cli{}
	err := yaml.Unmarshal(binary, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, err
}

func parseJson(binary []byte) (*structure.Cli, error) {

	if binary == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	descriptor := structure.Cli{}
	err := json.Unmarshal(binary, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, err
}
