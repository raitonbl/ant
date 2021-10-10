package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/project/structure"
)

type ExitLinter struct {
}

func (instance *ExitLinter) CanLint(ctx internal.ProjectContext, when Moment) bool {
	return ctx != nil && when == Document
}

func (instance *ExitLinter) Lint(ctx internal.ProjectContext, document *structure.Specification, when Moment) ([]Violation, error) {
	document, err := ctx.GetDocument()

	if err != nil {
		return nil, err
	}

	if document == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedState()
	}

	problems := make([]Violation, 0)

	if document.Exit == nil {
		return problems, nil
	}

	for index, exit := range document.Exit {
		context := &LintingContext{prefix: fmt.Sprintf("/exit/%d", index), schema: nil, when: when, resolve: false}

		v, prob := lintExit(context, exit, when)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)
	}

	return problems, nil
}
