package shell_executor

import (
	"bytes"
	"os/exec"
)

type ShellExecutor struct {
	shell string
}

func NewShellExecutor(shell string) IShellExecutor {
	return &ShellExecutor{
		shell: shell,
	}
}

func (p *ShellExecutor) Execute(args ...string) (stdOut string, stdErr string, err error) {
	cmd := exec.Command(p.shell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
