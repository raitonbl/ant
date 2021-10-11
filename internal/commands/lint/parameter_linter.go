package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/project/structure"
)

type ParameterLinter struct {
}

func (instance *ParameterLinter) CanLint(ctx internal.ProjectContext, when Moment) bool {
	return ctx != nil && when == Document
}

func (instance *ParameterLinter) Lint(ctx internal.ProjectContext, document *structure.Specification, when Moment) ([]Violation, error) {
	document, err := ctx.GetDocument()

	if err != nil {
		return nil, err
	}

	if document == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedState()
	}

	problems := make([]Violation, 0)

	if document.Parameters == nil {
		return problems, nil
	}

	for index, parameter := range document.Parameters {
		context := &LintingContext{prefix: fmt.Sprintf("/parameters/%d", index), schema: parameter.Schema, when: when, document: document, isLocal: false}

		v, prob := lintParameter(context, parameter, when)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)
	}

	return problems, nil
}
