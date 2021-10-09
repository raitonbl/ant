package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"testing"
)

func TestLintParameterIndex003(t *testing.T) {
	withValue(t, "index-003.json", Document, func(array []Violation) {
		if len(array) != 0 {
			t.Fatal(fmt.Sprintf("\nExpected:[]\nActual:[%s]", toText(array)))
		}
	})
}

func TestLintParameter_Where_name_Is_Missing(t *testing.T) {
	withSingleValue(t, "index-005.json", Document, &Violation{Path: "/parameters/0/name", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_id_Is_Missing(t *testing.T) {
	withSingleValue(t, "index-007.json", Document, &Violation{Path: "/parameters/1/id", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_schema_Is_Missing(t *testing.T) {
	object := &Violation{Path: "/parameters/0/schema", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document}
	withSingleValue(t, "index-008.json", Document, object)
}

func TestLintParameter_Where_index_Is_Negative(t *testing.T) {
	withSingleValue(t, "index-009.json", Document, &Violation{Path: "/parameters/0/index", Message: message.PARAMETER_INDEX_GT_ZERO_MESSAGE, Type: Document})
}

func TestLintParameter_Where_refersTo_Has_Value(t *testing.T) {
	withSingleValue(t, "index-010.json", Document, &Violation{Path: "/parameters/1/refers-to", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_description_Is_Missing(t *testing.T) {
	withSingleValue(t, "index-006.json", Document, &Violation{Path: "/parameters/1/description", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_index_Is_Missing_and_argument_Is_Null(t *testing.T) {
	withSingleValue(t, "index-004.json", Document, &Violation{Path: "/parameters/0/index", Message: message.REQUIRED_PARAMETER_FIELD_WHEN_IN_ARGUMENTS, Type: Document})
}

func TestLintParameter_Where_schema_type_Is_Missing(t *testing.T) {
	withSingleValue(t, "index-011.json", Document, &Violation{Path: "/parameters/0/schema/type", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_type_Is_string_and_format_is_int64(t *testing.T) {
	withSingleValue(t, "index-012.json", Document, &Violation{Path: "/parameters/0/schema/format", Message: message.PAREMETER_FORMAT_NOT_ALLOWED_IN_TYPE_STRING, Type: Document})
}

func TestLintParameter_Where_type_Is_number_and_format_is_date(t *testing.T) {
	withSingleValue(t, "index-013.json", Document, &Violation{Path: "/parameters/0/schema/format", Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: Document})
}

func TestLintParameter_Where_type_Is_number_and_format_is_datetime(t *testing.T) {
	withSingleValue(t, "index-014.json", Document, &Violation{Path: "/parameters/0/schema/format", Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: Document})
}

func TestLintParameter_Where_type_enum_and_example_not_part_of_enum(t *testing.T) {
	withSingleValue(t, "index-015.json", Document, &Violation{Path: "/parameters/1/schema/examples/0", Message: message.PARAMETER_EXAMPLE_MUST_BE_PART_OF_ENUM, Type: Document})
}

func TestLintParameter_Where_type_number_and_max_length_is_two(t *testing.T) {
	withSingleValue(t, "index-016.json", Document, &Violation{Path: "/parameters/2/schema/max-length", Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: Document})
}

func TestLintParameter_Where_type_number_and_min_length_is_two(t *testing.T) {
	withSingleValue(t, "index-017.json", Document, &Violation{Path: "/parameters/2/schema/min-length", Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING, Type: Document})
}

func TestLintParameter_Where_type_string_and_min_length_is_gt_max_length(t *testing.T) {
	withSingleValue(t, "index-018.json", Document, &Violation{Path: "/parameters/0/schema/min-length", Message: message.PARAMETER_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH, Type: Document})
}

func TestLintParameter_Where_type_number_and_pattern_defined(t *testing.T) {
	withSingleValue(t, "index-019.json", Document, &Violation{Path: "/parameters/2/schema/pattern", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_type_string_and_min_length_lt_zero(t *testing.T) {
	withSingleValue(t, "index-020.json", Document, &Violation{Path: "/parameters/0/schema/min-length", Message: message.PARAMETER_MIN_LENGTH_GT_ZERO, Type: Document})
}

func TestLintParameter_Where_type_string_and_max_length_lt_zero(t *testing.T) {
	withSingleValue(t, "index-021.json", Document, &Violation{Path: "/parameters/0/schema/max-length", Message: message.PARAMETER_MAX_LENGTH_GT_ZERO, Type: Document})
}

func TestLintParameter_Where_type_string_and_multiple_of(t *testing.T) {
	withSingleValue(t, "index-022.json", Document, &Violation{Path: "/parameters/0/schema/multiple-of", Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: Document})
}

func TestLintParameter_Where_type_string_and_maximum(t *testing.T) {
	withSingleValue(t, "index-023.json", Document, &Violation{Path: "/parameters/0/schema/maximum", Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: Document})
}

func TestLintParameter_Where_type_string_and_minimum(t *testing.T) {
	withSingleValue(t, "index-024.json", Document, &Violation{Path: "/parameters/0/schema/minimum", Message: message.PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER, Type: Document})
}

func TestLintParameter_Where_type_number_and_minimum_gt_maximum(t *testing.T) {
	withSingleValue(t, "index-025.json", Document, &Violation{Path: "/parameters/2/schema/minimum", Message: message.PARAMETER_MIN_MUST_NOT_BE_GT_MAX, Type: Document})
}

func TestLintParameter_Where_type_string_and_exclusive_minimum_and_minimum_missing(t *testing.T) {
	withSingleValue(t, "index-026.json", Document, &Violation{Path: "/parameters/0/schema/exclusive-minimum", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_type_string_and_exclusive_maximum_and_maximum_missing(t *testing.T) {
	withSingleValue(t, "index-027.json", Document, &Violation{Path: "/parameters/0/schema/exclusive-maximum", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_type_number_and_exclusive_minimum_and_minimum_missing(t *testing.T) {
	withSingleValue(t, "index-028.json", Document, &Violation{Path: "/parameters/2/schema/minimum", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_type_number_and_exclusive_maximum_and_maximum_missing(t *testing.T) {
	withSingleValue(t, "index-029.json", Document, &Violation{Path: "/parameters/2/schema/maximum", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_type_number_and_max_items(t *testing.T) {
	withSingleValue(t, "index-030.json", Document, &Violation{Path: "/parameters/0/schema/max-items", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_type_number_and_min_items(t *testing.T) {
	withSingleValue(t, "index-031.json", Document, &Violation{Path: "/parameters/0/schema/min-items", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_type_number_and_unique_items(t *testing.T) {
	withSingleValue(t, "index-032.json", Document, &Violation{Path: "/parameters/0/schema/unique-items", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_type_array_and_array_schema_is_undefined(t *testing.T) {
	withSingleValue(t, "index-033.json", Document, &Violation{Path: "/parameters/2/schema/items", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_type_array_and_format_not_defined(t *testing.T) {
	withSingleValue(t, "index-034.json", Document, &Violation{Path: "/parameters/2/schema/format", Message: message.PARAMETER_FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintParameter_Where_type_array_and_min_items_gt_max_items(t *testing.T) {
	withSingleValue(t, "index-035.json", Document, &Violation{Path: "/parameters/2/schema/min-items", Message: message.PARAMETER_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS, Type: Document})
}

func TestLintParameter_Where_type_array_and_min_items_lt_zero(t *testing.T) {
	withSingleValue(t, "index-037.json", Document, &Violation{Path: "/parameters/2/schema/min-items", Message: message.PARAMETER_MIN_ITEMS_GT_ZERO, Type: Document})
}

func TestLintParameter_Where_type_array_and_max_items_lt_zero(t *testing.T) {
	withSingleValue(t, "index-036.json", Document, &Violation{Path: "/parameters/2/schema/max-items", Message: message.PARAMETER_MAX_ITEMS_GT_ZERO, Type: Document})
}

func TestLintParameter_Where_type_array_and_array_type_undefined(t *testing.T) {
	withSingleValue(t, "index-038.json", Document, &Violation{Path: "/parameters/2/schema/items/type", Message: message.REQUIRED_PARAMETER_MESSAGE, Type: Document})
}

func TestLintParameter_Where_type_array_and_array_type_array(t *testing.T) {
	withSingleValue(t, "index-039.json", Document, &Violation{Path: "/parameters/2/schema/items/type", Message: message.PARAMETER_TYPE_NOT_ALLOWED, Type: Document})
}

func toString(each *Violation) string {

	if each == nil {
		return "nil"
	}

	return fmt.Sprintf("{\"path\":\"%s\" , \"message\":\"%s\" , \"type\": \"%s\"}", each.Path, each.Message, each.Type)

}

func toText(array []Violation) string {

	if array == nil || len(array) == 0 {
		return "[]"
	}

	if len(array) == 1 {
		return fmt.Sprintf("[%s]", toString(&array[0]))
	}

	text := "["

	for index, value := range array {
		if index != 0 {
			text += ","
		}

		text += toString(&value)
	}

	text += "]"

	return text
}

func withSingleValue(t *testing.T, filename string, when Moment, object *Violation) {
	withValue(t, filename, when, func(array []Violation) {

		if array == nil || len(array) == 0 {
			t.Fatal(fmt.Sprintf("\nExpected:[%s]\nActual:[nil]", toString(object)))
		}

		if len(array) > 1 {
			t.Fatal(fmt.Sprintf("\nExpected:[%s]\nActual:%s", toString(object), toText(array)))
		}

		singleValue := &array[0]

		if singleValue.Path != object.Path || singleValue.Message != object.Message || singleValue.Type != object.Type {
			t.Fatal(fmt.Sprintf("\nExpected:[%s]\nActual:[%s]", toString(object), toString(singleValue)))
		}

	})
}

func withValue(t *testing.T, filename string, when Moment, fn func(array []Violation)) {

	ctx, err := internal.GetContext(fmt.Sprintf("testdata/%s", filename))

	if err != nil {
		t.Fatal(err)
	}

	instance := ParameterSectionLinter{}

	document, err := ctx.GetDocument()

	if err != nil {
		t.Fatal(err)
	}

	array, err := instance.Lint(ctx, document, when)

	if err != nil {
		t.Fatal(err)
	}

	fn(array)
}
