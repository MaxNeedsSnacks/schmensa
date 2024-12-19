package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"time"
)

var (
	submitKey = key.NewBinding(
		key.WithKeys("enter", "ctrl+q", " "),
		key.WithHelp("enter/␣", "select item"),
	)

	previousMenuKey = key.NewBinding(
		key.WithKeys("esc", "shift+esc"),
		key.WithHelp("esc", "back to previous menu"),
	)
)

func FromRemoteJson[T any](url string) (*T, error) {
	c := &http.Client{Timeout: 1 * time.Second}

	res, err := c.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var ret T
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func DelayCmd(d time.Duration, cb tea.Cmd) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return cb()
	})
}

func FormatRelativeToday(targetDate time.Time) string {
	nowDate := time.Now()

	daysApart := targetDate.Day() - nowDate.Day()
	nowYear, nowWeek := nowDate.ISOWeek()
	targetYear, targetWeek := targetDate.ISOWeek()

	var weeksApart int
	switch targetYear - nowYear {
	case -1:
		weeksApart = targetWeek - 52 - nowWeek
	case 0:
		weeksApart = targetWeek - nowWeek
	case 1:
		weeksApart = targetWeek + 52 - nowWeek
	default:
		panic("date results should not have been more than a year apart!")
	}

	weekday := targetDate.Weekday()

	switch daysApart {
	case -1:
		return "yesterday"
	case 0:
		return "today"
	case 1:
		return "tomorrow"
	default:
		if targetDate.Before(nowDate) {
			return fmt.Sprintf("%d days ago", daysApart)
		} else {
			switch weeksApart {
			case 0:
				return fmt.Sprintf("this %s", weekday)
			case 1:
				return fmt.Sprintf("next %s", weekday)
			default:
				return fmt.Sprintf("%s in %d weeks", weekday, weeksApart)
			}
		}
	}
}

func ChangeModel(m tea.Model) (tea.Model, tea.Cmd) {
	return m, m.Init()
}
