package lint

import (
	"fmt"
	"github.com/raitonbl/ant/internal/commands/lint/lint_message"
	"github.com/raitonbl/ant/internal/project"
	"github.com/raitonbl/ant/internal/utils"
)

const (
	schema_format_pattern     = "%s/format"
	minimum_format_pattern    = "%s/minimum"
	min_length_format_pattern = "%s/min-length"
	min_items_format_pattern  = "%s/min-items"
	index_format_pattern      = "%s/index"
	refers_to_format_pattern  = "%s/refers-to"
	name_format_pattern       = "%s/name"
)

type LintContext struct {
	prefix   string
	document *project.CliObject
	schemas  map[string]*project.Schema
}

func doLintSchema(ctx *LintContext, schema *project.Schema) []Violation {

	problems, skip := doLintSchemaRefersTo(ctx, schema)

	if skip {
		return problems
	}

	if schema.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	typeOf := *schema.TypeOf

	if typeOf != project.String && schema.Format != nil && (*schema.Format == project.Date || *schema.Format == project.DateTime || *schema.Format == project.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: lint_message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
	}

	if typeOf == project.String && schema.Format != nil && !(*schema.Format == project.Date || *schema.Format == project.DateTime || *schema.Format == project.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: lint_message.FIELD_FORMAT_NOT_ALLOWED_IN_TYPE_STRING})
	}

	problems = append(problems, doLintTextSchema(ctx, schema, typeOf)...)
	problems = append(problems, doLintArraySchema(ctx, schema, typeOf)...)
	problems = append(problems, doLintNumberSchema(ctx, schema, typeOf)...)

	if schema.Enum != nil && schema.Examples != nil && len(schema.Examples) > 0 {
		for i, example := range schema.Examples {
			if !belongsTo(schema.Enum, example) {
				problems = append(problems, Violation{Path: fmt.Sprintf("%s/examples/%d", ctx.prefix, i), Message: lint_message.FIELD_EXAMPLE_MUST_BE_PART_OF_ENUM})
			}
		}
	}

	return problems
}

func doLintSchemaRefersTo(ctx *LintContext, schema *project.Schema) ([]Violation, bool) {
	problems := make([]Violation, 0)

	if schema.TypeOf == nil && schema.RefersTo == nil {
		return []Violation{{Path: fmt.Sprintf("%s/type", ctx.prefix), Message: lint_message.REQUIRED_FIELD}}, true
	} else if schema.TypeOf == nil && schema.RefersTo != nil {

		fromCache := ctx.schemas[*schema.RefersTo]

		if fromCache == nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
		}

		return problems, true
	}
	return problems, false
}

func doLintTextSchema(ctx *LintContext, schema *project.Schema, typeOf project.SchemaType) []Violation {

	problems := make([]Violation, 0)

	if typeOf != project.String {
		if schema.MaxLength != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: lint_message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
		}

		if schema.MinLength != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: lint_message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
		}

		if schema.Pattern != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/pattern", ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
		}

	} else {
		problems = append(problems, doLintTextSchemaLength(ctx, schema)...)
	}

	return problems
}

func doLintTextSchemaLength(ctx *LintContext, schema *project.Schema) []Violation {
	problems := make([]Violation, 0)

	if schema.MinLength != nil && *schema.MinLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: lint_message.FIELD_MIN_LENGTH_GT_ZERO})
	}

	if schema.MaxLength != nil && *schema.MaxLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: lint_message.FIELD_MAX_LENGTH_GT_ZERO})
	}

	if schema.MaxLength != nil && schema.MinLength != nil {
		maximum := *schema.MaxLength
		minimum := *schema.MinLength

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: lint_message.FIELD_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH})
		}
	}
	return problems
}

func doLintNumberSchema(ctx *LintContext, schema *project.Schema, typeOf project.SchemaType) []Violation {

	problems := make([]Violation, 0)

	if schema.MultipleOf != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/multiple-of", ctx.prefix), Message: lint_message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
	}

	if typeOf == project.Number {
		if schema.Maximum != nil && schema.Minimum != nil {
			maximum := *schema.Maximum
			minimum := *schema.Minimum

			if minimum > maximum {
				problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: lint_message.FIELD_MIN_MUST_NOT_BE_GT_MAX})
			}

		}

		if schema.Maximum == nil && schema.ExclusiveMaximum != nil && *schema.ExclusiveMaximum {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
		}

		if schema.Minimum == nil && schema.ExclusiveMinimum != nil && *schema.ExclusiveMinimum {
			problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: lint_message.REQUIRED_FIELD})
		}
	} else {
		problems = append(problems, doLintNumberSchemaBoundary(ctx, schema, typeOf)...)
	}

	return problems
}

func doLintNumberSchemaBoundary(ctx *LintContext, schema *project.Schema, typeOf project.SchemaType) []Violation {
	problems := make([]Violation, 0)

	if typeOf != project.Number {

		if schema.Maximum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: lint_message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
		}

		if schema.Minimum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: lint_message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
		}

		if schema.ExclusiveMaximum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-maximum", ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
		}

		if schema.ExclusiveMinimum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-minimum", ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
		}

	}

	return problems
}

func doLintArraySchema(ctx *LintContext, schema *project.Schema, typeOf project.SchemaType) []Violation {

	problems := doLintArraySchemaLength(ctx, schema, typeOf)

	if typeOf != project.Array && schema.UniqueItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/unique-items", ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	if typeOf == project.Array && schema.Items != nil && schema.Items.TypeOf != nil && *schema.Items.TypeOf == project.Array {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items/type", ctx.prefix), Message: lint_message.ARRAY_FIELD_TYPE_NOT_ALLOWED})
		return problems
	}

	if typeOf == project.Array && schema.Items == nil && schema.RefersTo == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if schema.Items != nil {
		copyOf := &LintContext{prefix: ctx.prefix + "/items", schemas: ctx.schemas}
		problems = append(problems, doLintSchema(copyOf, schema.Items)...)
	}

	if typeOf == project.Array && schema.Format != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	return problems
}

func doLintArraySchemaLength(ctx *LintContext, schema *project.Schema, typeOf project.SchemaType) []Violation {

	problems := make([]Violation, 0)

	if typeOf != project.Array && schema.MaxItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	if typeOf != project.Array && schema.MinItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	if typeOf == project.Array && schema.MinItems != nil && *schema.MinItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: lint_message.FIELD_MIN_ITEMS_GT_ZERO})
	}

	if typeOf == project.Array && schema.MaxItems != nil && *schema.MaxItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: lint_message.FIELD_MAX_ITEMS_GT_ZERO})
	}

	if typeOf == project.Array && schema.MinItems != nil && schema.MaxItems != nil {
		maximum := *schema.MaxItems
		minimum := *schema.MinItems

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: lint_message.FIELD_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS})
		}

	}

	return problems
}

func doLintParameter(ctx *LintContext, parameter *project.ParameterObject) ([]Violation, error) {

	schema := parameter.Schema
	problems := make([]Violation, 0)
	problems = append(problems, doLintParameterFields(ctx, parameter)...)

	if parameter.Schema != nil && parameter.Schema.RefersTo != nil {
		schema = ctx.schemas[*parameter.Schema.RefersTo]

		if schema == nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/schema/refers-to", ctx.prefix), Message: lint_message.UNRESOLVABLE_FIELD})
		}

		return problems, nil
	}

	array, err := doLintParameterSchema(ctx, parameter)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	return problems, nil
}

func doLintParameterFields(ctx *LintContext, parameter *project.ParameterObject) []Violation {
	problems := make([]Violation, 0)

	if parameter == nil {
		return problems
	}

	if parameter.Description == nil || utils.IsBlank(*parameter.Description) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/description", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if parameter.Name == nil || utils.IsBlank(*parameter.Name) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/name", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if parameter.Index != nil && *parameter.Index < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(index_format_pattern, ctx.prefix), Message: lint_message.FIELD_INDEX_GT_ZERO})
	}

	if (parameter.In == nil || *parameter.In == project.Flags) && parameter.Index != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(index_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	if parameter.In != nil && *parameter.In == project.Arguments && parameter.Index == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(index_format_pattern, ctx.prefix), Message: lint_message.FIELD_WHEN_IN_ARGUMENTS})
	}

	if parameter.In != nil && *parameter.In == project.Arguments && parameter.ShortForm != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/short-form", ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	if parameter.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	return problems
}

func doLintParameterSchema(ctx *LintContext, parameter *project.ParameterObject) ([]Violation, error) {

	problems := make([]Violation, 0)

	if parameter.Schema == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/schema", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if parameter.Schema != nil {
		context := LintContext{prefix: fmt.Sprintf("%s/schema", ctx.prefix), schemas: ctx.schemas}
		problems = append(problems, doLintSchema(&context, parameter.Schema)...)
	}

	return problems, nil
}

func doLintExit(ctx *LintContext, exit *project.ExitObject) ([]Violation, error) {
	problems := make([]Violation, 0)

	if exit.Message == nil || utils.IsBlank(*exit.Message) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/message", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if exit.Code == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/code", ctx.prefix), Message: lint_message.REQUIRED_FIELD})
	}

	if exit.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: lint_message.FIELD_NOT_ALLOWED})
	}

	return problems, nil
}

func belongsTo(array []string, value string) bool {
	for _, each := range array {
		if each == value {
			return true
		}

	}
	return false
}
