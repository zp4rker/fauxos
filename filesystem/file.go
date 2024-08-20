package filesystem

type File struct {
	Name string
	Data []byte
}

func (f File) GetName() string {
	return f.Name
}
