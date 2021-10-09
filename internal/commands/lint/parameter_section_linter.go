package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"github.com/raitonbl/cli/internal/project/structure"
	"github.com/raitonbl/cli/internal/utils"
)

const (
	parameterIndexFormatPattern = "/parameters/%d/index"
)

type ParameterSectionLinter struct {
}

func (instance *ParameterSectionLinter) CanLint(ctx internal.ProjectContext, when Moment) bool {
	return ctx != nil && when == Document
}

func (instance *ParameterSectionLinter) Lint(ctx internal.ProjectContext, document *structure.Specification, when Moment) ([]Violation, error) {
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
		v, prob := doLintParameter(index, parameter, when)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)
	}

	return problems, nil
}

func doLintParameter(index int, parameter structure.Parameter, when Moment) ([]Violation, error) {

	problems := make([]Violation, 0)

	problems = doLintParameterAttributes(index, parameter, when, problems)

	if parameter.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("/parameters/%d/refers-to", index), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	return doLintSchema(index, parameter, when, problems)
}

func doLintParameterAttributes(index int, parameter structure.Parameter, when Moment, problems []Violation) []Violation {
	if parameter.Id == nil || utils.IsBlank(*parameter.Id) {
		problems = append(problems, Violation{Path: fmt.Sprintf("/parameters/%d/id", index), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when})
	}

	if parameter.Description == nil || utils.IsBlank(*parameter.Description) {
		problems = append(problems, Violation{Path: fmt.Sprintf("/parameters/%d/description", index), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when})
	}

	if parameter.Name == nil || utils.IsBlank(*parameter.Name) {
		problems = append(problems, Violation{Path: fmt.Sprintf("/parameters/%d/name", index), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when})
	}

	if parameter.Index != nil && *parameter.Index < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(parameterIndexFormatPattern, index), Message: message.PARAMETER_INDEX_GT_ZERO_MESSAGE, Type: when})
	}

	if parameter.In == nil || *parameter.In == structure.Flags && parameter.Index != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(parameterIndexFormatPattern, index), Message: message.PARAMETER_FIELD_NOT_ALLOWED_IN_FLAGS, Type: when})
	}

	if parameter.In != nil && *parameter.In == structure.Arguments && parameter.Index == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(parameterIndexFormatPattern, index), Message: message.REQUIRED_PARAMETER_FIELD_WHEN_IN_ARGUMENTS, Type: when})
	}

	if parameter.In != nil && *parameter.In == structure.Arguments && parameter.ShortForm != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("/parameters/%d/short-form", index), Message: message.PARAMETER_FIELD_NOT_ALLOWED_IN_ARGUMENTS, Type: when})
	}
	return problems
}

func doLintSchema(index int, parameter structure.Parameter, when Moment, problems []Violation) ([]Violation, error) {
	if parameter.Schema == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("/parameters/%d/schema", index), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when})
	}

	if parameter.Schema != nil {
		ctx := LintingContext{prefix: fmt.Sprintf("/parameters/%d/schema", index), schema: parameter.Schema, when: when}
		problems = append(problems, ValidateSchema(&ctx)...)
	}

	return problems, nil
}
