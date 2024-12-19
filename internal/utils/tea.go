package utils

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func DelayCmd(d time.Duration, cb tea.Cmd) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return cb()
	})
}

func ChangeModel(m tea.Model) (tea.Model, tea.Cmd) {
	return m, tea.Batch(m.Init(), tea.WindowSize())
}

func SubmitKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("enter", "ctrl+q", " "),
		key.WithHelp("enter/␣", "select item"),
	)
}

func PreviousMenuKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("backspace", "shift+backspace"),
		key.WithHelp("⌫", "back to previous menu"),
	)
}
