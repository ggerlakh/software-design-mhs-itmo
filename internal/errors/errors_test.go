package errors

import "testing"

func TestCommandNotFoundError_Error(t *testing.T) {
	err := &CommandNotFoundError{Command: "foo"}
	if err.Error() != "go-cli: command not found: foo" {
		t.Fatalf("неверный текст ошибки: %s", err.Error())
	}
}

func TestIsWrapper(t *testing.T) {
	if !Is(ErrExit, ErrExit) {
		t.Fatalf("Is должен возвращать true для ErrExit")
	}

	if Is(ErrExit, &CommandNotFoundError{Command: "foo"}) {
		t.Fatalf("Is не должен возвращать true для разных ошибок")
	}
}
