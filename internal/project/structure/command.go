package structure

type Command struct {
	Id          *string     `yaml:"id" json:"id"`
	Name        *string     `yaml:"name" json:"name"`
	Description *string     `yaml:"description" json:"description"`
	Subcommands []*Command   `yaml:"commands" json:"commands"`
	Parameters  []Parameter `yaml:"parameters" json:"parameters"`
	Exit        []Exit      `yaml:"exit" json:"exit"`
}
