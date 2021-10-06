package internal

import "fmt"

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

type Problem struct {
	Code    int
	Message string
}

func (instance *Problem) Error() string {
	return instance.Message
}
