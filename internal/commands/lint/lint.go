package lint

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/qri-io/jsonschema"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/commands/lint/lint_message"
	"github.com/raitonbl/cli/internal/project/structure"
	"github.com/raitonbl/cli/internal/utils"
	"strings"
)

var (
	//go:embed schema.json
	resources embed.FS
)

type Violation struct {
	Path    string
	Message string
}

type CommandLintingContext struct {
	path           string
	cache          map[string]*structure.Command
	exitCache      map[string]*structure.Exit
	parameterCache map[string]*structure.Parameter
}

func Lint(context internal.ProjectContext) ([]Violation, error) {

	if context == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	if context.GetProjectFile() == nil || context.GetProjectFile().GetName() == "" {
		return nil, internal.GetProblemFactory().GetConfigurationFileNotFound()
	}

	problems, err := doLint(context)

	if err != nil {
		return nil, err
	}

	return problems, nil
}

func doLint(context internal.ProjectContext) ([]Violation, error) {
	problems := make([]Violation, 0)

	if strings.HasSuffix(context.GetProjectFile().GetName(), ".json") {
		array, err := doLintFile(context)

		if err != nil {
			return nil, err
		}

		if len(array) > 0 {
			return array, nil
		}

	}

	array, err := doLintObject(context)

	if err != nil {
		return nil, err
	}

	return append(problems, array...), nil
}

func doLintFile(ctx internal.ProjectContext) ([]Violation, error) {
	goContext := context.Background()

	binary, err := resources.ReadFile("schema.json")

	if err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	rs := &jsonschema.Schema{}

	if err = json.Unmarshal(binary, rs); err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	errs, err := rs.ValidateBytes(goContext, ctx.GetProjectFile().GetContent())

	if err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	problems := make([]Violation, len(errs))

	for index, each := range errs {
		problems[index] = Violation{Path: each.PropertyPath, Message: each.Message}
	}

	return problems, nil
}

func doLintObject(ctx internal.ProjectContext) ([]Violation, error) {

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

	parameterCache, array, err := doLintParameterSection(document)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	exitCache, array, err := doLintExitSection(document)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	array, err = doLintCommandSection(document, parameterCache, exitCache)

	if err != nil {
		return nil, err
	}

	return append(problems, array...), nil
}

func doLintExitSection(document *structure.Specification) (map[string]*structure.Exit, []Violation, error) {
	problems := make([]Violation, 0)
	cache := make(map[string]*structure.Exit)

	if document.Exit == nil {
		return cache, problems, nil
	}

	for index, exit := range document.Exit {

		ctx := &LintContext{prefix: fmt.Sprintf("/exit/%d", index), schema: nil, document: document}

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

func doLintParameterSection(document *structure.Specification) (map[string]*structure.Parameter, []Violation, error) {
	problems := make([]Violation, 0)
	cache := make(map[string]*structure.Parameter)

	if document.Parameters == nil {
		return cache, problems, nil
	}

	for index, parameter := range document.Parameters {

		ctx := &LintContext{prefix: fmt.Sprintf("/parameters/%d", index), schema: parameter.Schema, document: document}

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

func doLintCommandSection(document *structure.Specification, parameterCache map[string]*structure.Parameter, exitCache map[string]*structure.Exit) ([]Violation, error) {
	problems := make([]Violation, 0)
	cache := make(map[string]*structure.Command)

	if document.Subcommands == nil {
		return problems, nil
	}

	for index, command := range document.Subcommands {
		ctx := &CommandLintingContext{path: fmt.Sprintf("/commands/%d", index), parameterCache: parameterCache, exitCache: exitCache, cache: cache}

		v, prob := doLintCommand(ctx, command, document)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)
	}

	return problems, nil
}

func doLintCommand(commandContext *CommandLintingContext, instance *structure.Command, document *structure.Specification) ([]Violation, error) {

	cache := commandContext.cache
	prefix := commandContext.path
	problems := make([]Violation, 0)

	if instance.Id != nil && utils.IsBlank(*instance.Id) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", prefix), Message: lint_message.BLANK_FIELD})
	}

	if instance.Id != nil && cache[*instance.Id] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", prefix), Message: lint_message.REPEATED_VALUE})
	}

	if instance.Name == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(name_format_pattern, prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if instance.Name != nil && utils.IsBlank(*instance.Name) {
		problems = append(problems, Violation{Path: fmt.Sprintf(name_format_pattern, prefix), Message: lint_message.BLANK_FIELD})
	}

	if instance.Description == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/description", prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if instance.Description != nil && utils.IsBlank(*instance.Description) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/description", prefix), Message: lint_message.BLANK_FIELD})
	}

	if instance.Subcommands != nil && instance.Exit != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/exit", prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	if instance.Subcommands != nil && instance.Parameters != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/parameters", prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	array, err := doLintCommandConfiguration(commandContext, instance, document)

	if err != nil {
		return nil, err
	}

	return append(problems, array...), nil
}

func doLintCommandConfiguration(commandContext *CommandLintingContext, instance *structure.Command, document *structure.Specification) ([]Violation, error) {

	cache := commandContext.cache
	problems := make([]Violation, 0)

	array, err := doLintCommandExitSection(commandContext, document, instance)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	array, err = doLintCommandParameters(commandContext, document, instance)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	if instance.Id != nil {
		cache[*instance.Id] = instance
	}

	array, err = doLintSubcommands(commandContext, document, instance)

	if err != nil {
		return nil, err
	}

	return append(problems, array...), nil
}

func doLintCommandExitSection(commandContext *CommandLintingContext, document *structure.Specification, instance *structure.Command) ([]Violation, error) {

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

func doLintCommandExit(commandContext *CommandLintingContext, document *structure.Specification, index int, each structure.Exit) ([]Violation, error) {

	prefix := commandContext.path
	problems := make([]Violation, 0)
	exitCache := commandContext.exitCache

	exit := &each
	isReference := isExitReference(&each)
	ctx := &LintContext{prefix: fmt.Sprintf("%s/exit/%d", prefix, index), schema: nil, document: document}

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

func doLintCommandParameters(commandContext *CommandLintingContext, document *structure.Specification, instance *structure.Command) ([]Violation, error) {

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

func doLintCommandParameter(commandContext *CommandLintingContext, instance *structure.Command, prefix string, document *structure.Specification) ([]Violation, error) {

	problems := make([]Violation, 0)
	argNames := make(map[string]*structure.Parameter)
	flagNames := make(map[string]*structure.Parameter)
	shortForms := make(map[string]*structure.Parameter)

	for index, each := range instance.Parameters {
		param := &each
		ctx := &LintContext{prefix: fmt.Sprintf("%s/parameters/%d", prefix, index), schema: each.Schema, document: document}

		array, skipLintParameter, isUnresolvable := doLintCommandParameterRefersTo(commandContext, ctx, each)

		problems = append(problems, array...)

		if isUnresolvable {
			continue
		}

		if param.In == nil || *param.In == structure.Flags {
			problems = append(problems, doLintCommandFlag(ctx, param, flagNames, shortForms)...)
		} else if param.In != nil && *param.In == structure.Arguments && param.Name != nil && argNames[*param.Name] != nil {
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

func doLintCommandParameterRefersTo(commandContext *CommandLintingContext, ctx *LintContext, each structure.Parameter) ([]Violation, bool, bool) {
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
		} else if (param.In == nil || *param.In == structure.Flags) && each.Index != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
			isUnresolvable = true
		} else {
			skipLintParameter = true
		}
	}

	return problems, skipLintParameter, isUnresolvable
}

func doLintCommandFlag(ctx *LintContext, param *structure.Parameter, flagNames map[string]*structure.Parameter, shortForms map[string]*structure.Parameter) []Violation {
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

func doLintSubcommands(commandContext *CommandLintingContext, document *structure.Specification, instance *structure.Command) ([]Violation, error) {
	cache := commandContext.cache
	prefix := commandContext.path
	exitCache := commandContext.exitCache
	parameterCache := commandContext.parameterCache
	problems := make([]Violation, 0)

	if instance.Subcommands != nil {
		for index, command := range instance.Subcommands {
			path := fmt.Sprintf("%s/commands/%d", prefix, index)
			ctx := &CommandLintingContext{path: path, parameterCache: parameterCache, exitCache: exitCache, cache: cache}
			array, err := doLintCommand(ctx, command, document)

			if err != nil {
				return nil, err
			}
			problems = append(problems, array...)
		}
	}

	return problems, nil
}

func isExitReference(each *structure.Exit) bool {

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

func isParameterReference(each *structure.Parameter) bool {

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
