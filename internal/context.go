package internal

type ProjectContext interface {
	GetProjectFile() *File
}

type DefaultContext struct {
	ProjectFile *File
}

func (instance *DefaultContext) GetProjectFile() *File {
	return instance.ProjectFile
}
