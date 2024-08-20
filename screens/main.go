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

	user    string
	machine string

	fs  map[string]filesystem.Node
	cwd string

	input textinput.Model

	history       []string
	historyCursor int
}

func MainScreenModel() MainScreen {
	user := "fox"
	cwd := fmt.Sprintf("/home/%s/", user)
	input := textinput.New()
	input.Focus()
	input.TextStyle = styles.Ansi[7]

	return MainScreen{
		err: nil,

		user:    user,
		machine: "fos",

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

		input: input,

		history:       make([]string, 0),
		historyCursor: -1,
	}
}

func (m MainScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.input.Prompt = getPrompt(m)

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			parts := strings.Split(m.input.Value(), " ")
			return m.handleCommand(parts[0], parts[1:])
		}

		if msg.Type == tea.KeyUp {
			if m.historyCursor > 0 && m.historyCursor <= len(m.history) {
				m.historyCursor--
				m.input.SetValue(m.history[m.historyCursor])
				m.input.CursorEnd()
			}
		}

		if msg.Type == tea.KeyDown {
			if m.historyCursor < len(m.history) {
				m.historyCursor++
				if m.historyCursor == len(m.history) {
					m.input.Reset()
				} else {
					m.input.SetValue(m.history[m.historyCursor])
					m.input.CursorEnd()
				}
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m MainScreen) View() string {
	return m.input.View()
}

func (m MainScreen) handleCommand(command string, args []string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.input.Cursor.Blink = true
	output(m.input.View() + "\n")

	switch command {
	case "":
		// do nothing

	case "quit", "exit", "logout":
		output("Quitting...\n")
		m.quitting = true
		cmd = tea.Quit

	case "lf", "list-files":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		output(commands.ListFiles(path, m.fs, m.cwd))

	case "cd", "change-directory":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		cwd, output := commands.ChangeDirectory(path, m.fs, m.cwd)
		if output != "" {
			fmt.Print(output)
		} else {
			if cwd == "" {
				cwd = "/"
			}
			m.cwd = cwd
			m.input.Prompt = getPrompt(m)
		}

	case "pf", "print-file":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		output(commands.PrintFile(path, m.fs, m.cwd))

	case "ad", "add-dir", "add-directory":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		out, fs := commands.AddDirectory(path, &m.fs, m.cwd)
		m.fs = fs
		output(out)

	case "rd", "remove-dir", "remove-directory":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		out, fs := commands.RemoveDirectory(path, &m.fs, m.cwd)
		m.fs = fs
		output(out)

	default:
		output(fmt.Sprintf("Could not find command '%s'\n", command))
	}

	m.history = append(m.history, command)
	m.historyCursor = len(m.history)
	m.input.Reset()

	return m, cmd
}

func getPrompt(m MainScreen) string {
	path := strings.ReplaceAll(m.cwd, "/home/"+m.user, "~")
	return styles.Prompt1.Render(fmt.Sprintf("%s@%s ", m.user, m.machine)) + styles.Prompt2.Render(path) + " "
}

func output(out string) {
	out = strings.ReplaceAll(out, "\n", "\n\r")
	out = "\r" + out
	fmt.Print(out)
}
