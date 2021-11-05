package lint

import (
	"fmt"
	"github.com/raitonbl/ant/internal/commands/lint/lint_message"
	"github.com/raitonbl/ant/internal/project"
	"github.com/thoas/go-funk"
)

func doLintSchemaSection(document *project.CliObject) ([]Violation, error) {
	problems := make([]Violation, 0)

	if document.Components == nil && document.Components.Exits != nil {
		return problems, nil
	}

	keys := make([]string, 0)

	for key, schema := range document.Components.Schemas {
		if !funk.Contains(keys, key) {
			ctx := &LintContext{prefix: fmt.Sprintf("/components/schemas/%s", key), document: document}

			array, err := doLintSchemaFromSchemaSection(ctx, keys, schema)

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)
		}
	}

	return problems, nil
}

func doLintSchemaFromSchemaSection(ctx *LintContext, keys []string, schema *project.Schema) ([]Violation, error) {
	fromConfig := make(map[string]*project.Schema)

	if ctx.document.Components != nil && ctx.document.Components.Schemas != nil {
		fromConfig = ctx.document.Components.Schemas
	}

	problems := make([]Violation, 0)

	if schema.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	if schema.TypeOf == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/type", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if *schema.TypeOf == project.Array && schema.Items != nil && schema.Items.RefersTo != nil {

		if fromConfig[*schema.Items.RefersTo] == nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
		} else if funk.Contains(keys, *schema.Items.RefersTo) {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
		} else {
			keys = append(keys, *schema.Items.RefersTo)
			newCtx := &LintContext{prefix: fmt.Sprintf("/components/schemas/%s", *schema.Items.RefersTo), document: ctx.document}
			array, err := doLintSchemaFromSchemaSection(newCtx, keys, fromConfig[*schema.Items.RefersTo])

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)
		}

		return problems, nil
	}

	array := doLintSchema(ctx, schema)

	problems = append(problems, array...)

	return problems, nil
}
