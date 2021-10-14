package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/lint_message"
	"github.com/raitonbl/cli/internal/project"
	"github.com/raitonbl/cli/internal/utils"
	"github.com/thoas/go-funk"
	"sort"
)

type CommandCacheContext struct {
	args       map[string]*project.Parameter
	flags      map[string]*project.Parameter
	shortForms map[string]*project.Parameter
}

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

		if cache[*parameter.Id] != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", ctx.prefix), Message: lint_message.DUPLICATED_FIELD_VALUE})
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

func doLintCommandParameterSection(commandContext *CommandLintingContext, document *project.Specification, instance *project.Command) ([]Violation, error) {

	prefix := commandContext.path
	problems := make([]Violation, 0)

	if instance.Parameters != nil {

		array, err := doLintCommandParameters(commandContext, instance, prefix, document)

		if err != nil {
			return nil, err
		}

		problems = append(problems, array...)
	}

	return problems, nil
}

func doLintCommandParameters(commandContext *CommandLintingContext, instance *project.Command, prefix string, document *project.Specification) ([]Violation, error) {

	problems := make([]Violation, 0)
	args := make(map[string]*project.Parameter)
	flags := make(map[string]*project.Parameter)
	shortForms := make(map[string]*project.Parameter)
	cacheContext := &CommandCacheContext{args: args, flags: flags, shortForms: shortForms}

	for index, each := range instance.Parameters {
		ctx := &LintContext{prefix: fmt.Sprintf("%s/parameters/%d", prefix, index), document: document, schemas: commandContext.schemaCache}
		array, err := doLintCommandParameter(commandContext, ctx, cacheContext, &each)

		if err != nil {
			return nil, err
		}


		problems = append(problems, array...)

	}

	ctx := &LintContext{prefix: fmt.Sprintf("%s/parameters", prefix), document: document, schemas: commandContext.schemaCache}
	problems = append(problems, doLintCommandParameterInArguments(ctx, args)...)

	return problems, nil
}

func doLintCommandParameter(commandContext *CommandLintingContext, ctx *LintContext, cacheContext *CommandCacheContext, each *project.Parameter) ([]Violation, error) {
	args := cacheContext.args
	flags := cacheContext.flags
	shortForms := cacheContext.shortForms

	param := each
	problems := make([]Violation, 0)
	array, skipLintParameter, param := doLintCommandParameterRefersTo(commandContext, ctx, *each)

	problems = append(problems, array...)

	if param == nil {
		return problems, nil
	}

	if param.In != nil && *param.In == project.Arguments && param.Name != nil && args[*param.Name] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s", ctx.prefix), Message: lint_message.NOT_AVAILABLE_IN_USE})
	}

	if param.In != nil && *param.In == project.Arguments {
		args[*param.Id] = param
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

func doLintCommandParameterRefersTo(commandContext *CommandLintingContext, ctx *LintContext, each project.Parameter) ([]Violation, bool, *project.Parameter) {
	param := &each
	skipLintParameter := false
	problems := make([]Violation, 0)
	isReference := isParameterReference(param)
	parameterCache := commandContext.parameterCache

	if each.RefersTo != nil && !isReference {
		param = nil
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	} else if each.RefersTo != nil && isReference {
		param = parameterCache[*each.RefersTo]

		if param == nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
			param = nil
		} else if (param.In == nil || *param.In == project.Flags) && each.Index != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
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

func doLintCommandParameterInFlags(ctx *LintContext, param *project.Parameter, flags map[string]*project.Parameter, shortForms map[string]*project.Parameter) []Violation {
	problems := make([]Violation, 0)

	if param.Name != nil && flags[*param.Name] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(name_format_pattern, ctx.prefix), Message: lint_message.NOT_AVAILABLE_IN_USE})
	} else if param.Name != nil {
		flags[*param.Name] = param
	}

	if param.ShortForm != nil && shortForms[*param.ShortForm] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/short-form", ctx.prefix), Message: lint_message.DUPLICATED_FIELD_VALUE})
	} else if param.ShortForm != nil {
		shortForms[*param.ShortForm] = param
	}

	return problems
}

func doLintCommandParameterInArguments(ctx *LintContext, args map[string]*project.Parameter) []Violation {
	problems := make([]Violation, 0)
	seq := make([]*project.Parameter, 0)

	for _, value := range args {
		if value.Index != nil {
			seq = append(seq, value)
		}
	}

	sort.Sort(ArgParameter(seq))

	indexes := make([]int, 0)

	for _, each := range seq {

		if len(indexes) == 0 && *each.Index != 0 {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s", ctx.prefix), Message: lint_message.ARGS_INDEX_NOT_ORDERED})
		} else if funk.Contains(indexes, *each.Index) {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s", ctx.prefix), Message: lint_message.ARGS_INDEX_NOT_UNIQUE})
		}

		indexes = append(indexes, *each.Index)
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

type ArgParameter []*project.Parameter

func (s ArgParameter) Len() int {
	return len(s)
}

func (s ArgParameter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ArgParameter) Less(i, j int) bool {
	return *s[i].Index < *s[j].Index
}
