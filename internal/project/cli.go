package project

type CliObject struct {
	Name        *string           `yaml:"name" json:"name,omitempty"`
	Version     *string           `yaml:"version" json:"version,omitempty"`
	Subcommands []CommandObject   `yaml:"commands" json:"commands,omitempty"`
	Description *string           `yaml:"description" json:"description,omitempty"`
	Components  *ComponentsObject  `yaml:"components" json:"components,omitempty"`
}
