package lint

import (
	"fmt"
	"github.com/raitonbl/ant/internal"
	"github.com/raitonbl/ant/internal/commands/lint/lint_message"
	"github.com/raitonbl/ant/internal/project"
	"github.com/thoas/go-funk"
	"sort"
)

type CommandCacheContext struct {
	args       map[string]*project.ParameterObject
	flags      map[string]*project.ParameterObject
	shortForms map[string]*project.ParameterObject
}

func doLintParameterSection(document *project.CliObject) ([]internal.Violation, error) {
	problems := make([]internal.Violation, 0)

	if document.Components == nil && document.Components.Parameters == nil {
		return problems, nil
	}

	for key, parameter := range document.Components.Parameters {

		ctx := &LintContext{prefix: fmt.Sprintf("/components/parameters/%s", key), document: document}

		array, prob := doLintParameter(ctx, parameter)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, array...)

	}

	return problems, nil
}

func doLintCommandParameterSection(commandContext *CommandLintingContext, document *project.CliObject, instance *project.CommandObject) ([]internal.Violation, error) {

	prefix := commandContext.path
	problems := make([]internal.Violation, 0)

	if instance.Parameters != nil {

		array, err := doLintCommandParameters(commandContext, instance, prefix, document)

		if err != nil {
			return nil, err
		}

		problems = append(problems, array...)
	}

	return problems, nil
}

func doLintCommandParameters(commandContext *CommandLintingContext, instance *project.CommandObject, prefix string, document *project.CliObject) ([]internal.Violation, error) {

	problems := make([]internal.Violation, 0)
	args := make(map[string]*project.ParameterObject)
	flags := make(map[string]*project.ParameterObject)
	shortForms := make(map[string]*project.ParameterObject)
	cacheContext := &CommandCacheContext{args: args, flags: flags, shortForms: shortForms}

	for index, each := range instance.Parameters {
		ctx := &LintContext{prefix: fmt.Sprintf("%s/parameters/%d", prefix, index), document: document}
		array, err := doLintCommandParameter(commandContext, ctx, cacheContext, &each)

		if err != nil {
			return nil, err
		}

		problems = append(problems, array...)

	}

	ctx := &LintContext{prefix: fmt.Sprintf("%s/parameters", prefix), document: document}
	problems = append(problems, doLintCommandParameterInArguments(ctx, args)...)

	return problems, nil
}

func doLintCommandParameter(commandContext *CommandLintingContext, ctx *LintContext, cacheContext *CommandCacheContext, each *project.ParameterObject) ([]internal.Violation, error) {
	args := cacheContext.args
	flags := cacheContext.flags
	shortForms := cacheContext.shortForms

	param := each
	problems := make([]internal.Violation, 0)
	array, skipLintParameter, param := doLintCommandParameterRefersTo(commandContext, ctx, *each)

	problems = append(problems, array...)

	if param == nil {
		return problems, nil
	}

	if param.In != nil && *param.In == project.Arguments && param.Name != nil && args[*param.Name] != nil {
		problems = append(problems, internal.Violation{Path: fmt.Sprintf("%s", ctx.prefix), Message: lint_message.NOT_AVAILABLE_IN_USE})
	}

	if param.In != nil && *param.In == project.Arguments && param.Name!=nil {
		args[*param.Name] = param //TODO FIX ID CHANGED INTO NAME
	}

	if param.In == nil || *param.In == project.Flags {
		problems = append(problems, doLintCommandParameterInFlags(ctx, param, flags, shortForms)...)
	}

	if skipLintParameter {
		return problems, nil
	}

	array, err := doLintParameter(ctx, param)

	if err != nil {
		return nil, err
	}

	return append(problems, array...), nil
}

func doLintCommandParameterRefersTo(commandContext *CommandLintingContext, ctx *LintContext, each project.ParameterObject) ([]internal.Violation, bool, *project.ParameterObject) {
	var param = &each
	skipLintParameter := false
	problems := make([]internal.Violation, 0)
	isReference := isParameterReference(param)

	if each.RefersTo != nil && !isReference {
		param = nil
		problems = append(problems, internal.Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	} else if each.RefersTo != nil && isReference {
		param = nil

		if ctx.document.Components != nil && ctx.document.Components.Parameters != nil {
			param = ctx.document.Components.Parameters[*each.RefersTo]
		}

		if param == nil {
			problems = append(problems, internal.Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
			param = nil
		} else if (param.In == nil || *param.In == project.Flags) && each.Index != nil {
			problems = append(problems, internal.Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
			param = nil
		} else if param.In != nil && *param.In == project.Arguments && each.Index != nil {
			param = param.Clone()
			param.RefersTo = nil
			param.Index = each.Index
			skipLintParameter = false
		} else if param.In != nil && *param.In == project.Arguments && each.Index == nil {
			skipLintParameter = true
		} else {
			param = nil
			skipLintParameter = true
		}
	}

	return problems, skipLintParameter, param
}

func doLintCommandParameterInFlags(ctx *LintContext, param *project.ParameterObject, flags map[string]*project.ParameterObject, shortForms map[string]*project.ParameterObject) []internal.Violation {
	problems := make([]internal.Violation, 0)

	if param.Name != nil && flags[*param.Name] != nil {
		problems = append(problems, internal.Violation{Path: fmt.Sprintf(name_format_pattern, ctx.prefix), Message: lint_message.NOT_AVAILABLE_IN_USE})
	} else if param.Name != nil {
		flags[*param.Name] = param
	}

	if param.ShortForm != nil && shortForms[*param.ShortForm] != nil {
		problems = append(problems, internal.Violation{Path: fmt.Sprintf("%s/short-form", ctx.prefix), Message: lint_message.DUPLICATED_FIELD_VALUE})
	} else if param.ShortForm != nil {
		shortForms[*param.ShortForm] = param
	}

	return problems
}

func doLintCommandParameterInArguments(ctx *LintContext, args map[string]*project.ParameterObject) []internal.Violation {
	problems := make([]internal.Violation, 0)
	seq := make([]*project.ParameterObject, 0)

	for _, value := range args {
		if value.Index != nil {
			seq = append(seq, value)
		}
	}

	sort.Sort(ArgParameter(seq))

	indexes := make([]int, 0)

	for _, each := range seq {

		if len(indexes) == 0 && *each.Index != 0 {
			problems = append(problems, internal.Violation{Path: fmt.Sprintf("%s", ctx.prefix), Message: lint_message.ARGS_INDEX_NOT_ORDERED})
		} else if funk.Contains(indexes, *each.Index) {
			problems = append(problems, internal.Violation{Path: fmt.Sprintf("%s", ctx.prefix), Message: lint_message.ARGS_INDEX_NOT_UNIQUE})
		}

		indexes = append(indexes, *each.Index)
	}

	return problems
}

func isParameterReference(each *project.ParameterObject) bool {

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

type ArgParameter []*project.ParameterObject

func (s ArgParameter) Len() int {
	return len(s)
}

func (s ArgParameter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ArgParameter) Less(i, j int) bool {
	return *s[i].Index < *s[j].Index
}
