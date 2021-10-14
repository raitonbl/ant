package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/lint_message"
	"github.com/raitonbl/cli/internal/project"
	"github.com/raitonbl/cli/internal/utils"
)

func doLintParameterSection(document *project.Specification, schemaCache map[string]*project.Schema) (map[string]*project.Parameter, []Violation, error) {
	problems := make([]Violation, 0)
	cache := make(map[string]*project.Parameter)

	if document.Parameters == nil {
		return cache, problems, nil
	}

	for index, parameter := range document.Parameters {

		ctx := &LintContext{prefix: fmt.Sprintf("/parameters/%d", index), document: document, schemas: schemaCache}

		if parameter.Id == nil || utils.IsBlank(*parameter.Id) {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
		}

		array, prob := doLintParameter(ctx, &parameter)

		if prob != nil {
			return nil, nil, prob
		}

		problems = append(problems, array...)

		param := parameter

		cache[*parameter.Id] = &param

	}

	return cache, problems, nil
}

func doLintCommandParameters(commandContext *CommandLintingContext, document *project.Specification, instance *project.Command) ([]Violation, error) {

	prefix := commandContext.path
	problems := make([]Violation, 0)

	if instance.Parameters != nil {

		array, err := doLintCommandParameter(commandContext, instance, prefix, document)

		if err != nil {
			return nil, err
		}

		problems = append(problems, array...)
	}

	return problems, nil
}

func doLintCommandParameter(commandContext *CommandLintingContext, instance *project.Command, prefix string, document *project.Specification) ([]Violation, error) {

	problems := make([]Violation, 0)
	argNames := make(map[string]*project.Parameter)
	flagNames := make(map[string]*project.Parameter)
	shortForms := make(map[string]*project.Parameter)

	for index, each := range instance.Parameters {
		param := &each
		ctx := &LintContext{prefix: fmt.Sprintf("%s/parameters/%d", prefix, index), document: document, schemas: commandContext.schemaCache}

		array, skipLintParameter, isUnresolvable := doLintCommandParameterRefersTo(commandContext, ctx, each)

		problems = append(problems, array...)

		if isUnresolvable {
			continue
		}

		if param.In == nil || *param.In == project.Flags {
			problems = append(problems, doLintCommandFlag(ctx, param, flagNames, shortForms)...)
		} else if param.In != nil && *param.In == project.Arguments && param.Name != nil && argNames[*param.Name] != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s", ctx.prefix), Message: lint_message.NOT_AVAILABLE_IN_USE})
		}

		if skipLintParameter {
			continue
		}

		array, err := doLintParameter(ctx, param)

		if err != nil {
			return nil, err
		}

		problems = append(problems, array...)

	}

	return problems, nil
}

func doLintCommandParameterRefersTo(commandContext *CommandLintingContext, ctx *LintContext, each project.Parameter) ([]Violation, bool, bool) {
	param := &each
	isUnresolvable := false
	skipLintParameter := false
	problems := make([]Violation, 0)
	isReference := isParameterReference(param)
	parameterCache := commandContext.parameterCache

	if each.RefersTo != nil && !isReference {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	} else if each.RefersTo != nil && isReference {
		param = parameterCache[*each.RefersTo]

		if param == nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
			isUnresolvable = true
		} else if (param.In == nil || *param.In == project.Flags) && each.Index != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
			isUnresolvable = true
		} else {
			skipLintParameter = true
		}
	}

	return problems, skipLintParameter, isUnresolvable
}

func doLintCommandFlag(ctx *LintContext, param *project.Parameter, flagNames map[string]*project.Parameter, shortForms map[string]*project.Parameter) []Violation {
	problems := make([]Violation, 0)

	if param.Name != nil && flagNames[*param.Name] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(name_format_pattern, ctx.prefix), Message: lint_message.NOT_AVAILABLE_IN_USE})
	} else if param.Name != nil {
		flagNames[*param.Name] = param
	}

	if param.ShortForm != nil && shortForms[*param.ShortForm] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/short-form", ctx.prefix), Message: lint_message.REPEATED_VALUE})
	} else if param.ShortForm != nil {
		shortForms[*param.ShortForm] = param
	}

	return problems
}

func isParameterReference(each *project.Parameter) bool {

	if each.Id != nil {
		return false
	}

	if each.Description != nil {
		return false
	}

	if each.Name != nil {
		return false
	}

	if each.In != nil {
		return false
	}

	if each.Required != nil {
		return false
	}

	if each.ShortForm != nil {
		return false
	}

	if each.DefaultValue != nil {
		return false
	}

	if each.Schema != nil {
		return false
	}

	if each.Index != nil && each.RefersTo == nil {
		return false
	}

	return true
}
