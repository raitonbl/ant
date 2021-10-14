package lint

import (
	"fmt"
	"github.com/raitonbl/ant/internal/commands/lint/lint_message"
	"github.com/raitonbl/ant/internal/project"
	"github.com/raitonbl/ant/internal/utils"
	"github.com/thoas/go-funk"
)

func doLintSchemaSection(document *project.Specification) (map[string]*project.Schema, []Violation, error) {
	problems := make([]Violation, 0)
	cache := make(map[string]*project.Schema)

	if document.Schemas == nil {
		return cache, problems, nil
	}

	keys := make([]string, 0)
	fromConfig := make(map[string]*project.Schema)

	for _, schema := range document.Schemas {
		if schema.Id != nil {
			fromConfig[*schema.Id] = schema
		}
	}

	for index, schema := range document.Schemas {
		if !funk.Contains(keys, *schema.Id) {
			ctx := &LintContext{prefix: fmt.Sprintf("/schemas/%d", index), document: document, schemas: cache}

			array, err := doLintSchemaFromSchemaSection(ctx, fromConfig, keys, schema)

			if err != nil {
				return nil, nil, err
			}

			problems = append(problems, array...)
		}
	}

	return cache, problems, nil
}

func doLintSchemaFromSchemaSection(ctx *LintContext, fromConfig map[string]*project.Schema, keys []string, schema *project.Schema) ([]Violation, error) {

	cache := ctx.schemas
	problems := make([]Violation, 0)

	if schema.Id == nil || utils.IsBlank(*schema.Id) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

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
			newCtx := &LintContext{prefix: fmt.Sprintf("/schemas/%d", funk.IndexOf(ctx.document.Schemas, schema)), document: ctx.document, schemas: ctx.schemas}
			array, err := doLintSchemaFromSchemaSection(newCtx, fromConfig, keys, fromConfig[*schema.Items.RefersTo])

			if err != nil {
				return nil, err
			}

			problems = append(problems, array...)
		}

		cache[*schema.Id] = schema
		return problems, nil
	}

	array := doLintSchema(ctx, schema)

	problems = append(problems, array...)

	cache[*schema.Id] = schema

	return problems, nil
}
