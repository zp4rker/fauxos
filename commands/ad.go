package commands

import (
	"fauxos/filesystem"
	"strings"
)

func AddDirectory(path string, fs *map[string]filesystem.Node, cwd string) (string, map[string]filesystem.Node) {
	var output string

	path = filesystem.ResolvePath(path, cwd)

	components := strings.Split(strings.Trim(path, "/"), "/")
	current := fs
	for i, component := range components {
		if i == len(components)-1 {
			if _, ok := (*current)[component]; ok {
				if _, ok := (*current)[component].(filesystem.File); ok {
					output += path + " already contains a file!"
				} else {
					output += path + " already exists!"
				}
			} else {
				(*current)[component] = filesystem.Directory{Name: component}
			}
		} else {
			if dir, ok := (*current)[component].(filesystem.Directory); ok {
				if dir.Files == nil {
					dir.Files = map[string]filesystem.Node{}
					(*current)[component] = dir
				}
				current = &dir.Files
			} else {
				output += path + " is not valid!"
				return output, *fs
			}
		}
	}

	return output, *fs
}
