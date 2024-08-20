package commands

import (
	"fauxos/filesystem"
	"strings"
)

func ChangeDirectory(path string, fs map[string]filesystem.Node, cwd string) (string, string) {
	var output string
	dir := "/"

	path = filesystem.ResolvePath(path, cwd)
	if path == "/" {
		return path, ""
	}

	components := strings.Split(strings.Trim(path, "/"), "/")
	current := fs
	for i, component := range components {
		if i == len(components)-1 {
			if d, ok := current[component].(filesystem.Directory); ok {
				dir += d.Name + "/"
			} else {
				output += path + " is not a valid directory!"
				return "", output
			}
		} else {
			if d, ok := current[component].(filesystem.Directory); ok {
				current = d.Files
				dir += d.GetName() + "/"
			} else {
				output += path + " is not a valid directory!"
				return "", output
			}
		}
	}

	return dir, output
}
