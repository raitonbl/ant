package lint

import "github.com/raitonbl/cli/internal"

var Builder LinterBuilder

func Lint(context internal.ProjectContext) ([]Violation, error) {

	if context == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	if context.GetProjectFile() == nil || context.GetProjectFile().GetName() == "" {
		return nil, internal.GetProblemFactory().GetConfigurationFileNotFound()
	}

	builder := Builder

	if builder == nil {
		builder = getLinter()
	}

	object, err := builder.Build()

	if err != nil {
		return nil, err
	}

	problems, err := lint(context, object, Binary)

	if err != nil {
		return nil, err
	}

	if len(problems) > 0 {
		return problems, nil
	}

	return lint(context, object, Document)
}

func getLinter() LinterBuilder {
	builder := DelegatedLinterBuilder{}
	return builder.Append(&JsonSchemaLinter{})
}

func lint(context internal.ProjectContext, object Linter, when Moment) ([]Violation, error) {
	problems := make([]Violation, 0)

	if object.CanLint(context, when) {
		array, prob := object.Lint(context, nil, when)

		if prob != nil {
			return nil, prob
		}
		problems = append(problems, array...)
	}
	return problems, nil
}
