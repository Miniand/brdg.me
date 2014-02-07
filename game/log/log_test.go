package log

import (
	"testing"
)

func TestSortedMessages(t *testing.T) {
	l := Log{}
	m1 := NewPublicMessage("Test 1")
	m2 := NewPublicMessage("Test 2")
	m2.Time -= 100000
	l.Add(m1)
	l.Add(m2)
	messages := l.SortedMessages()
	if messages[0].Text != m2.Text {
		t.Fatal("Messages aren't sorted")
	}
}

func TestPublicMessages(t *testing.T) {
	l := Log{}
	m1 := NewPublicMessage("Test 1")
	m2 := NewPrivateMessage("Test 2", []string{"Bob"})
	l.Add(m1)
	l.Add(m2)
	messages := l.PublicMessages()
	if len(messages) != 1 {
		t.Fatal("Did not get 1 public message")
	}
	if messages[0].Text != m1.Text {
		t.Fatal("Message is not public")
	}
}

func TestMessagesFor(t *testing.T) {
	l := Log{}
	m1 := NewPublicMessage("Test 1")
	m2 := NewPrivateMessage("Test 2", []string{"Bob"})
	l.Add(m1)
	l.Add(m2)
	messages := l.MessagesFor("Bob")
	if len(messages) != 2 {
		t.Fatal("Did not get 2 messages for Bob")
	}
}

func TestNewMessagesFor(t *testing.T) {
	l := Log{}
	m1 := NewPublicMessage("Test 1")
	m2 := NewPrivateMessage("Test 2", []string{"Bob"})
	m2.Time += 1
	l.Add(m1)
	l.Add(m2)
	messages := l.NewMessagesFor("Bob")
	if len(messages) != 2 {
		t.Fatal("Did not get all messages when not yet read")
	}
	l.LastReadTimeFor = map[string]int64{}
	l.LastReadTimeFor["Bob"] = m1.Time
	messages = l.NewMessagesFor("Bob")
	if len(messages) != 1 {
		t.Fatal("Did not get 1 message for Bob")
	}
	if messages[0].Text != m2.Text {
		t.Fatal("Unread message wasn't m2")
	}
	l.MarkReadFor("Bob")
	messages = l.NewMessagesFor("Bob")
	if len(messages) != 0 {
		t.Fatal("Expected to get no more unread messages")
	}
}
