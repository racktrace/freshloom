package system

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor is the single entry point for running external commands.
// It is an interface so tests can replace it with a fake implementation.
type Executor interface {
	Run(name string, args ...string) (string, error)
}

// ShellExecutor runs real commands via os/exec.
type ShellExecutor struct{}

func (ShellExecutor) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	output := strings.TrimSpace(out.String())
	if err != nil {
		return output, fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return output, nil
}