package shell_executor

import (
	"bytes"
	"os/exec"
)

type PowerShell struct {
	shell string
}

func NewPowerShellExecutor() IShellExecutor {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		shell: ps,
	}
}

func (p *PowerShell) Execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.shell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
