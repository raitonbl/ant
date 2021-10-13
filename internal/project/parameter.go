package project


type In string

const (
	Flags     In = "flags"
	Arguments In = "arguments"
)

type Parameter struct {
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
