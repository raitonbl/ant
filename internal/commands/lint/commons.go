package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"github.com/raitonbl/cli/internal/project/structure"
	"github.com/raitonbl/cli/internal/utils"
)

const (
	schema_format_pattern          = "%s/format"
	minimum_format_pattern         = "%s/minimum"
	min_length_format_pattern      = "%s/min-length"
	min_items_format_pattern       = "%s/min-items"
	parameter_index_format_pattern = "%s/index"
	refers_to_format_pattern       = "%s/refers-to"
)

type LintingContext struct {
	isLocal  bool
	prefix   string
	when     Moment
	schema   *structure.Schema
	document *structure.Specification
}

func lintSchema(ctx *LintingContext) []Violation {

	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if schema.TypeOf == nil {
		return []Violation{{Path: fmt.Sprintf("%s/type", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when}}
	}

	typeOf := *schema.TypeOf

	if typeOf != structure.String && schema.Format != nil && (*schema.Format == structure.Date || *schema.Format == structure.DateTime || *schema.Format == structure.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	if typeOf == structure.String && schema.Format != nil && !(*schema.Format == structure.Date || *schema.Format == structure.DateTime || *schema.Format == structure.Binary) {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_NOT_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	problems = append(problems, lintTextSchema(ctx, typeOf)...)
	problems = append(problems, lintArraySchema(ctx, typeOf)...)
	problems = append(problems, lintNumberSchema(ctx, typeOf)...)

	if schema.Enum != nil && schema.Examples != nil && len(schema.Examples) > 0 {
		for i, example := range schema.Examples {
			if !belongsTo(schema.Enum, example) {
				problems = append(problems, Violation{Path: fmt.Sprintf("%s/examples/%d", ctx.prefix, i), Message: message.FIELD_EXAMPLE_MUST_BE_PART_OF_ENUM, Type: when})
			}
		}
	}

	return problems
}

func lintTextSchema(ctx *LintingContext, typeOf structure.SchemaType) []Violation {
	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if typeOf != structure.String && schema.MaxLength != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	if typeOf != structure.String && schema.MinLength != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: when})
	}

	if typeOf == structure.String && schema.MinLength != nil && *schema.MinLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.FIELD_MIN_LENGTH_GT_ZERO, Type: when})
	}

	if typeOf == structure.String && schema.MaxLength != nil && *schema.MaxLength < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-length", ctx.prefix), Message: message.FIELD_MAX_LENGTH_GT_ZERO, Type: when})
	}

	if typeOf == structure.String && schema.MaxLength != nil && schema.MinLength != nil {
		maximum := *schema.MaxLength
		minimum := *schema.MinLength

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_length_format_pattern, ctx.prefix), Message: message.FIELD_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH, Type: when})
		}

	}

	if typeOf != structure.String && schema.Pattern != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/pattern", ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	return problems
}

func lintNumberSchema(ctx *LintingContext, typeOf structure.SchemaType) []Violation {
	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if typeOf != structure.Number && schema.MultipleOf != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/multiple-of", ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: when})
	}

	if typeOf != structure.Number && schema.Maximum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: when})
	}

	if typeOf != structure.Number && schema.Minimum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: when})
	}

	if typeOf == structure.Number && schema.Maximum != nil && schema.Minimum != nil {
		maximum := *schema.Maximum
		minimum := *schema.Minimum

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.FIELD_MIN_MUST_NOT_BE_GT_MAX, Type: when})
		}

	}

	if typeOf != structure.Number && schema.ExclusiveMaximum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-maximum", ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf != structure.Number && schema.ExclusiveMinimum != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/exclusive-minimum", ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf == structure.Number && schema.Maximum == nil && schema.ExclusiveMaximum != nil && *schema.ExclusiveMaximum {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/maximum", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if typeOf == structure.Number && schema.Minimum == nil && schema.ExclusiveMinimum != nil && *schema.ExclusiveMinimum {
		problems = append(problems, Violation{Path: fmt.Sprintf(minimum_format_pattern, ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	return problems
}

func lintArraySchema(ctx *LintingContext, typeOf structure.SchemaType) []Violation {
	when := ctx.when
	schema := ctx.schema

	problems := make([]Violation, 0)

	if typeOf != structure.Array && schema.MaxItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf != structure.Array && schema.MinItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf != structure.Array && schema.UniqueItems != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/unique-items", ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf == structure.Array && schema.Items == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if typeOf == structure.Array && schema.Format != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(schema_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	if typeOf == structure.Array && schema.MinItems != nil && *schema.MinItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.FIELD_MIN_ITEMS_GT_ZERO, Type: when})
	}

	if typeOf == structure.Array && schema.MaxItems != nil && *schema.MaxItems < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/max-items", ctx.prefix), Message: message.FIELD_MAX_ITEMS_GT_ZERO, Type: when})
	}

	if typeOf == structure.Array && schema.MinItems != nil && schema.MaxItems != nil {
		maximum := *schema.MaxItems
		minimum := *schema.MinItems

		if minimum > maximum {
			problems = append(problems, Violation{Path: fmt.Sprintf(min_items_format_pattern, ctx.prefix), Message: message.FIELD_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS, Type: when})
		}

	}

	if typeOf == structure.Array && schema.Items != nil && schema.Items.TypeOf != nil && *schema.Items.TypeOf == structure.Array {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/items/type", ctx.prefix), Message: message.ARRAY_FIELD_TYPE_NOT_ALLOWED, Type: when})
		return problems
	}

	if ctx.schema.Items != nil {
		copyOf := &LintingContext{prefix: ctx.prefix + "/items", when: when, schema: ctx.schema.Items}
		problems = append(problems, lintSchema(copyOf)...)
	}

	return problems
}

func lintParameter(ctx *LintingContext, parameter structure.Parameter, when Moment) ([]Violation, error) {

	problems := make([]Violation, 0)
	isValue := isValueParameter(parameter)

	if !isValue && ctx.isLocal {

		fromReference := findParameterById(ctx.document, parameter.RefersTo)

		if fromReference == nil {
			return []Violation{{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: message.UNRESOLVABLE_FIELD, Type: when}}, nil
		}

		if parameter.Index != nil && fromReference.In == nil {
			return []Violation{{Path: fmt.Sprintf(parameter_index_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when}}, nil
		}

		if parameter.Index != nil && fromReference.In != nil && *fromReference.In != structure.Arguments {
			return []Violation{{Path: fmt.Sprintf(parameter_index_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when}}, nil
		}

		return problems, nil
	}

	problems = append(problems, lintParameterInformation(ctx, parameter, when)...)

	if isValue && parameter.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	array, err := doLintSchema(ctx, parameter, when)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	return problems, nil
}

func lintParameterInformation(ctx *LintingContext, parameter structure.Parameter, when Moment) []Violation {
	problems := make([]Violation, 0)

	if !ctx.isLocal && (parameter.Id == nil || utils.IsBlank(*parameter.Id)) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if parameter.Description == nil || utils.IsBlank(*parameter.Description) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/description", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if parameter.Name == nil || utils.IsBlank(*parameter.Name) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/name", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if parameter.Index != nil && *parameter.Index < 0 {
		problems = append(problems, Violation{Path: fmt.Sprintf(parameter_index_format_pattern, ctx.prefix), Message: message.FIELD_INDEX_GT_ZERO, Type: when})
	}

	if parameter.In == nil || *parameter.In == structure.Flags && parameter.Index != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(parameter_index_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED_IN_FLAGS, Type: when})
	}

	if parameter.In != nil && *parameter.In == structure.Arguments && parameter.Index == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(parameter_index_format_pattern, ctx.prefix), Message: message.FIELD_WHEN_IN_ARGUMENTS, Type: when})
	}

	if parameter.In != nil && *parameter.In == structure.Arguments && parameter.ShortForm != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/short-form", ctx.prefix), Message: message.FIELD_NOT_ALLOWED_IN_ARGUMENTS, Type: when})
	}

	return problems
}

func lintExit(ctx *LintingContext, exit structure.Exit, when Moment) ([]Violation, error) {
	problems := make([]Violation, 0)
	isValue := isValueExit(exit)

	if !isValue && ctx.isLocal {

		fromReference := findExitById(ctx.document, exit.RefersTo)

		if fromReference == nil {
			return []Violation{{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: message.UNRESOLVABLE_FIELD, Type: when}}, nil
		}

		return problems, nil
	}

	if ctx.isLocal && exit.Id != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	} else if !ctx.isLocal && (exit.Id == nil || utils.IsBlank(*exit.Id)) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/id", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if exit.Message == nil || utils.IsBlank(*exit.Message) {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/message", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if exit.Code == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/code", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if exit.RefersTo != nil {
		problems = append(problems, Violation{Path: fmt.Sprintf(refers_to_format_pattern, ctx.prefix), Message: message.FIELD_NOT_ALLOWED, Type: when})
	}

	return problems, nil
}

func doLintSchema(ctx *LintingContext, parameter structure.Parameter, when Moment) ([]Violation, error) {

	problems := make([]Violation, 0)

	if parameter.Schema == nil {
		problems = append(problems, Violation{Path: fmt.Sprintf("%s/schema", ctx.prefix), Message: message.REQUIRED_FIELD, Type: when})
	}

	if parameter.Schema != nil {
		ctx := LintingContext{prefix: fmt.Sprintf("%s/schema", ctx.prefix), schema: parameter.Schema, when: when}
		problems = append(problems, lintSchema(&ctx)...)
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

func isValueParameter(parameter structure.Parameter) bool {

	if parameter.Id != nil {
		return true
	}

	if parameter.Description != nil {
		return true
	}

	if parameter.Name != nil {
		return true
	}

	if parameter.In != nil {
		return true
	}

	if parameter.Required != nil {
		return true
	}

	if parameter.ShortForm != nil {
		return true
	}

	if parameter.DefaultValue != nil {
		return true
	}

	if parameter.Schema != nil {
		return true
	}

	if parameter.Index != nil && parameter.RefersTo == nil {
		return true
	}

	return false
}

func isValueExit(exit structure.Exit) bool {

	if exit.Id != nil {
		return true
	}

	if exit.Description != nil {
		return true
	}

	if exit.Code != nil {
		return true
	}

	if exit.Message != nil {
		return true
	}

	return false
}

func findParameterById(document *structure.Specification, idRef *string) *structure.Parameter {

	if document.Parameters == nil || idRef == nil {
		return nil
	}

	id := *idRef

	for _, each := range document.Parameters {
		if each.Id != nil && *each.Id == id {
			return &each
		}
	}

	return nil
}

func findExitById(document *structure.Specification, idRef *string) *structure.Exit {

	if document.Exit == nil || idRef == nil {
		return nil
	}

	id := *idRef

	for _, each := range document.Exit {
		if each.Id != nil && *each.Id == id {
			return &each
		}
	}

	return nil
}
