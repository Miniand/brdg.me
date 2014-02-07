package log

import (
	"sort"
	"time"
)

type Log struct {
	sort.Interface
	Messages        []Message
	LastReadTimeFor map[string]int64
}

func New() *Log {
	return &Log{
		Messages:        []Message{},
		LastReadTimeFor: map[string]int64{},
	}
}

func (l *Log) SortedMessages() []Message {
	sort.Sort(l)
	return l.Messages
}

func (l *Log) PublicMessages() []Message {
	messages := []Message{}
	for _, m := range l.SortedMessages() {
		if !m.Private {
			messages = append(messages, m)
		}
	}
	return messages
}

func (l *Log) MessagesFor(player string) []Message {
	messages := []Message{}
	for _, m := range l.SortedMessages() {
		if m.CanRead(player) {
			messages = append(messages, m)
		}
	}
	return messages
}

func (l *Log) NewMessagesFor(player string) []Message {
	messages := []Message{}
	if l.LastReadTimeFor == nil {
		l.LastReadTimeFor = map[string]int64{}
	}
	for _, m := range l.MessagesFor(player) {
		if l.LastReadTimeFor[player] == 0 || l.LastReadTimeFor[player] < m.Time {
			messages = append(messages, m)
		}
	}
	return messages
}

func (l *Log) MarkReadFor(player string) {
	if l.LastReadTimeFor == nil {
		l.LastReadTimeFor = map[string]int64{}
	}
	l.LastReadTimeFor[player] = time.Now().UnixNano()
}

func (l *Log) Len() int {
	return len(l.Messages)
}

func (l *Log) Less(i, j int) bool {
	return l.Messages[i].Time < l.Messages[j].Time
}

func (l *Log) Swap(i, j int) {
	l.Messages[i], l.Messages[j] = l.Messages[j], l.Messages[i]
}

func (l *Log) Add(message Message) {
	l.Messages = append(l.Messages, message)
}
