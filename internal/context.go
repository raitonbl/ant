package internal

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/raitonbl/ant/internal/project"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
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

type LintContext interface {
	GetProjectFile() *File
	GetDocument() (*project.CliObject, error)
}

type DefaultContext struct {
	projectFile           *File
	processedDocument     *project.CliObject
	document              *project.CliObject
	targetProjectLocation string
	targetProjectLanguage LanguageType
	targetProjectType     ProjectType
	configuration         *properties.Properties
	directory             *string
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

type GenerateContext interface {
	GetDirectory() string
	GetProjectFile() *File
	GetTargetProjectType() ProjectType
	GetTargetProjectLocation() string
	GetProjectTargetLanguage() LanguageType
	GetDocument() (*project.CliObject, error)
	GetConfiguration(key string) *string
	BindConfiguration(value interface{}) error
	Write(filename string, binary []byte) error
	WriteTo(directory []string, filename string, binary []byte) error
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

func (instance *DefaultContext) GetConfiguration(key string) *string {

	if key == "" {
		return nil
	}

	if instance.configuration == nil {
		return nil
	}

	value, exits := instance.configuration.Get(key)

	if !exits {
		return nil
	}

	return &value
}

func (instance *DefaultContext) BindConfiguration(value interface{}) error {
	if instance.configuration == nil {
	}

	if err := instance.configuration.Decode(value); err != nil {
		return err
	}

	return nil
}

func (instance *DefaultContext) Write(filename string, binary []byte) error {
	return instance.WriteTo([]string{}, filename, binary)
}

func (instance *DefaultContext) WriteTo(directory []string, filename string, binary []byte) error {
	if instance.directory == nil {
		directory, err := ioutil.TempDir(fmt.Sprintf("%d", os.Getpid()), "ant-cli")

		if err != nil {
			return GetProblemFactory().GetProblem(err)
		}

		instance.directory = &directory
	}

	destination := make([]string, 0)
	destination = append(destination, *instance.directory)

	if directory != nil && len(directory) > 0 {
		destination = append(destination, directory...)
	}

	destination = append(destination, filename)

	err := os.WriteFile(path.Join(destination...), binary, 0700)

	if err != nil {
		return GetProblemFactory().GetProblem(err)
	}

	return nil
}

func (instance *DefaultContext) GetDirectory() string {

	if instance.directory == nil {
		directory, err := ioutil.TempDir(fmt.Sprintf("%d", os.Getpid()), "ant-cli")

		if err != nil {
			return path.Join(os.TempDir(), "ant-cli", fmt.Sprintf("%d", os.Getpid()))
		}

		instance.directory = &directory
	}

	return *instance.directory
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

type ContextFactory struct {
	projectDestination string
	filename           string
	projectLanguage    LanguageType
	projectType        ProjectType
	configuration      *properties.Properties
}

func (instance *ContextFactory) SetFilename(filename string) *ContextFactory {
	instance.filename = filename
	return instance
}

func (instance *ContextFactory) SetProjectLanguage(value LanguageType) *ContextFactory {
	instance.projectLanguage = value
	return instance
}

func (instance *ContextFactory) SetProjectType(value ProjectType) *ContextFactory {
	instance.projectType = value
	return instance
}

func (instance *ContextFactory) SetProjectDestination(value string) *ContextFactory {
	instance.projectDestination = value
	return instance
}

func (instance *ContextFactory) SetProperties(properties *properties.Properties) *ContextFactory {
	instance.configuration = properties
	return instance
}

func (instance *ContextFactory) GetLintContext() (LintContext, error) {
	file, err := GetFile(instance.filename)

	if err != nil {
		return nil, err
	}

	return &DefaultContext{projectFile: file}, nil
}

func (instance *ContextFactory) GetGenerateContext() (GenerateContext, error) {
	file, err := GetFile(instance.filename)

	if err != nil {
		return nil, err
	}

	if !isLangSupported(instance.projectType, instance.projectLanguage) {
		return nil, GetProblemFactory().GetUnsupportedLanguage(string(instance.projectType), string(instance.projectLanguage))
	}

	return &DefaultContext{projectFile: file, targetProjectLanguage: instance.projectLanguage,
		targetProjectLocation: instance.projectDestination, targetProjectType: instance.projectType,
		configuration: instance.configuration}, nil
}
