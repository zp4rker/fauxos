package commands

import (
	"fauxos/filesystem"
	"strings"
)

func RemoveDirectory(path string, fs *map[string]filesystem.Node, cwd string) (string, map[string]filesystem.Node) {
	var output string

	path = filesystem.ResolvePath(path, cwd)

	components := strings.Split(strings.Trim(path, "/"), "/")
	current := fs
	for i, component := range components {
		if i == len(components)-1 {
			if _, ok := (*current)[component]; ok {
				if _, ok := (*current)[component].(filesystem.File); ok {
					output += path + " is not a directory!\n"
				} else {
					// TODO: Add remove directory code here
					delete(*current, component)
				}
			} else {
				output += path + " does not exist!\n"
			}
		} else {
			if dir, ok := (*current)[component].(filesystem.Directory); ok {
				current = &dir.Files
			}
		}
	}

	return output, *fs
}
