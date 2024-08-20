package screens

import (
	"fauxos/filesystem"
	"fauxos/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Login struct {
	complete bool
	errorMsg string

	success bool
	user    string

	userInput textinput.Model
	passInput textinput.Model

	options map[string]string
	debug   bool
}

func (m Login) View() string {
	var view string

	if m.complete {
		if m.errorMsg != "" {
			view = m.errorMsg
		}
		return view
	}

	view += "Log into "
	view += styles.Ansi[11].Copy().Bold(true).Render("fOS") + " "
	view += styles.Ansi[7].Render("v0.0.1") + "\n\n"

	view += m.userInput.View() + "\n"
	view += m.passInput.View()

	return view
}

func (m Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.complete {
		if m.success {
			main.currentModel = modelShell
			return m, m.openShell()
		} else {
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if m.userInput.Value() != "" && m.passInput.Value() != "" {
				if pass, ok := m.options[m.userInput.Value()]; ok {
					if pass == m.passInput.Value() {
						m.user = m.userInput.Value()
						m.success = true
					} else {
						m.errorMsg = "Password incorrect!\n"
					}
				} else {
					m.errorMsg = "User does not exist!\n"
				}
				m.complete = true
				return m, nil
			} else if m.debug {
				m.user = "fox"
				if m.userInput.Value() != "" {
					m.user = m.userInput.Value()
				}
				m.success = true
				m.complete = true
				return m, nil
			} else if m.userInput.Value() != "" && m.passInput.Value() == "" && m.userInput.Focused() {
				m.passInput.Focus()
				m.userInput.Blur()
			} else if m.userInput.Value() == "" && m.passInput.Focused() {
				m.userInput.Focus()
				m.passInput.Blur()
			}
		}

		if msg.Type == tea.KeyTab || msg.Type == tea.KeyShiftTab {
			if m.userInput.Focused() {
				m.userInput.Blur()
				m.passInput.Focus()
			} else {
				m.passInput.Blur()
				m.userInput.Focus()
			}
		}
	}

	var cmd tea.Cmd
	if m.userInput.Focused() {
		m.userInput, cmd = m.userInput.Update(msg)
	} else {
		m.passInput, cmd = m.passInput.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func LoginModel(logins map[string]string, debug bool) Login {
	userInput := textinput.New()
	userInput.Prompt = "Username: "
	userInput.Focus()
	userInput.TextStyle = styles.Ansi[7]
	userInput.PromptStyle = styles.Bold

	passInput := textinput.New()
	passInput.Prompt = "Password: "
	passInput.EchoMode = textinput.EchoPassword
	passInput.TextStyle = styles.Ansi[7]
	passInput.PromptStyle = styles.Bold

	return Login{
		userInput: userInput,
		passInput: passInput,

		options: logins,
		debug:   debug,
	}
}

func (m Login) Init() tea.Cmd {
	return nil
}

func (m Login) openShell() tea.Cmd {
	var cmd tea.Cmd
	if !m.success {
		return cmd
	}

	splash := styles.Ansi[8].Render(strings.Repeat("=", 26))
	splash += "\n    Running " + styles.Ansi[11].Copy().Bold(true).Render("fOS") + " v0.0.1"
	splash += "\n        by " + styles.Ansi[4].Copy().Bold(true).Render("zp4rker")
	splash += "\n" + styles.Ansi[8].Render(strings.Repeat("=", 26)) + "\n"
	cmd = tea.Println(splash)

	fs := map[string]filesystem.Node{
		"bin": filesystem.Directory{Name: "bin"},
		"home": filesystem.Directory{Name: "home", Files: map[string]filesystem.Node{
			m.user: filesystem.Directory{Name: m.user, Files: map[string]filesystem.Node{
				"readme.txt": filesystem.File{Name: "readme.txt", Data: []byte("this is a test file")},
				"work":       filesystem.Directory{Name: "work"},
			}},
		}},
	}

	main.SetShell(ShellModel(m.user, "fos", fs))

	return cmd
}
