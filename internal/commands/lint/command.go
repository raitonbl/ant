package lint

import (
	"fmt"
	"github.com/raitonbl/ant/internal/commands/lint/lint_message"
	"github.com/raitonbl/ant/internal/project"
	"github.com/raitonbl/ant/internal/utils"
)

func doLintCommandSection(document *project.Specification, parameterCache map[string]*project.Parameter, exitCache map[string]*project.Exit, schemas map[string]*project.Schema) ([]Violation, error) {
	problems := make([]Violation, 0)
	cache := make(map[string]*project.Command)

	if document.Subcommands == nil {
		return problems, nil
	}

	for index, command := range document.Subcommands {
		ctx := &CommandLintingContext{path: fmt.Sprintf("/commands/%d", index), parameterCache: parameterCache, exitCache: exitCache, schemaCache: schemas, commandCache: cache}

		v, prob := doLintCommand(ctx, &command, document)

		if prob != nil {
			return nil, prob
		}

		problems = append(problems, v...)
	}

	return problems, nil
}

func doLintCommand(commandContext *CommandLintingContext, instance *project.Command, document *project.Specification) ([]Violation, error) {

	cache := commandContext.commandCache
	prefix := commandContext.path
	problems := make([]Violation, 0)

	if instance.Id != nil && utils.IsBlank(*instance.Id) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", prefix), Message: lint_message.BLANK_FIELD})
	}

	if instance.Id != nil && cache[*instance.Id] != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", prefix), Message: lint_message.DUPLICATED_FIELD_VALUE})
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

func doLintCommandConfiguration(commandContext *CommandLintingContext, instance *project.Command, document *project.Specification) ([]Violation, error) {

	cache := commandContext.commandCache
	problems := make([]Violation, 0)

	array, err := doLintCommandExitSection(commandContext, document, instance)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	array, err = doLintCommandParameterSection(commandContext, document, instance)

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

func doLintSubcommands(commandContext *CommandLintingContext, document *project.Specification, instance *project.Command) ([]Violation, error) {
	cache := commandContext.commandCache
	prefix := commandContext.path
	exitCache := commandContext.exitCache
	schemaCache := commandContext.schemaCache
	parameterCache := commandContext.parameterCache

	problems := make([]Violation, 0)

	if instance.Subcommands != nil {
		for index, command := range instance.Subcommands {
			path := fmt.Sprintf("%s/commands/%d", prefix, index)
			ctx := &CommandLintingContext{path: path, parameterCache: parameterCache, exitCache: exitCache, schemaCache: schemaCache, commandCache: cache}
			array, err := doLintCommand(ctx, command, document)

			if err != nil {
				return nil, err
			}
			problems = append(problems, array...)
		}
	}

	return problems, nil
}
