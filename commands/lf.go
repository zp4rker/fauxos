package commands

import (
	"fauxos/filesystem"
	"sort"
	"strings"
)

func ListFiles(path string, fs map[string]filesystem.Node, cwd string) string {
	var output string

	path = filesystem.ResolvePath(path, cwd)
	if path == "/" {
		for k, v := range fs {
			output += k
			if _, ok := v.(filesystem.Directory); ok {
				output += "/"
			}
			output += "\n"
		}
		return output
	}

	components := strings.Split(strings.Trim(path, "/"), "/")
	current := fs
	for i, component := range components {
		if i == len(components)-1 {
			switch f := current[component].(type) {
			case filesystem.Directory:
				for _, name := range sortKeys(f.Files) {
					file := f.Files[name]
					output += file.GetName()
					if _, ok := file.(filesystem.Directory); ok {
						output += "/"
					}

					output += "\n"
				}
			default:
				output += path + " is not a valid directory!\n"
			}
		} else {
			if dir, ok := current[component].(filesystem.Directory); ok {
				current = dir.Files
			}
		}
	}

	return output
}

func sortKeys(m map[string]filesystem.Node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
