package lint

import (
	"fmt"
	"github.com/raitonbl/cli/internal"
	"github.com/raitonbl/cli/internal/project/structure"
)

type LinterBuilder interface {
	Build() (Linter, error)
}

type DelegatedLinterBuilder struct {
	sequence []Linter
}

func (instance *DelegatedLinterBuilder) Build() (Linter, error) {

	if instance.sequence == nil || len(instance.sequence) == 0 {
		return nil, internal.GetProblemFactory().GetProblem("can only lint if at least one (1) is available")
	}

	for index, each := range instance.sequence {
		if each == nil {
			return nil, internal.GetProblemFactory().GetProblem(fmt.Sprintf("linter[%d] mustn't be null", index))
		}
	}

	return &DelegatedLinter{sequence: instance.sequence}, nil
}

func (instance *DelegatedLinterBuilder) Append(object Linter) *DelegatedLinterBuilder {

	if object == nil {
		return instance
	}

	if instance.sequence == nil {
		instance.sequence = make([]Linter, 0)
	}

	instance.sequence = append(instance.sequence, object)

	return instance
}

type Moment string

const (
	Binary   Moment = "binary"
	Document Moment = "document"
)

type Linter interface {
	CanLint(ctx internal.ProjectContext, when Moment) bool
	Lint(ctx internal.ProjectContext, document *structure.Specification, when Moment) ([]Violation, error)
}

type Violation struct {
	Path    string
	Message string
	Type    Moment
}

type DelegatedLinter struct {
	sequence []Linter
}

func (instance *DelegatedLinter) CanLint(ctx internal.ProjectContext, when Moment) bool {
	return ctx != nil
}

func (instance *DelegatedLinter) Lint(ctx internal.ProjectContext, object *structure.Specification, when Moment) ([]Violation, error) {
	problems := make([]Violation, 0)

	for _, each := range instance.sequence {
		if each.CanLint(ctx, when) {
			array, err := each.Lint(ctx, object, when)

			if err != nil {
				return nil, err
			}

			if array != nil && len(array) > 0 {
				problems = append(problems, array...)
			}
		}
	}

	return problems, nil
}
