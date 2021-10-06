package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"testing"
)

func TestLint(t *testing.T) {
	file, _ := internal.GetFile("testdata/index-01.json")
	ctx := &internal.DefaultContext{ProjectFile: file}

	object := JsonSchemaLinter{}

	if !object.CanLint(ctx, Binary) {
		t.Fatal("schema linter must lint binary")
	}

	if object.CanLint(ctx, Document) {
		t.Fatal("schema linter mustn't lint document")
	}

	problems, err := object.Lint(ctx, nil, Binary)

	if err != nil {
		t.Fatal(err)
	}

	if len(problems) > 0 {
		t.Fatal(fmt.Sprintf("unexpected errors:%s", problems))
	}

}

func TestLintForInvalidContent(t *testing.T) {
	file, _ := internal.GetFile("testdata/index-02.json")
	ctx := &internal.DefaultContext{ProjectFile: file}

	object := JsonSchemaLinter{}

	if !object.CanLint(ctx, Binary) {
		t.Fatal("schema linter must lint binary")
	}

	if object.CanLint(ctx, Document) {
		t.Fatal("schema linter mustn't lint document")
	}

	problems, err := object.Lint(ctx, nil, Binary)

	if err != nil {
		t.Fatal(err)
	}

	if len(problems) == 0 {
		t.Fatal("expected one (1) error but got none")
	}

	if len(problems) > 1 {
		t.Fatal(fmt.Sprintf("unexpect errors:%s", problems))
	}

	problem := problems[0]

	if !(problem.Path == "/commands/0/exit/0/code" && problem.Message == "type should be integer, got string" && problem.Type == Binary) {
		t.Fatal(fmt.Sprintf("unexpect error:%s", problem))
	}

}
