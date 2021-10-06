package internal

import (
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
