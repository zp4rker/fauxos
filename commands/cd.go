package commands

import (
	"fauxos/filesystem"
	"strings"
)

func ChangeDirectory(path string, fs map[string]filesystem.Node, cwd string) (string, string) {
	var output string
	dir := "/"

	if strings.HasPrefix(path, "..") {
		parts := strings.Split(strings.Trim(cwd, "/"), "/")
		parent := strings.Join(parts[:len(parts)-1], "/")
		path = "/" + parent + "/" + path[2:]
	}
	if path[0] != '/' {
		path = cwd + path
	}
	if path == "/" || path == "//" {
		return "/", ""
	}

	components := strings.Split(strings.Trim(path, "/"), "/")
	current := fs
	for i, component := range components {
		if i == len(components)-1 {
			if d, ok := current[component].(filesystem.Directory); ok {
				dir += d.Name + "/"
			} else {
				output += path + " is not a valid directory!\n"
				return "", output
			}
		} else {
			if d, ok := current[component].(filesystem.Directory); ok {
				current = d.Files
				dir += d.GetName() + "/"
			} else {
				output += path + " is not a valid directory!\n"
				return "", output
			}
		}
	}

	return dir, output
}
