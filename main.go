package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type escapeState uint8

const (
	plain escapeState = iota
	escape
	csi
	osc
	oscEscape
)

func cleanTerminal(input io.Reader, output io.Writer) error {
	reader := bufio.NewReader(input)
	var line strings.Builder
	var state escapeState
	var rewrite bool

	for {
		char, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			_, err := io.WriteString(output, line.String())
			return err
		}
		if err != nil {
			return err
		}

		switch state {
		case escape:
			switch char {
			case '[':
				state = csi
			case ']':
				state = osc
			default:
				state = plain
			}
			continue
		case csi:
			if char >= 0x40 && char <= 0x7e {
				state = plain
			}
			continue
		case osc:
			switch char {
			case '\a':
				state = plain
			case '\x1b':
				state = oscEscape
			}
			continue
		case oscEscape:
			switch char {
			case '\\', '\a':
				state = plain
			case '\x1b':
			default:
				state = osc
			}
			continue
		}

		switch char {
		case '\x1b':
			state = escape
		case '\r':
			rewrite = true
		case '\n':
			line.WriteByte('\n')
			if _, err := io.WriteString(output, line.String()); err != nil {
				return err
			}
			line.Reset()
			rewrite = false
		default:
			if char != '\t' && unicode.IsControl(char) {
				continue
			}
			if rewrite {
				line.Reset()
				rewrite = false
			}
			line.WriteRune(char)
		}
	}
}

func main() {
	if err := cleanTerminal(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
