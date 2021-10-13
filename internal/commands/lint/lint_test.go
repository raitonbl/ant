package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"testing"
)

func TestLintFromJson(t *testing.T) {
	doLintTest(t, "index-003.json", nil)
}

func TestLintFromYaml(t *testing.T) {
	doLintTest(t, "index-003.yaml", nil)
}

func TestLint_where_name_is_missing(t *testing.T) {
	doLintTest(t, "index-005.json", &Violation{Path: "/parameters/0/name", Message: message.REQUIRED_FIELD})
}

func TestLint_where_in_flags_and_index_is_defined(t *testing.T) {
	doLintTest(t, "index-002.json", &Violation{Path: "/parameters/1/index", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_in_arguments_and_shortForm_is_defined(t *testing.T) {
	doLintTest(t, "index-001.json", &Violation{Path: "/parameters/0/short-form", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_schema_is_missing(t *testing.T) {
	object := &Violation{Path: "/parameters/0/schema", Message: message.REQUIRED_FIELD}
	doLintTest(t, "index-008.json", object)
}

func TestLint_where_index_Is_Negative(t *testing.T) {
	doLintTest(t, "index-009.json", &Violation{Path: "/parameters/0/index", Message: message.FIELD_INDEX_GT_ZERO})
}

func TestLint_where_description_is_missing(t *testing.T) {
	doLintTest(t, "index-006.json", &Violation{Path: "/parameters/1/description", Message: message.REQUIRED_FIELD})
}

func TestLint_where_index_is_missing_and_argument_Is_Null(t *testing.T) {
	doLintTest(t, "index-004.json", &Violation{Path: "/parameters/0/index", Message: message.FIELD_WHEN_IN_ARGUMENTS})
}

func TestLint_where_type_Is_string_and_format_is_int64(t *testing.T) {
	doLintTest(t, "index-012.json", &Violation{Path: "/parameters/0/schema/format", Message: message.FIELD_FORMAT_NOT_ALLOWED_IN_TYPE_STRING})
}

func TestLint_where_type_Is_number_and_format_is_date(t *testing.T) {
	doLintTest(t, "index-013.json", &Violation{Path: "/parameters/0/schema/format", Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
}

func TestLint_where_type_Is_number_and_format_is_datetime(t *testing.T) {
	doLintTest(t, "index-014.json", &Violation{Path: "/parameters/0/schema/format", Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
}

func TestLint_where_type_enum_and_example_not_part_of_enum(t *testing.T) {
	doLintTest(t, "index-015.json", &Violation{Path: "/parameters/1/schema/examples/0", Message: message.FIELD_EXAMPLE_MUST_BE_PART_OF_ENUM})
}

func TestLint_where_type_number_and_max_length_is_two(t *testing.T) {
	doLintTest(t, "index-016.json", &Violation{Path: "/parameters/2/schema/max-length", Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
}

func TestLint_where_type_number_and_min_length_is_two(t *testing.T) {
	doLintTest(t, "index-017.json", &Violation{Path: "/parameters/2/schema/min-length", Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING})
}

func TestLint_where_type_string_and_min_length_is_gt_max_length(t *testing.T) {
	doLintTest(t, "index-018.json", &Violation{Path: "/parameters/0/schema/min-length", Message: message.FIELD_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH})
}

func TestLint_where_type_string_and_min_length_lt_zero(t *testing.T) {
	doLintTest(t, "index-020.json", &Violation{Path: "/parameters/0/schema/min-length", Message: message.FIELD_MIN_LENGTH_GT_ZERO})
}

func TestLint_where_type_string_and_max_length_lt_zero(t *testing.T) {
	doLintTest(t, "index-021.json", &Violation{Path: "/parameters/0/schema/max-length", Message: message.FIELD_MAX_LENGTH_GT_ZERO})
}

func TestLint_where_type_string_and_multiple_of(t *testing.T) {
	doLintTest(t, "index-022.json", &Violation{Path: "/parameters/0/schema/multiple-of", Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
}

func TestLint_where_type_string_and_maximum(t *testing.T) {
	doLintTest(t, "index-023.json", &Violation{Path: "/parameters/0/schema/maximum", Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
}

func TestLint_where_type_string_and_minimum(t *testing.T) {
	doLintTest(t, "index-024.json", &Violation{Path: "/parameters/0/schema/minimum", Message: message.FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER})
}

func TestLint_where_type_number_and_minimum_gt_maximum(t *testing.T) {
	doLintTest(t, "index-025.json", &Violation{Path: "/parameters/2/schema/minimum", Message: message.FIELD_MIN_MUST_NOT_BE_GT_MAX})
}

func TestLint_where_type_number_and_max_items(t *testing.T) {
	doLintTest(t, "index-030.json", &Violation{Path: "/parameters/0/schema/max-items", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_type_number_and_min_items(t *testing.T) {
	doLintTest(t, "index-031.json", &Violation{Path: "/parameters/0/schema/min-items", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_type_array_and_array_schema_is_undefined(t *testing.T) {
	doLintTest(t, "index-033.json", &Violation{Path: "/parameters/2/schema/items", Message: message.REQUIRED_FIELD})
}

func TestLint_where_type_array_and_format_not_defined(t *testing.T) {
	doLintTest(t, "index-034.json", &Violation{Path: "/parameters/2/schema/format", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_type_array_and_min_items_gt_max_items(t *testing.T) {
	doLintTest(t, "index-035.json", &Violation{Path: "/parameters/2/schema/min-items", Message: message.FIELD_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS})
}

func TestLint_where_type_array_and_min_items_lt_zero(t *testing.T) {
	doLintTest(t, "index-037.json", &Violation{Path: "/parameters/2/schema/min-items", Message: message.FIELD_MIN_ITEMS_GT_ZERO})
}

func TestLint_where_type_array_and_max_items_lt_zero(t *testing.T) {
	doLintTest(t, "index-036.json", &Violation{Path: "/parameters/2/schema/max-items", Message: message.FIELD_MAX_ITEMS_GT_ZERO})
}

func TestLint_where_type_array_and_array_type_array(t *testing.T) {
	doLintTest(t, "index-039.json", &Violation{Path: "/parameters/2/schema/items/type", Message: message.ARRAY_FIELD_TYPE_NOT_ALLOWED})
}

func TestLintCommand_where_command_exit_code_is_missing(t *testing.T) {
	doLintTest(t, "index-046.json", &Violation{Path: "/commands/0/exit/0/code", Message: message.REQUIRED_FIELD})
}

func TestLintCommand_where_command_exit_has_id(t *testing.T) {
	doLintTest(t, "index-047.json", &Violation{Path: "/commands/0/exit/0", Message: "did not match any of the specified OneOf schemas"})
}

func TestLintCommand_where_command_exit_message_is_missing(t *testing.T) {
	doLintTest(t, "index-048.json", &Violation{Path: "/commands/0/exit/0/message", Message: message.REQUIRED_FIELD})
}

func TestLintCommand_where_command_exit_refers_to_is_unresolvable(t *testing.T) {
	doLintTest(t, "index-049.json", &Violation{Path: "/commands/0/exit/1/refers-to", Message: message.UNRESOLVABLE_FIELD})
}

func TestLintCommand_where_command_parameter_refers_to_is_unresolvable(t *testing.T) {
	doLintTest(t, "index-050.json", &Violation{Path: "/commands/0/parameters/0/refers-to", Message: message.UNRESOLVABLE_FIELD})
}

func TestLintCommand_where_command_parameter_in_is_null_and_index_is_defined(t *testing.T) {
	doLintTest(t, "index-051.json", &Violation{Path: "/commands/0/parameters/0/refers-to", Message: message.FIELD_NOT_ALLOWED})
}

func TestLintCommand_where_command_parameter_in_flags_and_index_is_defined(t *testing.T) {
	doLintTest(t, "index-052.json", &Violation{Path: "/commands/0/parameters/0/refers-to", Message: message.FIELD_NOT_ALLOWED})
}

func TestLintExit_where_id_is_missing(t *testing.T) {
	doLintTest(t, "index-040.json", &Violation{Path: "/exit/0", Message: "\"id\" value is required"})
}

func TestLintExit_where_id_is_blank(t *testing.T) {
	doLintTest(t, "index-042.json", &Violation{Path: "/exit/0/id", Message: message.REQUIRED_FIELD})
}

func TestLintExit_where_code_is_missing(t *testing.T) {
	doLintTest(t, "index-041.json", &Violation{Path: "/exit/0", Message: "\"code\" value is required"})
}

func TestLintExit_where_message_is_missing(t *testing.T) {
	doLintTest(t, "index-043.json", &Violation{Path: "/exit/0", Message: "\"message\" value is required"})
}

func TestLintExit_where_message_is_blank(t *testing.T) {
	doLintTest(t, "index-044.json", &Violation{Path: "/exit/0/message", Message: message.REQUIRED_FIELD})
}

func TestLintExit_where_refers_to_is_defined(t *testing.T) {
	doLintTest(t, "index-045.json", &Violation{Path: "/exit/0", Message: "additional properties are not allowed"})
}

func TestLint_where_type_number_and_pattern_defined(t *testing.T) {
	doLintTest(t, "index-019.json", &Violation{Path: "/parameters/2/schema/pattern", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_type_string_and_exclusive_minimum_and_minimum_missing(t *testing.T) {
	doLintTest(t, "index-026.json", &Violation{Path: "/parameters/0/schema/exclusive-minimum", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_type_string_and_exclusive_maximum_and_maximum_missing(t *testing.T) {
	doLintTest(t, "index-027.json", &Violation{Path: "/parameters/0/schema/exclusive-maximum", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_type_number_and_exclusive_minimum_and_minimum_missing(t *testing.T) {
	doLintTest(t, "index-028.json", &Violation{Path: "/parameters/2/schema/minimum", Message: message.REQUIRED_FIELD})
}

func TestLint_where_type_number_and_exclusive_maximum_and_maximum_missing(t *testing.T) {
	doLintTest(t, "index-029.json", &Violation{Path: "/parameters/2/schema/maximum", Message: message.REQUIRED_FIELD})
}

func TestLint_where_type_number_and_unique_items(t *testing.T) {
	doLintTest(t, "index-032.json", &Violation{Path: "/parameters/0/schema/unique-items", Message: message.FIELD_NOT_ALLOWED})
}

func TestLint_where_type_array_and_array_type_undefined(t *testing.T) {
	doLintTest(t, "index-038.json", &Violation{Path: "/parameters/2/schema/items", Message: "\"type\" value is required"})
}

func doLintTest(t *testing.T, filename string, object *Violation) {
	doLintFrom(t, filename, func(array []Violation) {

		if (array == nil || len(array) == 0) && object != nil {
			t.Fatal(fmt.Sprintf("\nExpected:[%s]\nActual:[nil]", toString(object)))
		}

		if len(array) > 1 {
			t.Fatal(fmt.Sprintf("\nExpected:[%s]\nActual:%s", toString(object), toText(array)))
		}

		if object == nil && len(array) == 0 {
			return
		}

		if object == nil {
			t.Fatal(fmt.Sprintf("\nExpected:[]\nActual:%s", toText(array)))
		}

		singleValue := &array[0]

		if singleValue.Path != object.Path || singleValue.Message != object.Message  {
			t.Fatal(fmt.Sprintf("\nExpected:[%s]\nActual:[%s]", toString(object), toString(singleValue)))
		}

	})
}

func doLintFrom(t *testing.T, filename string, afterLint func(array []Violation)) {

	ctx, err := internal.GetContext(fmt.Sprintf("testdata/%s", filename))

	if err != nil {
		t.Fatal(err)
	}

	array, err := Lint(ctx)

	if err != nil {
		t.Fatal(err)
	}

	afterLint(array)
}
