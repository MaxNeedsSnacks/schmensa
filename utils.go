package main

import (
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"time"
)

func fromJsonUrl[T any](url string) (*T, error) {
	c := &http.Client{Timeout: 10 * time.Second}

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
