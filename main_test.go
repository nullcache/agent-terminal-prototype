package main

import (
	"strings"
	"testing"
	"testing/iotest"
)

func TestCleanTerminal(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"CSI and Unicode", "\x1b[32m开始\x1b[0m\x1b[2J\n", "开始\n"},
		{"OSC with BEL or ST", "\x1b]0;title\a\x1b]8;;https://example.org\x1b\\链接\x1b]8;;\x1b\\", "链接"},
		{"progress preserves logs", "start\r\n10%\r50%\r\x1b[32m100%\x1b[0m\r\ndone\n", "start\n100%\ndone\n"},
		{"trailing carriage return", "last\r", "last"},
		{"incomplete escape", "last\x1b[3", "last"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got strings.Builder
			input := iotest.OneByteReader(strings.NewReader(tc.input))
			if err := cleanTerminal(input, &got); err != nil {
				t.Fatal(err)
			}
			if got.String() != tc.want {
				t.Fatalf("got %q, want %q", got.String(), tc.want)
			}
		})
	}
}
