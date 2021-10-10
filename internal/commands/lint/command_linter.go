package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"github.com/raitonbl/cli/internal/project/structure"
	"github.com/raitonbl/cli/internal/utils"
)

type CommandLinter struct {
}

func (instance *CommandLinter) CanLint(ctx internal.ProjectContext, when Moment) bool {
	return ctx != nil && when == Document
}

func (instance *CommandLinter) Lint(ctx internal.ProjectContext, document *structure.Specification, when Moment) ([]Violation, error) {
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

	inCache := make(map[string]*structure.Command)

	for index, command := range document.Subcommands {
		v, prob := doLintCommand(index, command, when, "", inCache)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)
	}

	return problems, nil
}

func doLintCommand(index int, instance *structure.Command, when Moment, prefix string, cache map[string]*structure.Command) ([]Violation, error) {
	problems := make([]Violation, 0)

	if instance.Id != nil && utils.IsBlank(*instance.Id) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/id", prefix, index), Message: message.BLANK_FIELD_MESSAGE, Type: when})
	}

	if instance.Id != nil && cache[*instance.Id] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/id", prefix, index), Message: message.REPEATED_VALUE_MESSAGE, Type: when})
	}

	if instance.Name == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/name", prefix, index), Message: message.REQUIRED_FIELD_MESSAGE, Type: when})
	}

	if instance.Name != nil && utils.IsBlank(*instance.Name) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/name", prefix, index), Message: message.BLANK_FIELD_MESSAGE, Type: when})
	}

	if instance.Description == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/description", prefix, index), Message: message.REQUIRED_FIELD_MESSAGE, Type: when})
	}

	if instance.Description != nil && utils.IsBlank(*instance.Description) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/description", prefix, index), Message: message.BLANK_FIELD_MESSAGE, Type: when})
	}

	if instance.Subcommands != nil && instance.Exit != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/exit", prefix, index), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if instance.Subcommands != nil && instance.Parameters != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/commands/%d/parameters", prefix, index), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if instance.Exit != nil {
		for i, exit := range instance.Exit {
			context := &LintingContext{prefix: fmt.Sprintf("%s/commands/%d/exit/%d", prefix, index, i), schema: nil, when: when, resolve: true}
			array, err := lintExit(context, exit, when)

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)
		}
	}

	if instance.Parameters != nil {
		for i, param := range instance.Parameters {
			context := &LintingContext{prefix: fmt.Sprintf("%s/commands/%d/parameters/%d", prefix, index, i), schema: param.Schema, when: when, resolve: true}
			array, err := lintParameter(context, param, when)

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)
		}
	}

	if instance.Id != nil {
		cache[*instance.Id] = instance
	}

	return problems, nil
}
