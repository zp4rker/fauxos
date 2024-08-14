package filesystem

type File struct {
	Name     string
	Contents []byte
}

func (f File) GetName() string {
	return f.Name
}
