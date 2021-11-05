package lint

import (
	"fmt"
	"github.com/raitonbl/ant/internal/commands/lint/lint_message"
	"github.com/raitonbl/ant/internal/project"
)

func doLintExitSection(document *project.CliObject) ([]Violation, error) {
	problems := make([]Violation, 0)

	if document.Components == nil && document.Components.Exits == nil {
		return problems, nil
	}

	for key, exit := range document.Components.Exits {

		ctx := &LintContext{prefix: fmt.Sprintf("/components/exits/%s", key), document: document}

		v, prob := doLintExit(ctx, exit)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)

	}

	return problems, nil
}

func doLintCommandExitSection(commandContext *CommandLintingContext, document *project.CliObject, instance *project.CommandObject) ([]Violation, error) {

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

func doLintCommandExit(commandContext *CommandLintingContext, document *project.CliObject, index int, each project.ExitObject) ([]Violation, error) {

	exit := &each
	prefix := commandContext.path
	problems := make([]Violation, 0)
	isReference := isExitReference(&each)
	ctx := &LintContext{prefix: fmt.Sprintf("%s/exit/%d", prefix, index), document: document}

	if each.RefersTo != nil && !isReference {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	} else if each.RefersTo != nil && isReference {

		if document.Components != nil && document.Components.Exits != nil  {
			exit = document.Components.Exits[*each.RefersTo]
		}

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

func isExitReference(each *project.ExitObject) bool {

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
