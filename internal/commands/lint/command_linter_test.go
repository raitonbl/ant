package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"testing"
)

func TestLintCommand(t *testing.T) {
	doLintCommandSectionLinter(t, "index-003.json", Document, nil)
}

func TestLintCommand_where_command_exit_code_is_missing(t *testing.T) {
	doLintCommandSectionLinter(t, "index-046.json", Document, &Violation{Path: "/commands/0/exit/0/code", Message: message.REQUIRED_FIELD, Type: Document})
}

func TestLintCommand_where_command_exit_has_id(t *testing.T) {
	doLintCommandSectionLinter(t, "index-047.json", Document, &Violation{Path: "/commands/0/exit/0/id", Message: message.FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintCommand_where_command_exit_message_is_missing(t *testing.T) {
	doLintCommandSectionLinter(t, "index-048.json", Document, &Violation{Path: "/commands/0/exit/0/message", Message: message.REQUIRED_FIELD, Type: Document})
}

func TestLintCommand_where_command_exit_refers_to_is_unresolvable(t *testing.T) {
	doLintCommandSectionLinter(t, "index-049.json", Document, &Violation{Path: "/commands/0/exit/1/refers-to", Message: message.UNRESOLVABLE_FIELD, Type: Document})
}

func TestLintCommand_where_command_parameter_refers_to_is_unresolvable(t *testing.T) {
	doLintCommandSectionLinter(t, "index-050.json", Document, &Violation{Path: "/commands/0/parameters/0/refers-to", Message: message.UNRESOLVABLE_FIELD, Type: Document})
}

func TestLintCommand_where_command_parameter_in_is_null_and_index_is_defined(t *testing.T) {
	doLintCommandSectionLinter(t, "index-051.json", Document, &Violation{Path: "/commands/0/parameters/0/index", Message: message.FIELD_NOT_ALLOWED, Type: Document})
}

func TestLintCommand_where_command_parameter_in_flags_and_index_is_defined(t *testing.T) {
	doLintCommandSectionLinter(t, "index-052.json", Document, &Violation{Path: "/commands/0/parameters/0/index", Message: message.FIELD_NOT_ALLOWED, Type: Document})
}

func doLintCommandSectionLinter(t *testing.T, filename string, when Moment, object *Violation) {
	doLintFile(t, filename, when, &CommandLinter{}, func(array []Violation) {

		if object == nil && array != nil && len(array) == 0 {
			return
		}

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
