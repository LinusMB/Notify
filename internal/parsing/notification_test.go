package parsing

import "testing"

func TestParseNotification(t *testing.T) {
	tests := []struct {
		input string
		want  Notification
	}{
		{"[Title]Body", Notification{Title: "Title", Body: "Body"}},
		{
			"   [    Title      ]     Body     ",
			Notification{Title: "Title", Body: "Body"},
		},
		{"Body", Notification{Title: "", Body: "Body"}},
		{"[Title]", Notification{Title: "Title", Body: ""}},
		{"[Title][]Body", Notification{Title: "Title", Body: "[]Body"}},
		{"[Ti[]tle]Body", Notification{Title: "Ti[]tle", Body: "Body"}},
	}
	for _, tt := range tests {
		testname := tt.input
		t.Run(testname, func(t *testing.T) {
			got := ParseNotification(tt.input)
			if *got != tt.want {
				t.Errorf("got %v, want %v", *got, tt.want)
			}
		})
	}
}
