package main

import (
	"fauxos/screens"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	login := screens.LoginModel(map[string]string{"zp4rker": "bismillah"})
	prog := tea.NewProgram(screens.MainModel(login))
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
}
