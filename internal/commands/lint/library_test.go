package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"testing"
)

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

func doLintFile(t *testing.T, filename string, when Moment, instance Linter, afterLint func(array []Violation)) {

	ctx, err := internal.GetContext(fmt.Sprintf("testdata/%s", filename))

	if err != nil {
		t.Fatal(err)
	}

	document, err := ctx.GetDocument()

	if err != nil {
		t.Fatal(err)
	}

	array, err := instance.Lint(ctx, document, when)

	if err != nil {
		t.Fatal(err)
	}

	afterLint(array)
}
