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
					output += path + " already contains a file!\n"
				} else {
					output += path + " already exists!\n"
				}
			} else {
				(*current)[component] = filesystem.Directory{Name: component, Files: map[string]filesystem.Node{}}
			}
		} else {
			if dir, ok := (*current)[component].(filesystem.Directory); ok {
				current = &dir.Files
			}
		}
	}

	return output, *fs
}
