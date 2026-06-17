package process

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ExitCode int

func (e ExitCode) Error() string {
	return fmt.Sprintf("終了コード %d", int(e))
}

func Run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return commandError(cmd.Run())
}

func FindCommand(name string, fallback string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	if fallback != "" {
		if stat, err := os.Stat(fallback); err == nil && !stat.IsDir() {
			return fallback, nil
		}
	}
	return "", fmt.Errorf("%s が見つかりません", name)
}

func ExitError(err error) {
	var code ExitCode
	if errors.As(err, &code) {
		os.Exit(int(code))
	}

	var buffer bytes.Buffer
	buffer.WriteString(err.Error())
	if !strings.HasSuffix(buffer.String(), "\n") {
		buffer.WriteString("\n")
	}
	fmt.Fprint(os.Stderr, buffer.String())
	os.Exit(1)
}

func commandError(err error) error {
	if err == nil {
		return nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return ExitCode(exitErr.ExitCode())
	}
	return err
}
