package lint

import (
	"context"
	"encoding/json"
	"github.com/qri-io/jsonschema"
	"github.com/raitonbl/ant/internal"
	"github.com/raitonbl/ant/internal/project"
	"github.com/raitonbl/ant/pkg/resources"
	"sigs.k8s.io/yaml"
	"strings"
)

type Violation struct {
	Path    string
	Message string
}

type CommandLintingContext struct {
	path           string
	commandCache   map[string]*project.Command
	exitCache      map[string]*project.Exit
	parameterCache map[string]*project.Parameter
	schemaCache    map[string]*project.Schema
}

func Lint(context internal.ProjectContext) ([]Violation, error) {

	if context == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedContext()
	}

	if context.GetProjectFile() == nil || context.GetProjectFile().GetName() == "" {
		return nil, internal.GetProblemFactory().GetConfigurationFileNotFound()
	}

	problems, err := doLint(context)

	if err != nil {
		return nil, err
	}

	return problems, nil
}

func doLint(context internal.ProjectContext) ([]Violation, error) {

	binary := make([]byte, 0)
	problems := make([]Violation, 0)

	if strings.HasSuffix(context.GetProjectFile().GetName(), ".json") {
		binary = context.GetProjectFile().GetContent()
	} else if strings.HasSuffix(context.GetProjectFile().GetName(), ".yaml") || strings.HasSuffix(context.GetProjectFile().GetName(), ".yml") {
		binary = context.GetProjectFile().GetContent()
		content, err := yaml.YAMLToJSON(binary)

		if err != nil {
			return nil, err
		}

		binary = content
	} else {
		return nil, internal.GetProblemFactory().GetProblem("the specified doesn't meet the expected extension[json|yaml|yml]")
	}

	array, err := doLintFile(binary)

	if err != nil {
		return nil, err
	}

	if len(array) > 0 {
		return array, nil
	}

	array, err = doLintObject(context)

	if err != nil {
		return nil, err
	}

	return append(problems, array...), nil
}

func doLintFile(binary []byte) ([]Violation, error) {
	goContext := context.Background()

	schema, err := resources.GetResource("schema.json")

	if err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	rs := &jsonschema.Schema{}

	if err := json.Unmarshal(schema, rs); err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	errs, err := rs.ValidateBytes(goContext, binary)

	if err != nil {
		return nil, internal.GetProblemFactory().GetProblem(err)
	}

	problems := make([]Violation, len(errs))

	for index, each := range errs {
		problems[index] = Violation{Path: each.PropertyPath, Message: each.Message}
	}

	return problems, nil
}

func doLintObject(ctx internal.ProjectContext) ([]Violation, error) {

	document, err := ctx.GetDocument()

	if err != nil {
		return nil, err
	}

	if document == nil {
		return nil, internal.GetProblemFactory().GetUnexpectedState()
	}

	problems := make([]Violation, 0)

	schemaCache, array, err := doLintSchemaSection(document)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	if document.Parameters == nil {
		return problems, nil
	}

	parameterCache, array, err := doLintParameterSection(document, schemaCache)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	exitCache, array, err := doLintExitSection(document)

	if err != nil {
		return nil, err
	}

	problems = append(problems, array...)

	array, err = doLintCommandSection(document, parameterCache, exitCache, schemaCache)

	if err != nil {
		return nil, err
	}

	return append(problems, array...), nil
}
