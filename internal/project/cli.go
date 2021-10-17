package project

type CliObject struct {
	Name        *string     `yaml:"name" json:"name,omitempty"`
	Version     *string         `yaml:"version" json:"version,omitempty"`
	Subcommands []CommandObject `yaml:"commands" json:"commands,omitempty"`
	Description *string           `yaml:"description" json:"description,omitempty"`
	Parameters  []ParameterObject `yaml:"parameters" json:"parameters,omitempty"`
	Exit        []ExitObject      `yaml:"exit" json:"exit,omitempty"`
	Schemas     []*Schema         `yaml:"schemas" json:"schemas,omitempty"`
}
