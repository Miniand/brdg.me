package log

import (
	"time"
)

type Message struct {
	Text    string
	Private bool
	To      []string
	Time    int64
}

// Create a new message timed at now
func NewMessage() Message {
	return Message{
		Time: time.Now().UnixNano(),
	}
}

func NewPublicMessage(text string) Message {
	m := NewMessage()
	m.Text = text
	return m
}

func NewPrivateMessage(text string, to []string) Message {
	m := NewMessage()
	m.Text = text
	m.Private = true
	m.To = to
	return m
}

func (m Message) CanRead(player string) bool {
	if !m.Private {
		return true
	}
	for _, to := range m.To {
		if to == player {
			return true
		}
	}
	return false
}
