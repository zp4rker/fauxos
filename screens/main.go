package screens

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState int

const (
	modelLogin modelState = iota
	modelShell
)

var main *Main

type Main struct {
	currentModel modelState

	loginModel, shellModel tea.Model
}

func (m *Main) View() string {
	if m.currentModel == modelLogin {
		return m.loginModel.View()
	} else {
		return m.shellModel.View()
	}
}

func (m *Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.currentModel == modelLogin {
		model, cmd := m.loginModel.Update(msg)
		m.loginModel = model
		return m, cmd
	} else {
		model, cmd := m.shellModel.Update(msg)
		m.shellModel = model
		return m, cmd
	}
}

func (m *Main) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Main) SetShell(shellModel Shell) {
	shellModel.input.Prompt = getPrompt(shellModel)
	m.shellModel = shellModel
}

func MainModel(loginModel Login) tea.Model {
	main = &Main{
		currentModel: modelLogin,

		loginModel: loginModel,
	}

	return main
}
