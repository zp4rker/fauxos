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
	login := screens.LoginScreenModel(map[string]string{"zp4rker": "bismillah"})
	prog := tea.NewProgram(&login)
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}

	openShell(login.Success, login.User)
}

func openShell(success bool, user string) {
	if !success {
		return
	}

	splash := styles.Ansi[8].Render(strings.Repeat("=", 26))
	splash += "\n    Running " + styles.Ansi[11].Copy().Bold(true).Render("fOS") + " v0.0.1"
	splash += "\n        by " + styles.Ansi[4].Copy().Bold(true).Render("zp4rker")
	splash += "\n" + styles.Ansi[8].Render(strings.Repeat("=", 26))
	fmt.Printf("%s\n\n", splash)

	fs := map[string]filesystem.Node{
		"bin": filesystem.Directory{Name: "bin"},
		"home": filesystem.Directory{Name: "home", Files: map[string]filesystem.Node{
			user: filesystem.Directory{Name: user, Files: map[string]filesystem.Node{
				"readme.txt": filesystem.File{Name: "readme.txt", Data: []byte("this is a test file")},
				"work":       filesystem.Directory{Name: "work"},
			}},
		}},
	}

	prog := tea.NewProgram(screens.MainScreenModel(user, "fos", fs))
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
}
