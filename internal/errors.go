package internal

import (
	"fmt"
	"github.com/raitonbl/ant/internal/commands/lint"
)

var factory *ProblemFactory

func GetProblemFactory() *ProblemFactory {

	if factory == nil {
		factory = &ProblemFactory{}
	}

	return factory
}

type ProblemFactory struct {
}

func (instance *ProblemFactory) GetUnexpectedContext() *Problem {
	return &Problem{Code: 1, Message: "unexpected context"}
}

func (instance *ProblemFactory) GetProblem(value interface{}) *Problem {
	return &Problem{Code: 1, Message: fmt.Sprintf("unexpected problem occurred\ncaused by:%s", value)}
}

func (instance *ProblemFactory) GetConfigurationFileNotFound() error {
	return &Problem{Code: 2, Message: "specification not found"}
}

func (instance *ProblemFactory) GetUnsupportedDescriptor() error {
	return &Problem{Code: 3, Message: "unsupported descriptor"}
}

func (instance *ProblemFactory) GetFileNotFound(path string) error {
	return &Problem{Code: 1, Message: fmt.Sprintf("file '%s' cannot be found", path)}
}

func (instance *ProblemFactory) GetFileCannotBeOpened(path string, error error) error {
	return &Problem{Code: 1, Message: fmt.Sprintf("file '%s' cannot be opened\ncaused by:%s", path, error)}
}

func (instance *ProblemFactory) GetUnexpectedState() error {
	return &Problem{Code: 1, Message: "unexpected application state"}
}

func (instance *ProblemFactory) GetMissingExit(name string) error {
	return &Problem{Code: 101, Message: fmt.Sprintf("missing exit[\"id\":\"%s\"]", name)}
}

func (instance *ProblemFactory) GetMissingParameter(name string) error {
	return &Problem{Code: 101, Message: fmt.Sprintf("missing parameters[\"id\":\"%s\"]", name)}
}

func (instance *ProblemFactory) GetUnsupportedLanguage(projectType string, projectLang string) error {
	return &Problem{Code: 103, Message: fmt.Sprintf("language \"%s\" isn't supported for project of type \"%s\"", projectLang, projectType)}
}

func (instance *ProblemFactory) GetValidationConstraintViolation(problems []lint.Violation) error {
	txt := ""
	for index, each := range problems {
		txt += fmt.Sprintf("%d.path:%s\n message:%s", index, each.Path, each.Message)
	}

	return &Problem{Code: 104, Message: txt}
}

func (instance *ProblemFactory) NotImplemented() error {
	return instance.GetUnexpectedContext()
}

type Problem struct {
	Code    int
	Message string
}

func (instance *Problem) Error() string {
	return instance.Message
}
