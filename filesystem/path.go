package filesystem

import "strings"

func ResolvePath(path string, cwd string) string {
	if strings.HasPrefix(path, "..") {
		parts := strings.Split(strings.Trim(cwd, "/"), "/")
		parent := strings.Join(parts[:len(parts)-1], "/")
		path = "/" + parent + "/" + path[2:]
	}
	if path[0] != '/' {
		path = cwd + path
	}
	if path == "//" {
		path = "/"
	}
	return path
}
