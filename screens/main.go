package screens

import (
	"fauxos/commands"
	"fauxos/filesystem"
	"fauxos/styles"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type MainScreen struct {
	quitting bool
	err      error

	fs  map[string]filesystem.Node
	cwd string

	input  textinput.Model
	output string

	history       []string
	historyCursor int
}

func MainScreenModel() MainScreen {
	cwd := "/home/fox/"
	input := textinput.New()
	input.Focus()
	input.Prompt = getPrompt(cwd)
	input.TextStyle = styles.Ansi[7]

	return MainScreen{
		err: nil,

		fs: map[string]filesystem.Node{
			"bin": filesystem.Directory{Name: "bin"},
			"home": filesystem.Directory{Name: "home", Files: map[string]filesystem.Node{
				"fox": filesystem.Directory{Name: "fox", Files: map[string]filesystem.Node{
					"readme.txt": filesystem.File{Name: "readme.txt", Contents: []byte("this is a test file")},
					"work":       filesystem.Directory{Name: "work"},
				}},
			}},
		},
		cwd: cwd,

		input:  input,
		output: "",

		history:       make([]string, 0),
		historyCursor: -1,
	}
}

func (m MainScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m.handleCommand(m.input.Value())

		case tea.KeyUp:
			if m.historyCursor > 0 && m.historyCursor <= len(m.history) {
				m.historyCursor--
				m.input.SetValue(m.history[m.historyCursor])
				m.input.CursorEnd()
			}

		case tea.KeyDown:
			if m.historyCursor < len(m.history) {
				m.historyCursor++
				if m.historyCursor == len(m.history) {
					m.input.Reset()
				} else {
					m.input.SetValue(m.history[m.historyCursor])
					m.input.CursorEnd()
				}
			}

		case tea.KeyRunes, tea.KeyBackspace, tea.KeyDelete:
			if m.historyCursor < len(m.history) {
				m.historyCursor++
			}
		}
	}

	m.input.Prompt = getPrompt(m.cwd)

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m MainScreen) View() string {
	finalOutput := m.output
	if !m.quitting {
		finalOutput += m.input.View()
	}

	splash := styles.Ansi[8].Render("*********************")
	splash += "\n" + styles.Ansi[11].Copy().Bold(true).Render("FoxOS") + " v0.0.1-alpha"
	splash += "\nby " + styles.Ansi[4].Copy().Bold(true).Render("zp4rker")
	splash += "\n" + styles.Ansi[8].Render("*********************")
	return fmt.Sprintf("%s\n%s", strings.TrimSpace(splash), finalOutput)
}

func (m MainScreen) handleCommand(command string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.input.Cursor.Blink = true
	m.output += m.input.View() + "\n"

	switch command {
	case "quit", "exit", "logout":
		m.output += "Quitting FoxOS...\n"
		m.quitting = true
		cmd = tea.Quit
	case "lf", "list-files":
		m.output += commands.ListFiles(strings.Trim(m.cwd, "/"), m.fs)

	default:
		m.output += "Could not find command '" + command + "'\n"
	}

	m.history = append(m.history, command)
	m.historyCursor = len(m.history)
	m.input.Reset()

	return m, cmd
}

func getPrompt(path string) string {
	return styles.Prompt1.Render("user@host ") + styles.Prompt2.Render(path) + " "
}
