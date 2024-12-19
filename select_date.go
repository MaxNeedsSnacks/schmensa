package main

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"maps"
	"slices"
	"time"
)

var (
	dateSpinner = gloss.NewStyle().
		Foreground(gloss.Color("#D96546"))
)

func dateListStyle() list.Styles {
	s := list.DefaultStyles()

	s.Title = gloss.NewStyle().
		Foreground(gloss.Color("#FFFDF5")).
		Background(gloss.Color("#D96546")).
		Padding(0, 1)

	return s
}

type selectDateModel struct {
	mensa mensa

	list    list.Model
	spinner spinner.Model
	loading bool
}

func SelectDate(mensa mensa) tea.Model {
	model := selectDateModel{
		mensa:   mensa,
		list:    list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(dateSpinner)),
		loading: true,
	}

	model.list.Title = fmt.Sprintf("Viewing dates for %s", mensa.Title())
	model.list.Styles = dateListStyle()

	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{submitKey, previousMenuKey}
	}
	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{submitKey, previousMenuKey}
	}

	return model
}

func allMenuUrl(id string) string {
	return "https://mobil.itmc.tu-dortmund.de/canteen-menu/v3/canteens/" + id
}

func (m selectDateModel) loadDateMap() tea.Msg {
	menuMap, err := FromRemoteJson[map[string]menu](allMenuUrl(m.mensa.Id))

	if err != nil {
		return errMsg{err}
	}

	dateMap := make(map[date]menu, len(*menuMap))
	for str, value := range *menuMap {
		parsed, err := time.Parse(time.DateOnly, str)

		if err != nil {
			return errMsg{err}
		}

		dateMap[date(parsed)] = value
	}

	return dateMapMsg(dateMap)
}

type date time.Time

func (d date) FilterValue() string {
	return time.Time(d).Format("Monday January") + " " + d.Title()
}
func (d date) Title() string {
	return time.Time(d).Format(time.DateOnly)
}
func (d date) Description() string {
	return FormatRelativeToday(time.Time(d))
}

type dateMapMsg map[date]menu

func (m selectDateModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadDateMap)
}

func (m selectDateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case dateMapMsg:
		dates := slices.SortedFunc(maps.Keys(msg), func(date date, other date) int {
			return time.Time(date).Day() - time.Time(date).Day()
		})

		var items = make([]list.Item, 0, len(msg))
		for _, date := range dates {
			items = append(items, date)
		}

		m.loading = false
		return m, m.list.SetItems(items)

	case errMsg:
		m.loading = false
		return m, tea.Quit

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, submitKey):
			return m, func() tea.Msg {
				if v, ok := m.list.SelectedItem().(date); ok {
					return errMsg{fmt.Errorf("date selection not implemented! %s", v.Title())}
				} else {
					return errMsg{errors.New("selected item was not a date! (somehow...?)")}
				}
			}
		case key.Matches(msg, previousMenuKey):
			return ChangeModel(SelectMensa())
		}
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m selectDateModel) View() string {
	if m.loading {
		return fmt.Sprintf("%s Loading menu for %s...", m.spinner.View(), m.mensa.Title())
	}
	return m.list.View()
}
