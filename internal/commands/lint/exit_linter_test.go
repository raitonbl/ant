package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal/commands/lint/message"
	"testing"
)

func TestLintExitIndex003(t *testing.T) {
	doLintFile(t, "index-003.json", Document, &ExitLinter{}, func(array []Violation) {
		if len(array) != 0 {
			t.Fatal(fmt.Sprintf("\nExpected:[]\nActual:[%s]", toText(array)))
		}
	})
}

func TestLintExit_Where_id_is_missing(t *testing.T) {
	doLintExitSectionLinter(t, "index-040.json", Document, &Violation{Path: "/exit/0/id", Message: message.REQUIRED_FIELD_MESSAGE, Type: Document})
}

func TestLintExit_Where_id_is_blank(t *testing.T) {
	doLintExitSectionLinter(t, "index-042.json", Document, &Violation{Path: "/exit/0/id", Message: message.REQUIRED_FIELD_MESSAGE, Type: Document})
}

func TestLintExit_Where_code_is_missing(t *testing.T) {
	doLintExitSectionLinter(t, "index-041.json", Document, &Violation{Path: "/exit/0/code", Message: message.REQUIRED_FIELD_MESSAGE, Type: Document})
}

func TestLintExit_Where_message_is_missing(t *testing.T) {
	doLintExitSectionLinter(t, "index-043.json", Document, &Violation{Path: "/exit/0/message", Message: message.REQUIRED_FIELD_MESSAGE, Type: Document})
}

func TestLintExit_Where_message_is_blank(t *testing.T) {
	doLintExitSectionLinter(t, "index-044.json", Document, &Violation{Path: "/exit/0/message", Message: message.REQUIRED_FIELD_MESSAGE, Type: Document})
}

func TestLintExit_Where_refers_to_is_defined(t *testing.T) {
	doLintExitSectionLinter(t, "index-045.json", Document, &Violation{Path: "/exit/0/refers-to", Message: message.FIELD_NOT_ALLOWED, Type: Document})
}

func doLintExitSectionLinter(t *testing.T, filename string, when Moment, object *Violation) {
	doLintFile(t, filename, when, &ExitLinter{}, func(array []Violation) {

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