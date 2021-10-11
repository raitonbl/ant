package lint

import (
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/project/structure"
)

var factory func() LinterBuilder

func Lint(context internal.ProjectContext) ([]Violation, error) {

	if context == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	if context.GetProjectFile() == nil || context.GetProjectFile().GetName() == "" {
		return nil, internal.GetProblemFactory().GetConfigurationFileNotFound()
	}

	builder := getFactory()()

	object, err := builder.Build()

	if err != nil {
		return nil, err
	}

	problems, err := lint(context, object, nil, Binary)

	if err != nil {
		return nil, err
	}

	if len(problems) > 0 {
		return problems, nil
	}

	document, err := context.GetDocument()

	if err != nil {
		return nil, err
	}

	return lint(context, object, document, Document)
}

func getFactory() func() LinterBuilder {
	if factory != nil {
		return factory
	} else {
		return func() LinterBuilder {
			instance := DelegatedLinterBuilder{}
			return instance.Append(&JsonSchemaLinter{}).Append(&ExitLinter{}).Append(&ParameterLinter{}).Append(&CommandLinter{})
		}
	}

}

func lint(context internal.ProjectContext, object Linter, document *structure.Specification, when Moment) ([]Violation, error) {
	problems := make([]Violation, 0)

	if object.CanLint(context, when) {
		array, prob := object.Lint(context, document, when)

		if prob != nil {
			return nil, prob
		}
		problems = append(problems, array...)
	}
	return problems, nil
}
