package structure

type In string

const (
	Flags     In = "flags"
	Arguments In = "arguments"
)

type Parameter struct {
	Id           string  `yaml:"id"`
	In           *In     `yaml:"in"`
	Index        *int    `yaml:"index"`
	Required     *bool   `yaml:"required"`
	Name         *string `yaml:"name"`
	ShortForm   *string `yaml:"short-form"`
	Description  *string `yaml:"description"`
	RefersTo     *string `yaml:"refers-to"`
	DefaultValue *string `yaml:"default"`
	Schema       *Schema `yaml:"schema"`
}
