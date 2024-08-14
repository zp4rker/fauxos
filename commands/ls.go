package commands

import (
	"fauxos/filesystem"
	"strings"
)

func ListFiles(path string, fs map[string]filesystem.Node) string {
	var output string

	components := strings.Split(path, "/")
	current := fs
	for i, component := range components {
		if i == len(components)-1 {
			switch f := current[component].(type) {
			case filesystem.Directory:
				for _, file := range f.Files {
					output += file.GetName()
					if _, ok := file.(filesystem.Directory); ok {
						output += "/"
					}

					output += "\n"
				}
			}
		} else {
			if dir, ok := current[component].(filesystem.Directory); ok {
				current = dir.Files
			}
		}
	}

	return output
}
