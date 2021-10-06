package internal

var factory *ProblemFactory

func GetProblemFactory() *ProblemFactory {

	if factory == nil {
		factory = &ProblemFactory{}
	}

	return factory
}

type ProblemFactory struct {
}

func (d *ProblemFactory) GetUnexpectedContext() *Problem {
	return &Problem{Code: 1, Message: "Unexpect context"}
}

func (d *ProblemFactory) GetConfigurationFileNotFound() error {
	return &Problem{Code: 2, Message: "Unexpect context"}
}

type Problem struct {
	Code    int
	Message string
}

func (instance *Problem) Error() string {
	return instance.Message
}
