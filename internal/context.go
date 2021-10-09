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
	GetFinalDocument() (*structure.Specification, error)
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

func (instance *DefaultContext) GetFinalDocument() (*structure.Specification, error) {

	if instance.processedDocument != nil {
		return instance.processedDocument, nil
	}

	fromFile, err := instance.GetDocument()

	if err != nil {
		return nil, err
	}

	if fromFile == nil {
		return nil, GetProblemFactory().GetUnexpectedState()
	}

	return finalOf(fromFile)
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

func finalOf(spec *structure.Specification) (*structure.Specification, error) {
	document := structure.Specification{Name: spec.Name, Version: spec.Version, Description: spec.Description}

	exitCache := make(map[string]structure.Exit)
	parameterCache := make(map[string]structure.Parameter)

	if spec.Exit != nil && len(spec.Exit) > 0 {
		for _, each := range spec.Exit {
			id := *each.Id
			exitCache[id] = each
		}
	}

	if spec.Parameters != nil && len(spec.Parameters) > 0 {
		for _, each := range spec.Parameters {
			parameterCache[*each.Id] = each
		}
	}

	if spec.Subcommands == nil {
		return &document, nil
	}

	commands, err := resolvedCopyOf(spec.Subcommands, parameterCache, exitCache)

	if err != nil {
		return nil, err
	}

	document.Subcommands = commands

	return &document, nil
}

func resolvedCopyOf(commands []*structure.Command, parameterCache map[string]structure.Parameter, exitCache map[string]structure.Exit) ([]*structure.Command, error) {
	array := make([]*structure.Command, 0)

	for _, each := range commands {
		v, err := resolveAndCopy(each, parameterCache, exitCache)

		if err != nil {
			return nil, err
		}

		array = append(array, v)
	}

	return array, nil
}

func resolveAndCopy(command *structure.Command, parameterCache map[string]structure.Parameter, exitCache map[string]structure.Exit) (*structure.Command, error) {

	if command == nil {
		return nil, nil
	}

	copyOf := structure.Command{Id: command.Id, Name: command.Name, Description: command.Description}
	copyOf.Exit = make([]structure.Exit, 0)
	copyOf.Parameters = make([]structure.Parameter, 0)
	copyOf.Subcommands = make([]*structure.Command, 0)

	if command.Exit != nil {

		err := resolveAndCopyExit(command, exitCache, copyOf)

		if err != nil {
			return nil, err
		}

	}

	if command.Parameters != nil {
		err := resolveAndCopyParameters(command, parameterCache, copyOf)
		if err != nil {
			return nil, err
		}
	}

	if command.Subcommands != nil {
		err := resolveAndCopySubcommands(command, parameterCache, exitCache, copyOf)
		if err != nil {
			return nil, err
		}
	}

	return &copyOf, nil
}

func resolveAndCopyExit(command *structure.Command, exitCache map[string]structure.Exit, instance structure.Command) error {
	for _, exit := range command.Exit {
		if exit.RefersTo != nil {
			name := *exit.RefersTo
			if fromCache, ok := exitCache[name]; ok {
				instance.Exit = append(instance.Exit, fromCache)
			} else {
				return GetProblemFactory().GetMissingExit(name)
			}
		} else {
			instance.Exit = append(instance.Exit, exit)
		}
	}
	return nil
}

func resolveAndCopyParameters(command *structure.Command, parameterCache map[string]structure.Parameter, copyOf structure.Command) error {
	for _, parameter := range command.Parameters {
		if parameter.RefersTo != nil {
			name := *parameter.RefersTo
			if fromCache, ok := parameterCache[name]; ok {
				copyOf.Parameters = append(copyOf.Parameters, fromCache)
			} else {
				return GetProblemFactory().GetMissingParameter(name)
			}
		} else {
			copyOf.Parameters = append(copyOf.Parameters, parameter)
		}
	}
	return nil
}

func resolveAndCopySubcommands(command *structure.Command, parameterCache map[string]structure.Parameter, exitCache map[string]structure.Exit, copyOf structure.Command) error {
	for _, subCommand := range command.Subcommands {
		v, err := resolveAndCopy(subCommand, parameterCache, exitCache)

		if err != nil {
			return err
		}

		copyOf.Subcommands = append(copyOf.Subcommands, v)
	}
	return nil
}
