package mensa

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"schmensa/internal/data"
	"schmensa/internal/model/date"
	"schmensa/internal/utils"
)

var (
	mensaSpinner = gloss.NewStyle().
			Foreground(gloss.Color("#25A065"))

	submitKey = utils.SubmitKey()
)

func mensaListStyle() list.Styles {
	s := list.DefaultStyles()

	s.Title = gloss.NewStyle().
		Foreground(gloss.Color("#FFFDF5")).
		Background(gloss.Color("#25A065")).
		Padding(0, 1)

	return s
}

type selectMensaModel struct {
	list    list.Model
	spinner spinner.Model
	loading bool
}

func SelectMensa() tea.Model {
	model := selectMensaModel{
		list:    list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(mensaSpinner)),
		loading: true,
	}

	model.list.Title = "Select a Mensa"
	model.list.Styles = mensaListStyle()

	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{submitKey}
	}
	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{submitKey}
	}

	return model
}

const canteensUrl = "https://mobil.itmc.tu-dortmund.de/canteen-menu/v3/canteens/"

func loadMensaList() tea.Msg {
	mensas, err := utils.FromRemoteJson[[]mensa](canteensUrl)

	if err != nil {
		return data.WrapError(err)
	}

	return mensaMsg(*mensas)
}

type mensa data.Mensa

func (m mensa) FilterValue() string { return m.Name.De }
func (m mensa) Title() string       { return m.Name.De }
func (m mensa) Description() string { return fmt.Sprintf("(ID %s)", m.Id) }

type mensaMsg []mensa
type mensaSelectMsg mensa

func (m selectMensaModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, loadMensaList)
}

func (m selectMensaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := utils.AppStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case mensaMsg:
		var items = make([]list.Item, 0, len(msg))
		for _, item := range msg {
			items = append(items, item)
		}

		m.loading = false
		return m, m.list.SetItems(items)

	case data.ErrMsg:
		m.loading = false
		return m, tea.Quit

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		if key.Matches(msg, submitKey) {
			return m, func() tea.Msg {
				if v, ok := m.list.SelectedItem().(mensa); ok {
					return mensaSelectMsg(v)
				} else {
					return data.NewError("selected item was not a mensa! (somehow...?)")
				}
			}
		}

	case mensaSelectMsg:
		return utils.ChangeModel(date.SelectDate(data.Mensa(msg)))
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m selectMensaModel) View() string {
	if m.loading {
		return fmt.Sprintf("%s Loading Mensa List...", m.spinner.View())
	}
	return m.list.View()
}
