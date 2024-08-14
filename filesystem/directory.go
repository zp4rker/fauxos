package filesystem

type Directory struct {
	Name  string
	Files map[string]Node
}

func (d Directory) GetName() string {
	return d.Name
}
