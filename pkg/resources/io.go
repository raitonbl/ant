package resources

import "embed"

var (
	//go:embed schema.json
	resources embed.FS
)

func GetResource(filename string) ([]byte, error) {
	binary, err := resources.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return binary, nil
}
