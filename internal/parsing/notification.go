package parsing

import (
	"strings"
)

type Notification struct {
	Title string
	Body  string
}

func ParseNotification(input string) *Notification {
	var notif Notification

	input = strings.TrimSpace(input)
	s := newState(input)

	var ok bool
	s, ok = consumeIf(s, func(r rune) bool { return r == '[' })
	if ok {
		var (
			title string
			err   error
		)
		title, s, err = lexBalanced(s, '[', ']')
		notif.Title = strings.TrimSpace(title)
		if err == errEndOfInput {
			return &notif
		}
	}
	notif.Body = strings.TrimSpace(s.remaining())
	return &notif
}
