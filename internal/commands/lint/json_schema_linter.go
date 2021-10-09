package lint

import (
	"context"
	"encoding/json"
	"github.com/qri-io/jsonschema"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/project/structure"
	"os"
	"strings"
)

type JsonSchemaLinter struct {
	schema []byte
}

func (instance *JsonSchemaLinter) CanLint(ctx internal.ProjectContext, when Moment) bool {
	return ctx != nil && strings.HasSuffix(ctx.GetProjectFile().GetName(), ".json") && when == Binary
}

func (instance *JsonSchemaLinter) Lint(ctx internal.ProjectContext, document *structure.Specification, when Moment) ([]Violation, error) {
	goContext := context.Background()

	if instance.schema == nil {
		binary, err := os.ReadFile("schema.json")

		if err != nil {
			return nil, internal.GetProblemFactory().GetProblem(err)
		}

		instance.schema = binary
	}

	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(instance.schema, rs); err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	errs, err := rs.ValidateBytes(goContext, ctx.GetProjectFile().GetContent())

	if err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	problems := make([]Violation, len(errs))

	for index, each := range errs {
		problems[index] = Violation{Path: each.PropertyPath, Message: each.Message, Type: when}
	}

	return problems, nil
}
