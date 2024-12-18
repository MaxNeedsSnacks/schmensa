package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	submitKey = key.NewBinding(
		key.WithKeys("enter", "ctrl+q"),
		key.WithHelp("enter", "select item"),
	)
)

type SelectMensaModel struct {
	list list.Model
}

func SelectMensa() tea.Model {
	model := SelectMensaModel{
		list: list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
	}
	model.list.Title = "Select a Mensa"
	model.list.Styles.Title = titleStyle

	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{submitKey}
	}
	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{submitKey}
	}

	return model
}

const url = "https://mobil.itmc.tu-dortmund.de/canteen-menu/v3/canteens/"

func loadMensaList() tea.Msg {
	mensas, err := fromJsonUrl[[]mensa](url)

	if err != nil {
		return errMsg{err}
	}

	return mensaMsg(mensas)
}

func (m mensa) FilterValue() string {
	return m.Name.De
}
func (m mensa) Title() string {
	return m.Name.De
}
func (m mensa) Description() string {
	return fmt.Sprintf("(ID %s)", m.Id)
}

type mensaMsg *[]mensa
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m SelectMensaModel) Init() tea.Cmd {
	return tea.Sequence(m.list.StartSpinner(), loadMensaList)
}

func (m SelectMensaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case mensaMsg:
		var items = make([]list.Item, 0, len(*msg))
		for _, item := range *msg {
			items = append(items, item)
		}

		m.list.StopSpinner()
		return m, m.list.SetItems(items)
	case errMsg:
		return nil, tea.Quit
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		if key.Matches(msg, submitKey) {
			return nil, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectMensaModel) View() string {
	return appStyle.Render(m.list.View())
}
