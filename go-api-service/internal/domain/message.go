package domain

import (
	"errors"
)

type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (m *Message) Validate() error {
	if m.Author == "" {
		return errors.New("author is required")
	}
	if m.Body == "" {
		return errors.New("body is required")
	}
	return nil
}
