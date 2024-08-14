package screens

import (
	"fauxos/styles"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainScreen struct {
	quitting      bool
	input         textinput.Model
	output        string
	history       []string
	historyCursor int
	err           error
}

func MainScreenModel() MainScreen {
	input := textinput.New()
	input.Focus()
	input.Prompt = styles.Prompt1.Render("user@host ") + styles.Prompt2.Render("~/some/path ")

	return MainScreen{
		input:         input,
		output:        "",
		history:       make([]string, 0),
		historyCursor: -1,
		err:           nil,
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

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m MainScreen) View() string {
	finalOutput := m.output
	if !m.quitting {
		finalOutput += m.input.View()
	}

	return fmt.Sprintf("Welcome to FoxOS!\n%s", finalOutput)
}

func (m MainScreen) handleCommand(command string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch command {
	case "quit", "exit", "logout":
		m.input.Cursor.Blink = true
		m.output += m.input.View() + "\n"
		m.output += "Quitting FoxOS...\n"
		m.quitting = true
		cmd = tea.Quit
	default:
		m.input.Cursor.Blink = true
		m.output += m.input.View() + "\n"
		m.output += "Could not find command '" + command + "'!\n"
	}

	m.history = append(m.history, command)
	m.historyCursor = len(m.history)
	m.input.Reset()

	return m, cmd
}
