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
		v, prob := doLintCommand(command, when, document, fmt.Sprintf("/commands/%d", index), inCache)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)
	}

	return problems, nil
}

func doLintCommand( instance *structure.Command, when Moment, document *structure.Specification, prefix string, cache map[string]*structure.Command) ([]Violation, error) {
	problems := make([]Violation, 0)

	if instance.Id != nil && utils.IsBlank(*instance.Id) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", prefix), Message: message.BLANK_FIELD, Type: when})
	}

	if instance.Id != nil && cache[*instance.Id] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", prefix), Message: message.REPEATED_VALUE, Type: when})
	}

	if instance.Name == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/name", prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if instance.Name != nil && utils.IsBlank(*instance.Name) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/name", prefix), Message: message.BLANK_FIELD, Type: when})
	}

	if instance.Description == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/description", prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if instance.Description != nil && utils.IsBlank(*instance.Description) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/description", prefix), Message: message.BLANK_FIELD, Type: when})
	}

	if instance.Subcommands != nil && instance.Exit != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/exit", prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if instance.Subcommands != nil && instance.Parameters != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/parameters", prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if instance.Exit != nil {
		for index, each := range instance.Exit {
			context := &LintingContext{prefix: fmt.Sprintf("%s/exit/%d", prefix, index), schema: nil, when: when, document: document, isLocal: true}
			array, err := lintExit(context, each, when)

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)
		}
	}

	if instance.Parameters != nil {
		for index, each := range instance.Parameters {
			context := &LintingContext{prefix: fmt.Sprintf("%s/parameters/%d", prefix, index), schema: each.Schema, when: when, document: document, isLocal: true}
			array, err := lintParameter(context, each, when)

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)
		}
	}

	if instance.Id != nil {
		cache[*instance.Id] = instance
	}

	if instance.Subcommands != nil {
		for index, each := range instance.Subcommands {
			path := fmt.Sprintf("%s/commands/%d", prefix, index)
			array, err := doLintCommand( each, when, document, path, cache)

			if err != nil {
				return nil, err
			}
			problems = append(problems, array...)
		}
	}

	return problems, nil
}
