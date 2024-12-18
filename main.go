package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"os"
	"time"
)

var (
	appStyle = gloss.NewStyle().Padding(1, 2)

	errorStyle = gloss.NewStyle().
			Foreground(gloss.Color("15")).
			Background(gloss.Color("9")).
			Padding(0, 1)

	errorMessageStyle = gloss.NewStyle().
				Foreground(gloss.Color("9"))
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type model struct {
	inner tea.Model
	error error
}

func (m model) Init() tea.Cmd {
	return m.inner.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.error = msg.err
		return m, DelayCmd(3*time.Second, tea.Quit)
	default:
		m.inner, cmd = m.inner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.error != nil {
		return appStyle.Render(
			errorStyle.Render("ERROR"),
			errorMessageStyle.Render(m.error.Error()),
		)
	}

	return appStyle.Render(m.inner.View())
}

func main() {
	p := tea.NewProgram(model{inner: SelectMensa()})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
