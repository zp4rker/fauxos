package main

import (
	"fauxos/filesystem"
	"fauxos/screens"
	"fauxos/styles"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

func main() {
	splash := styles.Ansi[8].Render(strings.Repeat("=", 26))
	splash += "\n    Running " + styles.Ansi[11].Copy().Bold(true).Render("fOS") + " v0.0.1"
	splash += "\n        by " + styles.Ansi[4].Copy().Bold(true).Render("zp4rker")
	splash += "\n" + styles.Ansi[8].Render(strings.Repeat("=", 26))
	fmt.Printf("%s\n\n", splash)

	fs := map[string]filesystem.Node{
		"bin": filesystem.Directory{Name: "bin"},
		"home": filesystem.Directory{Name: "home", Files: map[string]filesystem.Node{
			"fox": filesystem.Directory{Name: "fox", Files: map[string]filesystem.Node{
				"readme.txt": filesystem.File{Name: "readme.txt", Contents: []byte("this is a test file")},
				"work":       filesystem.Directory{Name: "work"},
			}},
		}},
	}

	prog := tea.NewProgram(screens.MainScreenModel("fox", "fos", fs))
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
}
