package project

type Specification struct {
	Name        *string     `yaml:"name" json:"name,omitempty"`
	Version     *string     `yaml:"version" json:"version,omitempty"`
	Subcommands []Command   `yaml:"commands" json:"commands,omitempty"`
	Description *string     `yaml:"description" json:"description,omitempty"`
	Parameters  []Parameter `yaml:"parameters" json:"parameters,omitempty"`
	Exit        []Exit      `yaml:"exit" json:"exit,omitempty"`
	Schemas     []*Schema    `yaml:"schemas" json:"schemas,omitempty"`
}
