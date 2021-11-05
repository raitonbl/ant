package project

type ExitObject struct {
	Code        *int    `yaml:"code" json:"code,omitempty"`
	Message     *string `yaml:"message" json:"message,omitempty"`
	RefersTo    *string `yaml:"refers-to" json:"refers-to,omitempty"`
	Description *string `yaml:"description" json:"description,omitempty"`
}
