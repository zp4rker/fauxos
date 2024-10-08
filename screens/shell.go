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

type Shell struct {
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

func (m Shell) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	return m.input.View()
}

func (m Shell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Shell) handleCommand(command string, args []string) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var out string

	switch command {
	case "":
		// do nothing

	case "quit", "exit", "logout":
		m.quitting = true

	case "lf", "list-files":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		out = commands.ListFiles(path, m.fs, m.cwd)

	case "cd", "change-directory":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		var cwd string
		cwd, out = commands.ChangeDirectory(path, m.fs, m.cwd)
		if out == "" {
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
		out = commands.PrintFile(path, m.fs, m.cwd)

	case "ad", "add-dir", "add-directory":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		var fs map[string]filesystem.Node
		out, fs = commands.AddDirectory(path, &m.fs, m.cwd)
		m.fs = fs

	case "rd", "remove-dir", "remove-directory":
		path := m.cwd
		if len(args) > 0 {
			path = args[0]
		}
		var fs map[string]filesystem.Node
		out, fs = commands.RemoveDirectory(path, &m.fs, m.cwd)
		m.fs = fs

	default:
		cmds = append(cmds, tea.Printf("Could not find command '%s'", command))
	}

	if out != "" {
		cmds = append(cmds, tea.Println(out))
	}
	m.input.Cursor.Blink = true
	cmds = append(cmds, tea.Println(m.input.View()))

	m.history = append(m.history, command)
	m.historyCursor = len(m.history)
	m.input.Reset()

	return m, tea.Batch(cmds...)
}

func (m Shell) Init() tea.Cmd {
	return nil
}

func ShellModel(user, machine string) Shell {
	cwd := fmt.Sprintf("/home/%s/", user)
	input := textinput.New()
	input.Focus()
	input.TextStyle = styles.Ansi[7]

	return Shell{
		err: nil,

		user:    user,
		machine: machine,

		fs:  defaultFS(user),
		cwd: cwd,

		input: input,

		history:       make([]string, 0),
		historyCursor: -1,
	}
}

func getPrompt(m Shell) string {
	path := strings.ReplaceAll(m.cwd, "/home/"+m.user, "~")
	return styles.Prompt1.Render(fmt.Sprintf("%s@%s ", m.user, m.machine)) + styles.Prompt2.Render(path) + " "
}

func splash() tea.Cmd {
	splash := styles.Ansi[8].Render(strings.Repeat("=", 26))
	splash += "\n    Running " + styles.Ansi[11].Copy().Bold(true).Render("fOS") + " v0.0.1"
	splash += "\n        by " + styles.Ansi[4].Copy().Bold(true).Render("zp4rker")
	splash += "\n" + styles.Ansi[8].Render(strings.Repeat("=", 26)) + "\n"
	return tea.Println(splash)
}

func defaultFS(user string) map[string]filesystem.Node {
	fs := map[string]filesystem.Node{}

	fs["home"] = filesystem.Directory{Name: "home", Files: map[string]filesystem.Node{
		user: filesystem.Directory{Name: user, Files: map[string]filesystem.Node{
			"readme.txt": filesystem.File{Name: "readme.txt", Data: []byte("this is a test file")},
			"work":       filesystem.Directory{Name: "work"},
		}},
	}}

	fs["bin"] = filesystem.Directory{Name: "bin", Files: map[string]filesystem.Node{
		"ad": filesystem.File{Name: "ad"},
		"rd": filesystem.File{Name: "rd"},
		"cd": filesystem.File{Name: "cd"},

		"lf": filesystem.File{Name: "lf"},

		"pf": filesystem.File{Name: "pf"},
	}}

	return fs
}
