//go:build windows

package handler

import "github.com/xo/usql/rline"

func promptPasswordFallback(string) (string, error) {
	return "", rline.ErrPasswordNotAvailable
}
