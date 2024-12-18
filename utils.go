package main

import (
	"encoding/json"
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
