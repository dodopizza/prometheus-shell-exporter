package script_executor

import (
	"bytes"
	"os/exec"
)

const (
	ShellTypeDefault = iota
	ShellTypePowerShell
	ShellTypeBash
)

type ScriptExecutor struct {
	shellType int
}

func NewScriptExecutor(shellType int) *ScriptExecutor {
	return &ScriptExecutor{
		shellType: shellType,
	}
}

func (se *ScriptExecutor) executePowerShellScript(script string, args ...string) (cmd *exec.Cmd) {
	ps, _ := exec.LookPath("powershell.exe")
	args = append([]string{"-NoProfile", "-NonInteractive", script}, args...)
	return exec.Command(ps, args...)
}

func (se *ScriptExecutor) executeScript(script string, args ...string) (cmd *exec.Cmd) {
	return exec.Command(script, args...)
}

func (se *ScriptExecutor) Execute(script string, args ...string) (stdOut string, stdErr string, err error) {

	var cmd *exec.Cmd

	switch se.shellType {
	case ShellTypePowerShell:
		cmd = se.executePowerShellScript(script, args...)
	default:
		cmd = se.executeScript(script, args...)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
