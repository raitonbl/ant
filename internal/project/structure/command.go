package structure

type Command struct {
	Name        *string     `yaml:"name"`
	Description *string     `yaml:"description"`
	Subcommands []Command   `yaml:"commands"`
	Parameters  []Parameter `yaml:"parameters"`
	Exit        []Exit      `yaml:"exit"`
}
