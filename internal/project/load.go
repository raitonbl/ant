package project

import (
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/project/structure"
	"gopkg.in/yaml.v3"
	"os"
)

func Load(context internal.ProjectContext) (*structure.Cli, error) {

	if context == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	if context.GetDescriptor() == "" || context.GetDescriptor() == " " {
		return nil, internal.GetProblemFactory().GetConfigurationFileNotFound()
	}

	binary, err := os.ReadFile(context.GetDescriptor())

	if err != nil {
		return nil, err
	}

	descriptor := structure.Cli{}
	err = yaml.Unmarshal(binary, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, err
}
