package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/lint_message"
	"github.com/raitonbl/cli/internal/project"
	"github.com/raitonbl/cli/internal/utils"
)

func doLintExitSection(document *project.Specification) (map[string]*project.Exit, []Violation, error) {
	problems := make([]Violation, 0)
	cache := make(map[string]*project.Exit)

	if document.Exit == nil {
		return cache, problems, nil
	}

	for index, exit := range document.Exit {

		ctx := &LintContext{prefix: fmt.Sprintf("/exit/%d", index), document: document}

		if exit.Id == nil || utils.IsBlank(*exit.Id) {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
		}

		v, prob := doLintExit(ctx, &exit)

		if prob != nil {
			return nil, nil, prob
		}

		problems = append(problems, v...)

		cache[*exit.Id] = &exit
	}

	return cache, problems, nil
}

func doLintCommandExitSection(commandContext *CommandLintingContext, document *project.Specification, instance *project.Command) ([]Violation, error) {

	problems := make([]Violation, 0)

	if instance.Exit != nil {
		for index, each := range instance.Exit {
			array, err := doLintCommandExit(commandContext, document, index, each)

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)

		}
	}
	return problems, nil
}

func doLintCommandExit(commandContext *CommandLintingContext, document *project.Specification, index int, each project.Exit) ([]Violation, error) {

	prefix := commandContext.path
	problems := make([]Violation, 0)
	exitCache := commandContext.exitCache

	exit := &each
	isReference := isExitReference(&each)
	ctx := &LintContext{prefix: fmt.Sprintf("%s/exit/%d", prefix, index), document: document}

	if each.RefersTo != nil && !isReference {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	} else if each.RefersTo != nil && isReference {
		exit = exitCache[*each.RefersTo]

		if exit == nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
		}

		return problems, nil
	}

	if exit != nil {
		array, err := doLintExit(ctx, exit)

		if err != nil {
			return nil, err
		}

		problems = append(problems, array...)
	}

	return problems, nil
}

func isExitReference(each *project.Exit) bool {

	if each.Id != nil {
		return false
	}

	if each.Code != nil {
		return false
	}

	if each.Message != nil {
		return false
	}

	if each.Description != nil {
		return false
	}

	return true
}

