package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"github.com/raitonbl/cli/internal/project/structure"
	"github.com/raitonbl/cli/internal/utils"
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
	schema   *structure.Schema
	document *structure.Specification
}

func doLintSchema(ctx *LintContext) []Violation {

	schema := ctx.schema

	problems := make([]Violation, 0)

	if schema.TypeOf == nil {
		return []Violation{{Path: fmt.Sprintf("%s/type", ctx.prefix), Message: message.REQUIRED_FIELD}}
	}

	typeOf := *schema.TypeOf

	if typeOf != structure.String && schema.Format != nil && (*schema.Format == structure.Date || *schema.Format == structure.DateTime || *schema.Format == structure.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
	}

	if typeOf == structure.String && schema.Format != nil && !(*schema.Format == structure.Date || *schema.Format == structure.DateTime || *schema.Format == structure.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_NOT_ALLOWED_IN_TYPE_STRING})
	}

	problems = append(problems, doLintTextSchema(ctx, typeOf)...)
	problems = append(problems, doLintArraySchema(ctx, typeOf)...)
	problems = append(problems, doLintNumberSchema(ctx, typeOf)...)

	if schema.Enum != nil && schema.Examples != nil && len(schema.Examples) > 0 {
		for i, example := range schema.Examples {
			if !belongsTo(schema.Enum, example) {
				problems = append(problems, Violation{Path: fmt.Sprintf("%s/examples/%d", ctx.prefix, i), Message: message.FIELD_EXAMPLE_MUST_BE_PART_OF_ENUM})
			}
		}
	}

	return problems
}

func doLintTextSchema(ctx *LintContext, typeOf structure.SchemaType) []Violation {

	schema := ctx.schema
	problems := make([]Violation, 0)

	if typeOf != structure.String {
		if schema.MaxLength != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
		}

		if schema.MinLength != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
		}

		if schema.Pattern != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/pattern", ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
		}

	} else {
		problems = append(problems, doLintTextSchemaLength(ctx, schema)...)
	}

	return problems
}

func doLintTextSchemaLength(ctx *LintContext, schema *structure.Schema) []Violation {
	problems := make([]Violation, 0)

	if schema.MinLength != nil && *schema.MinLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.FIELD_MIN_LENGTH_GT_ZERO})
	}

	if schema.MaxLength != nil && *schema.MaxLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: message.FIELD_MAX_LENGTH_GT_ZERO})
	}

	if schema.MaxLength != nil && schema.MinLength != nil {
		maximum := *schema.MaxLength
		minimum := *schema.MinLength

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.FIELD_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH})
		}
	}
	return problems
}

func doLintNumberSchema(ctx *LintContext, typeOf structure.SchemaType) []Violation {

	schema := ctx.schema
	problems := make([]Violation, 0)

	if schema.MultipleOf != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/multiple-of", ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
	}

	if typeOf == structure.Number {
		if schema.Maximum != nil && schema.Minimum != nil {
			maximum := *schema.Maximum
			minimum := *schema.Minimum

			if minimum > maximum {
				problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.FIELD_MIN_MUST_NOT_BE_GT_MAX})
			}

		}

		if schema.Maximum == nil && schema.ExclusiveMaximum != nil && *schema.ExclusiveMaximum {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: message.REQUIRED_FIELD})
		}

		if schema.Minimum == nil && schema.ExclusiveMinimum != nil && *schema.ExclusiveMinimum {
			problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.REQUIRED_FIELD})
		}
	} else {
		problems = append(problems, doLintNumberSchemaBoundary(ctx, schema, typeOf)...)
	}

	return problems
}

func doLintNumberSchemaBoundary(ctx *LintContext, schema *structure.Schema, typeOf structure.SchemaType) []Violation {
	problems := make([]Violation, 0)

	if typeOf != structure.Number {

		if schema.Maximum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
		}

		if schema.Minimum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
		}

		if schema.ExclusiveMaximum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-maximum", ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
		}

		if schema.ExclusiveMinimum != nil {
			problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-minimum", ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
		}

	}

	return problems
}

func doLintArraySchema(ctx *LintContext, typeOf structure.SchemaType) []Violation {

	schema := ctx.schema
	problems := doLintArraySchemaLength(ctx, typeOf)

	if typeOf != structure.Array && schema.UniqueItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/unique-items", ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
	}

	if typeOf == structure.Array && schema.Items == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items", ctx.prefix), Message: message.REQUIRED_FIELD})
	}

	if typeOf == structure.Array && schema.Items != nil && schema.Items.TypeOf != nil && *schema.Items.TypeOf == structure.Array {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items/type", ctx.prefix), Message: message.ARRAY_FIELD_TYPE_NOT_ALLOWED})
		return problems
	}

	if ctx.schema.Items != nil {
		copyOf := &LintContext{prefix: ctx.prefix + "/items", schema: ctx.schema.Items}
		problems = append(problems, doLintSchema(copyOf)...)
	}

	if typeOf == structure.Array && schema.Format != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
	}

	return problems
}

func doLintArraySchemaLength(ctx *LintContext, typeOf structure.SchemaType) []Violation {
	schema := ctx.schema
	problems := make([]Violation, 0)

	if typeOf != structure.Array && schema.MaxItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
	}

	if typeOf != structure.Array && schema.MinItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
	}

	if typeOf == structure.Array && schema.MinItems != nil && *schema.MinItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.FIELD_MIN_ITEMS_GT_ZERO})
	}

	if typeOf == structure.Array && schema.MaxItems != nil && *schema.MaxItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: message.FIELD_MAX_ITEMS_GT_ZERO})
	}

	if typeOf == structure.Array && schema.MinItems != nil && schema.MaxItems != nil {
		maximum := *schema.MaxItems
		minimum := *schema.MinItems

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.FIELD_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS})
		}

	}

	return problems
}

func doLintParameter(ctx *LintContext, parameter *structure.Parameter) ([]Violation, error) {

	problems := make([]Violation, 0)
	problems = append(problems, doLintParameterFields(ctx, parameter)...)

	array, err := doLintParameterSchema(ctx, parameter)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	return problems, nil
}

func doLintParameterFields(ctx *LintContext, parameter *structure.Parameter) []Violation {
	problems := make([]Violation, 0)

	if parameter == nil {
		return problems
	}

	if parameter.Description == nil || utils.IsBlank(*parameter.Description) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/description", ctx.prefix), Message: message.REQUIRED_FIELD})
	}

	if parameter.Name == nil || utils.IsBlank(*parameter.Name) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/name", ctx.prefix), Message: message.REQUIRED_FIELD})
	}

	if parameter.Index != nil && *parameter.Index < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(index_format_pattern, ctx.prefix), Message: message.FIELD_INDEX_GT_ZERO})
	}

	if (parameter.In == nil || *parameter.In == structure.Flags) && parameter.Index != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(index_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
	}

	if parameter.In != nil && *parameter.In == structure.Arguments && parameter.Index == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(index_format_pattern, ctx.prefix), Message: message.FIELD_WHEN_IN_ARGUMENTS})
	}

	if parameter.In != nil && *parameter.In == structure.Arguments && parameter.ShortForm != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/short-form", ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
	}

	if parameter.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
	}

	return problems
}

func doLintParameterSchema(ctx *LintContext, parameter *structure.Parameter) ([]Violation, error) {

	problems := make([]Violation, 0)

	if parameter.Schema == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/schema", ctx.prefix), Message: message.REQUIRED_FIELD})
	}

	if parameter.Schema != nil {
		ctx := LintContext{prefix: fmt.Sprintf("%s/schema", ctx.prefix), schema: parameter.Schema}
		problems = append(problems, doLintSchema(&ctx)...)
	}

	return problems, nil
}

func doLintExit(ctx *LintContext, exit *structure.Exit) ([]Violation, error) {
	problems := make([]Violation, 0)

	if exit.Message == nil || utils.IsBlank(*exit.Message) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/message", ctx.prefix), Message: message.REQUIRED_FIELD})
	}

	if exit.Code == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/code", ctx.prefix), Message: message.REQUIRED_FIELD})
	}

	if exit.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED})
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
