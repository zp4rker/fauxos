package screens

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginScreen struct {
	complete bool

	Success bool
	User    string

	userInput textinput.Model
	passInput textinput.Model

	options map[string]string
}

func LoginScreenModel(logins map[string]string) LoginScreen {
	userInput := textinput.New()
	userInput.Prompt = "Username: "
	userInput.Focus()

	passInput := textinput.New()
	passInput.Prompt = "Password: "
	passInput.EchoMode = textinput.EchoPassword

	return LoginScreen{

		userInput: userInput,
		passInput: passInput,

		options: logins,
	}
}

func (m *LoginScreen) View() string {
	var view string

	if m.complete {
		return view
	}

	view += m.userInput.View() + "\n"
	view += m.passInput.View()

	return view
}

func (m *LoginScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.complete {
		return m, tea.Quit
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if m.userInput.Value() != "" && m.passInput.Value() != "" {
				if pass, ok := m.options[m.userInput.Value()]; ok {
					if pass == m.passInput.Value() {
						m.User = m.userInput.Value()
						m.Success = true
					}
				}
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

	if m.userInput.Focused() {
		m.userInput, cmd = m.userInput.Update(msg)
	} else if m.passInput.Focused() {
		m.passInput, cmd = m.passInput.Update(msg)
	}
	return m, cmd
}

func (m *LoginScreen) Init() tea.Cmd {
	return textinput.Blink
}
