package generate

import (
	"github.com/raitonbl/ant/internal"
	"github.com/raitonbl/ant/internal/commands/lint"
)

func Generate(context internal.GenerateContext) (string, error) {

	if context == nil {
		return "", internal.GetProblemFactory().GetUnexpectedContext()
	}

	if context.GetProjectFile() == nil || context.GetProjectFile().GetName() == "" {
		return "", internal.GetProblemFactory().GetConfigurationFileNotFound()
	}

	problems, err := lint.Lint(context)

	if err != nil {
		return "", err
	}

	if len(problems) > 0 {
		return "", internal.GetProblemFactory().GetValidationConstraintViolation(problems)
	}

	if context.GetTargetProjectType() == internal.ApplicationType {
		return doGenerateCLIProject(context)
	} else if context.GetTargetProjectType() == internal.TestsType {
		return doGenerateTestProject(context)
	} else {
		return "", internal.GetProblemFactory().NotImplemented()
	}

}
