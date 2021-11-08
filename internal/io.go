package internal

import (
	"github.com/magiconair/properties"
	"os"
)

func GetFile(path string) (*File, error) {

	if path == "" {
		return nil, GetProblemFactory().GetFileNotFound(path)
	}

	_, err := os.Stat(path)

	if err != nil {
		return nil, GetProblemFactory().GetFileCannotBeOpened(path, err)
	}

	binary, err := os.ReadFile(path)

	if err != nil {
		return nil, GetProblemFactory().GetFileCannotBeOpened(path, err)
	}

	return &File{path: path, content: binary}, nil
}

type File struct {
	path    string
	content []byte
}

func (instance *File) GetName() string {
	return instance.path
}

func (instance *File) GetContent() []byte {
	return instance.content
}

func ReadURI(uri string) ([]byte, error) {
	return nil, GetProblemFactory().NotImplemented()
}

func GetProperties(uri string) (*properties.Properties, error) {
	binary, err := ReadURI(uri)

	if err != nil {
		return nil, err
	}

	return properties.MustLoadString(string(binary)), nil
}
