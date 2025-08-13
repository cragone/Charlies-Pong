package tui

import (
	"os"

	"golang.org/x/term"
)

// used to read key strokes sits in a go routine and reads to a channel
func ReadSingleKey(inputChan chan string) {
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b []byte = make([]byte, 1)
	for {
		if _, err := os.Stdin.Read(b); err == nil {
			inputChan <- string(b[0])
		}
	}
}
