package internal

import (
	"encoding/json"
	"github.com/raitonbl/cli/internal/project/structure"
	"gopkg.in/yaml.v3"
	"strings"
)

type ProjectContext interface {
	GetProjectFile() *File
	GetDocument() (*structure.Specification, error)
}

func GetContext(filename string) (ProjectContext, error) {
	file, err := GetFile(filename)

	if err != nil {
		return nil, err
	}

	return &DefaultContext{projectFile: file}, nil
}

type DefaultContext struct {
	projectFile       *File
	processedDocument *structure.Specification
	document          *structure.Specification
}

func (instance *DefaultContext) GetProjectFile() *File {
	return instance.projectFile
}

func (instance *DefaultContext) GetDocument() (*structure.Specification, error) {

	if instance.document != nil {
		return instance.document, nil
	}

	if instance.GetProjectFile() == nil || instance.GetProjectFile().GetName() == "" {
		return nil, GetProblemFactory().GetConfigurationFileNotFound()
	}

	filename := instance.GetProjectFile().GetName()

	if !(strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml")) {
		return nil, GetProblemFactory().GetUnsupportedDescriptor()
	}

	binary := instance.GetProjectFile().GetContent()

	if strings.HasSuffix(filename, ".json") {
		return parseJson(binary)
	}

	return parseYaml(binary)
}

func parseYaml(binary []byte) (*structure.Specification, error) {

	if binary == nil {
		return nil, GetProblemFactory().GetUnexpectedContext()
	}

	descriptor := structure.Specification{}
	err := yaml.Unmarshal(binary, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, err
}

func parseJson(binary []byte) (*structure.Specification, error) {

	if binary == nil {
		return nil, GetProblemFactory().GetUnexpectedContext()
	}

	descriptor := structure.Specification{}
	err := json.Unmarshal(binary, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, err
}