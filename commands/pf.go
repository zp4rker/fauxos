package commands

import (
	"fauxos/filesystem"
	"strings"
)

func PrintFile(path string, fs map[string]filesystem.Node, cwd string) string {
	var output string

	path = filesystem.ResolvePath(path, cwd)

	components := strings.Split(strings.Trim(path, "/"), "/")
	current := fs
	for i, component := range components {
		if i == len(components)-1 {
			switch f := current[component].(type) {
			case filesystem.File:
				output += string(f.Contents) + "\n"
			default:
				output += path + " is not a valid file!\n"
			}
		} else {
			if dir, ok := current[component].(filesystem.Directory); ok {
				current = dir.Files
			}
		}
	}

	return output
}
