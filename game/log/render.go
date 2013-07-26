package log

import (
	"bytes"
	"math"
	"time"
)

func RenderMessages(messages []Message) string {
	buf := bytes.NewBuffer([]byte{})
	firstPass := true
	for _, m := range messages {
		if !firstPass {
			buf.WriteByte('\n')
		}
		firstPass = false
		buf.WriteString(RenderMessage(m))
	}
	return buf.String()
}

func RenderMessage(message Message) string {
	t := time.Unix(message.Time/int64(math.Pow10(9)), 0)
	return t.UTC().Format(time.RFC822) + ": " + message.Text
}
