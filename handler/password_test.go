package handler

import (
	"errors"
	"os/user"
	"testing"

	"github.com/xo/usql/rline"
	"github.com/xo/usql/text"
)

func TestPasswordUsesReadlinePrompt(t *testing.T) {
	t.Cleanup(func() {
		readPasswordFallback = promptPasswordFallback
	})
	readPasswordFallback = func(string) (string, error) {
		t.Fatal("unexpected fallback password prompt")
		return "", nil
	}

	h := &Handler{
		l: &rline.Rline{
			Pw: func(prompt string) (string, error) {
				if prompt != text.EnterPassword {
					t.Fatalf("unexpected prompt: %q", prompt)
				}
				return "secret", nil
			},
		},
		user: &user.User{Username: "tester"},
	}

	dsn, err := h.Password("postgres://tester@localhost/postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := dsn, "postgres://tester:secret@localhost/postgres"; got != want {
		t.Fatalf("unexpected dsn: got %q want %q", got, want)
	}
}

func TestPasswordFallsBackWhenReadlinePasswordUnavailable(t *testing.T) {
	t.Cleanup(func() {
		readPasswordFallback = promptPasswordFallback
	})
	readPasswordFallback = func(prompt string) (string, error) {
		if prompt != text.EnterPassword {
			t.Fatalf("unexpected prompt: %q", prompt)
		}
		return "fallback", nil
	}

	h := &Handler{
		l: &rline.Rline{
			Pw: func(string) (string, error) {
				return "", rline.ErrPasswordNotAvailable
			},
		},
		user: &user.User{Username: "tester"},
	}

	dsn, err := h.Password("postgres://tester@localhost/postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := dsn, "postgres://tester:fallback@localhost/postgres"; got != want {
		t.Fatalf("unexpected dsn: got %q want %q", got, want)
	}
}

func TestPasswordDoesNotFallbackOnOtherErrors(t *testing.T) {
	t.Cleanup(func() {
		readPasswordFallback = promptPasswordFallback
	})
	readPasswordFallback = func(string) (string, error) {
		t.Fatal("unexpected fallback password prompt")
		return "", nil
	}

	wantErr := errors.New("boom")
	h := &Handler{
		l: &rline.Rline{
			Pw: func(string) (string, error) {
				return "", wantErr
			},
		},
		user: &user.User{Username: "tester"},
	}

	_, err := h.Password("postgres://tester@localhost/postgres")
	if !errors.Is(err, wantErr) {
		t.Fatalf("unexpected error: got %v want %v", err, wantErr)
	}
}
