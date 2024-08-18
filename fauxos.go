package main

import (
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

	prog := tea.NewProgram(screens.MainScreenModel())
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
}
