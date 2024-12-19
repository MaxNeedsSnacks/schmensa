package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"schmensa/internal/data"
	"schmensa/internal/model/mensa"
	"schmensa/internal/utils"
	"time"
)

var (
	appStyle          = utils.AppStyle()
	errorStyle        = utils.ErrorStyle()
	errorMessageStyle = utils.ErrorMessageStyle()
)

type mainModel struct {
	inner tea.Model
	error error
}

func (m mainModel) Init() tea.Cmd {
	return m.inner.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case data.ErrMsg:
		m.error = msg
		return m, utils.DelayCmd(3*time.Second, tea.Quit)
	default:
		m.inner, cmd = m.inner.Update(msg)
		return m, cmd
	}
}

func (m mainModel) View() string {
	if m.error != nil {
		return appStyle.Render(
			errorStyle.Render("ERROR"),
			errorMessageStyle.Render(m.error.Error()),
		)
	}

	return appStyle.Render(m.inner.View())
}

func main() {
	p := tea.NewProgram(mainModel{inner: mensa.SelectMensa()})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
