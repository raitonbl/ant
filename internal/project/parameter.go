package project

type In string

const (
	Flags     In = "flags"
	Arguments In = "arguments"
)

type ParameterObject struct {
	Id           *string `yaml:"id" json:"id,omitempty"`
	In           *In     `yaml:"in" json:"in,omitempty"`
	Index        *int    `yaml:"index" json:"index,omitempty"`
	Required     *bool   `yaml:"required" json:"required,omitempty"`
	Name         *string `yaml:"name" json:"name,omitempty"`
	ShortForm    *string `yaml:"short-form" json:"short-form,omitempty"`
	Description  *string `yaml:"description" json:"description,omitempty"`
	RefersTo     *string `yaml:"refers-to" json:"refers-to,omitempty"`
	DefaultValue *string `yaml:"default" json:"default,omitempty"`
	Schema       *Schema `yaml:"schema" json:"schema,omitempty"`
}

func (instance *ParameterObject) Clone() *ParameterObject {
	object := ParameterObject{}

	if instance.Id != nil {
		object.Id = instance.Id
	}

	if instance.In != nil {
		object.In = instance.In
	}

	if instance.Index != nil {
		object.Index = instance.Index
	}

	if instance.Required != nil {
		object.Required = instance.Required
	}

	if instance.Name != nil {
		object.Name = instance.Name
	}

	if instance.ShortForm != nil {
		object.ShortForm = instance.ShortForm
	}

	if instance.Description != nil {
		object.Description = instance.Description
	}

	if instance.RefersTo != nil {
		object.RefersTo = instance.RefersTo
	}

	if instance.DefaultValue != nil {
		object.DefaultValue = instance.DefaultValue
	}

	if instance.Schema != nil {
		object.Schema = instance.Schema
	}

	return &object
}
