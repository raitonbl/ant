package project

import (
	"encoding/json"
	"fmt"
	"github.com/raitonbl/cli/internal"
	"testing"
)

func TestLoad(t *testing.T) {
	rv, err := Load(&internal.DefaultContext{Descriptor: "../../testdata/load/descriptor.yaml"})

	if err != nil {
		t.Fatal(err)
	}

	if rv == nil {
		t.Fatal("rv cannot be null")
	}

	binary, _ := json.Marshal(rv)


	fmt.Println(string(binary))
}
