package project

type ComponentsObject struct {
	Schemas    map[string]*Schema          `yaml:"schemas" json:"schemas,omitempty"`
	Exits      map[string]*ExitObject      `yaml:"exits" json:"exits,omitempty"`
	Parameters map[string]*ParameterObject `yaml:"parameters" json:"parameters,omitempty"`
}
