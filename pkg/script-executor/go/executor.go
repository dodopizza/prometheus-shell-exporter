package script_executor

import (
	"bytes"
	"os/exec"
	"path/filepath"
)

const (
	ShellTypeDefault = iota
	ShellTypeAutodetect
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

func (se *ScriptExecutor) preparePowerShellScript(script string, args ...string) (cmd *exec.Cmd, err error) {
	ps, err := exec.LookPath("powershell.exe")
	args = append([]string{"-NoProfile", "-NonInteractive", script}, args...)
	cmd = exec.Command(ps, args...)
	return
}

func (se *ScriptExecutor) prepareScript(script string, args ...string) (cmd *exec.Cmd, err error) {
	return exec.Command(script, args...), nil
}

func (se *ScriptExecutor) Execute(script string, args ...string) (stdOut string, stdErr string, err error) {

	var cmd *exec.Cmd

	if se.shellType == ShellTypeAutodetect {
		switch filepath.Ext(script) {
		case ".ps1":
			se.shellType = ShellTypePowerShell
		default:
			se.shellType = ShellTypeDefault
		}
	}

	switch se.shellType {
	case ShellTypePowerShell:
		cmd, err = se.preparePowerShellScript(script, args...)
	default:
		cmd, err = se.prepareScript(script, args...)
	}

	if err != nil {
		return
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
