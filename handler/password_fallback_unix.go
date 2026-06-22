//go:build !windows

package handler

import (
	"os"

	"golang.org/x/term"
)

func promptPasswordFallback(prompt string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", err
	}
	defer tty.Close()

	if _, err := tty.WriteString(prompt); err != nil {
		return "", err
	}
	buf, err := term.ReadPassword(int(tty.Fd()))
	if _, writeErr := tty.WriteString("\n"); err == nil && writeErr != nil {
		err = writeErr
	}
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
