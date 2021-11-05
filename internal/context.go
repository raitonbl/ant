package internal

import (
	"encoding/json"
	"github.com/raitonbl/ant/internal/project"
	"gopkg.in/yaml.v3"
	"strings"
)

type ProjectType string
type LanguageType string

const (
	TestsType       ProjectType  = "tests"
	ApplicationType ProjectType  = "application"
	GoLang          LanguageType = "golang"
	Python3         LanguageType = "python3"
)

type ProjectContext interface {
	GetProjectFile() *File
	GetDocument() (*project.CliObject, error)
}

func GetContext(filename string) (ProjectContext, error) {
	file, err := GetFile(filename)

	if err != nil {
		return nil, err
	}

	return &DefaultContext{projectFile: file}, nil
}

type DefaultContext struct {
	projectFile           *File
	processedDocument     *project.CliObject
	document              *project.CliObject
	targetProjectLocation string
	targetProjectLanguage LanguageType
	targetProjectType     ProjectType
}

func (instance *DefaultContext) GetProjectFile() *File {
	return instance.projectFile
}

func (instance *DefaultContext) GetDocument() (*project.CliObject, error) {

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

func parseYaml(binary []byte) (*project.CliObject, error) {

	if binary == nil {
		return nil, GetProblemFactory().GetUnexpectedContext()
	}

	descriptor := project.CliObject{}
	err := yaml.Unmarshal(binary, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, err
}

func parseJson(binary []byte) (*project.CliObject, error) {

	if binary == nil {
		return nil, GetProblemFactory().GetUnexpectedContext()
	}

	descriptor := project.CliObject{}
	err := json.Unmarshal(binary, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, err
}

type GenerateProjectContext interface {
	GetProjectFile() *File
	GetDocument() (*project.CliObject, error)
	GetTargetProjectType() ProjectType
	GetTargetProjectLocation() string
	GetProjectTargetLanguage() LanguageType
}

func (instance *DefaultContext) GetTargetProjectType() ProjectType {
	return instance.targetProjectType
}

func (instance *DefaultContext) GetTargetProjectLocation() string {
	return instance.targetProjectLocation
}

func (instance *DefaultContext) GetProjectTargetLanguage() LanguageType {
	return instance.targetProjectLanguage
}

func GetGenerateProjectContext(filename string, targetDirectory string, targetType ProjectType, targetLanguage LanguageType) (GenerateProjectContext, error) {
	file, err := GetFile(filename)

	if err != nil {
		return nil, err
	}

	if !isLangSupported(targetType, targetLanguage) {
		return nil, GetProblemFactory().GetUnsupportedLanguage(string(targetType), string(targetLanguage))
	}

	return &DefaultContext{projectFile: file, targetProjectLanguage: targetLanguage, targetProjectLocation: targetDirectory, targetProjectType: targetType}, nil
}

func isLangSupported(projectType ProjectType, projectLanguage LanguageType) bool {

	if projectType == ApplicationType && projectLanguage == GoLang {
		return true
	}

	if projectType == TestsType && projectLanguage == Python3 {
		return true
	}

	return false
}
