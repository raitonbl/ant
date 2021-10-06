package structure

type Exit struct {
	Code     int    `yaml:"code"`
	Message  string `yaml:"message"`
	Id       string `yaml:"id"`
	RefersTo string `yaml:"refers-to"`
}
