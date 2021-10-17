package project

type CommandObject struct {
	Id          *string          `yaml:"id" json:"id"`
	Name        *string          `yaml:"name" json:"name"`
	Description *string          `yaml:"description" json:"description"`
	Subcommands []*CommandObject  `yaml:"commands" json:"commands"`
	Parameters  []ParameterObject `yaml:"parameters" json:"parameters"`
	Exit        []ExitObject      `yaml:"exit" json:"exit"`
}
