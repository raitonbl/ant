package internal

import "testing"

func TestGetContext_where_json_exits(t *testing.T) {
	ctx, err := GetContext("commands/lint/testdata/index-003.json")

	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.GetDocument()

	if err != nil {
		t.Fatal(err)
	}

	if ctx.GetProjectFile() == nil {
		t.Fatal("GetProjectFile() returned nil")
	}

}

func TestGetContext_where_yaml_exits(t *testing.T) {
	ctx, err := GetContext("commands/lint/testdata/index-003.yaml")

	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.GetDocument()

	if err != nil {
		t.Fatal(err)
	}

	if ctx.GetProjectFile() == nil {
		t.Fatal("GetProjectFile() returned nil")
	}

}

func TestGetContext_where_json_doesnt_exit(t *testing.T) {
	_, err := GetContext("commands/lint/testdata/index.json")

	if err == nil {
		t.Fatal("error not caught")
	}

}

func TestGetContext_where_yaml_doesnt_exit(t *testing.T) {
	_, err := GetContext("commands/lint/testdata/index.yaml")

	if err == nil {
		t.Fatal("error not caught")
	}

}
