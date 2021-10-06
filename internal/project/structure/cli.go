package structure

type Cli struct {
	Name        *string     `yaml:"name"`
	Version     *string     `yaml:"version"`
	Subcommands []Command   `yaml:"commands"`
	Description *string     `yaml:"description"`
	Parameters  []Parameter `yaml:"parameters"`
	Exit        []Exit      `yaml:"exit"`
}
