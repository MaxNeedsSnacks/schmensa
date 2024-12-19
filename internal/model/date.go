package model

import (
	"cmp"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"maps"
	"schmensa/internal/data"
	"schmensa/internal/utils"
	"slices"
	"strings"
	"time"
)

var (
	dateSpinner = gloss.NewStyle().
			Foreground(gloss.Color("#D96546"))

	dateSubmitKey   = utils.SubmitKey()
	datePreviousKey = utils.PreviousMenuKey()
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
	mensa data.Mensa

	list    list.Model
	spinner spinner.Model
	loading bool
}

func SelectDate(mensa data.Mensa) tea.Model {
	model := selectDateModel{
		mensa:   mensa,
		list:    list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(dateSpinner)),
		loading: true,
	}

	model.list.Title = fmt.Sprintf("Viewing dates for %s", mensa.Name.De)
	model.list.Styles = dateListStyle()

	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{dateSubmitKey, datePreviousKey}
	}
	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{dateSubmitKey, datePreviousKey}
	}

	return model
}

func allMenuUrl(id string) string {
	return "https://mobil.itmc.tu-dortmund.de/canteen-menu/v3/canteens/" + id
}

func (m selectDateModel) loadDateMap() tea.Msg {
	menuMap, err := utils.FromRemoteJson[map[string]data.Menu](allMenuUrl(m.mensa.Id))

	if err != nil {
		return data.WrapError(err)
	}

	dateMap := make(map[date]data.Menu, len(*menuMap))
	for str, value := range *menuMap {
		parsed, err := time.Parse(time.DateOnly, str)

		if err != nil {
			return data.WrapError(err)
		}

		dateMap[date(parsed)] = value
	}

	return dateMapMsg(dateMap)
}

type date time.Time

func (d date) FilterValue() string {
	return strings.Join([]string{
		d.Title(),
		d.Description(),
		time.Time(d).Format("Monday January"),
	}, " ")
}

func (d date) Title() string {
	return time.Time(d).Format(time.DateOnly)
}

func (d date) Description() string {
	return utils.FormatRelativeToday(time.Time(d))
}

type dateMapMsg map[date]data.Menu

func (m selectDateModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadDateMap)
}

func (m selectDateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := utils.AppStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case dateMapMsg:
		dates := slices.SortedFunc(maps.Keys(msg), func(a, b date) int {
			first, second := time.Time(a), time.Time(b)
			return int(first.Unix() - second.Unix())
		})

		var items = make([]list.Item, 0, len(msg))
		for _, date := range dates {
			items = append(items, date)
		}

		closestDate := slices.MinFunc(dates, func(a, b date) int {
			currentTime := time.Now()
			diffA := utils.Abs(time.Time(a).Sub(currentTime))
			diffB := utils.Abs(time.Time(b).Sub(currentTime))
			return cmp.Compare(diffA, diffB)
		})

		m.list.Select(slices.Index(dates, closestDate))
		m.loading = false
		return m, m.list.SetItems(items)

	case data.ErrMsg:
		m.loading = false
		return m, tea.Quit

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, dateSubmitKey):
			return m, func() tea.Msg {
				if v, ok := m.list.SelectedItem().(date); ok {
					return data.NewError(fmt.Sprintf("date selection not implemented! %s", v.Title()))
				} else {
					return data.NewError("selected item was not a date! (somehow...?)")
				}
			}
		case key.Matches(msg, datePreviousKey):
			return utils.ChangeModel(SelectMensa())
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
		return fmt.Sprintf("%s Loading menu for %s...", m.spinner.View(), m.mensa.Name.De)
	}
	return m.list.View()
}
