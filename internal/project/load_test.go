package project

import (
	"encoding/json"
	"fmt"
	"github.com/raitonbl/cli/internal"
	"testing"
)

func TestLoad(t *testing.T) {
	file, _ := internal.GetFile("../../testdata/load/descriptor.yaml")

	rv, err := Load(&internal.DefaultContext{ProjectFile: file})

	if err != nil {
		t.Fatal(err)
	}

	if rv == nil {
		t.Fatal("rv cannot be null")
	}

	binary, _ := json.Marshal(rv)

	fmt.Println(string(binary))
}
