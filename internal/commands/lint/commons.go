package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"github.com/raitonbl/cli/internal/project/structure"
)

const (
	schema_format_pattern     = "%s/format"
	minimum_format_pattern    = "%s/minimum"
	min_length_format_pattern = "%s/min-length"
	min_items_format_pattern  = "%s/min-items"
)

type LintingContext struct {
	prefix string
	when   Moment
	schema *structure.Schema
}

func ValidateSchema(ctx *LintingContext) []Violation {

	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if schema.TypeOf == nil {
		return []Violation{{Path: fmt.Sprintf("%s/type", ctx.prefix), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when}}
	}

	typeOf := *schema.TypeOf

	if typeOf != structure.String && schema.Format != nil && (*schema.Format == structure.Date || *schema.Format == structure.DateTime || *schema.Format == structure.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	if typeOf == structure.String && schema.Format != nil && !(*schema.Format == structure.Date || *schema.Format == structure.DateTime || *schema.Format == structure.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.PAREMETER_FORMAT_NOT_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	problems = append(problems, ValidateTextSchema(ctx, typeOf)...)
	problems = append(problems, ValidateArraySchema(ctx, typeOf)...)
	problems = append(problems, ValidateNumberSchema(ctx, typeOf)...)

	if schema.Enum != nil && schema.Examples != nil && len(schema.Examples) > 0 {
		for i, example := range schema.Examples {
			if !belongsTo(schema.Enum, example) {
				problems = append(problems, Violation{Path: fmt.Sprintf("%s/examples/%d", ctx.prefix, i), Message: message.PARAMETER_EXAMPLE_MUST_BE_PART_OF_ENUM, Type: when})
			}
		}
	}

	return problems
}

func ValidateTextSchema(ctx *LintingContext, typeOf structure.SchemaType) []Violation {
	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if typeOf != structure.String && schema.MaxLength != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	if typeOf != structure.String && schema.MinLength != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	if typeOf == structure.String && schema.MinLength != nil && *schema.MinLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.PARAMETER_MIN_LENGTH_GT_ZERO, Type: when})
	}

	if typeOf == structure.String && schema.MaxLength != nil && *schema.MaxLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: message.PARAMETER_MAX_LENGTH_GT_ZERO, Type: when})
	}

	if typeOf == structure.String && schema.MaxLength != nil && schema.MinLength != nil {
		maximum := *schema.MaxLength
		minimum := *schema.MinLength

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.PARAMETER_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH, Type: when})
		}

	}

	if typeOf != structure.String && schema.Pattern != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/pattern", ctx.prefix), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	return problems
}

func ValidateNumberSchema(ctx *LintingContext, typeOf structure.SchemaType) []Violation {
	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if typeOf != structure.Number && schema.MultipleOf != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/multiple-of", ctx.prefix), Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: when})
	}

	if typeOf != structure.Number && schema.Maximum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: when})
	}

	if typeOf != structure.Number && schema.Minimum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: when})
	}

	if typeOf == structure.Number && schema.Maximum != nil && schema.Minimum != nil {
		maximum := *schema.Maximum
		minimum := *schema.Minimum

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.PARAMETER_MIN_MUST_NOT_BE_GT_MAX, Type: when})
		}

	}

	if typeOf != structure.Number && schema.ExclusiveMaximum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-maximum", ctx.prefix), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf != structure.Number && schema.ExclusiveMinimum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-minimum", ctx.prefix), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf == structure.Number && schema.Maximum == nil && schema.ExclusiveMaximum != nil && *schema.ExclusiveMaximum {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when})
	}

	if typeOf == structure.Number && schema.Minimum == nil && schema.ExclusiveMinimum != nil && *schema.ExclusiveMinimum {
		problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when})
	}

	return problems
}

func ValidateArraySchema(ctx *LintingContext, typeOf structure.SchemaType) []Violation {
	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if typeOf != structure.Array && schema.MaxItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf != structure.Array && schema.MinItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf != structure.Array && schema.UniqueItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/unique-items", ctx.prefix), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf == structure.Array && schema.Items == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items", ctx.prefix), Message: message.REQUIRED_PARAMETER_MESSAGE, Type: when})
	}

	if typeOf == structure.Array && schema.Format != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf == structure.Array && schema.MinItems != nil && *schema.MinItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.PARAMETER_MIN_ITEMS_GT_ZERO, Type: when})
	}

	if typeOf == structure.Array && schema.MaxItems != nil && *schema.MaxItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: message.PARAMETER_MAX_ITEMS_GT_ZERO, Type: when})
	}

	if typeOf == structure.Array && schema.MinItems != nil && schema.MaxItems != nil {
		maximum := *schema.MaxItems
		minimum := *schema.MinItems

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.PARAMETER_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS, Type: when})
		}

	}

	if typeOf == structure.Array && schema.Items != nil && schema.Items.TypeOf != nil && *schema.Items.TypeOf == structure.Array {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items/type", ctx.prefix), Message: message.PARAMETER_TYPE_NOT_ALLOWED, Type: when})
		return problems
	}

	if ctx.schema.Items != nil {
		copyOf := &LintingContext{prefix: ctx.prefix + "/items", when: when, schema: ctx.schema.Items}
		problems = append(problems, ValidateSchema(copyOf)...)
	}

	return problems
}

func belongsTo(array []string, value string) bool {
	for _, each := range array {
		if each == value {
			return true
		}

	}
	return false
}
