package main

import (
	scriptExecutor "github.com/dodopizza/prometheus-shell-exporter/pkg/script-executor/go"
)

func main() {
	ps := scriptExecutor.NewScriptExecutor(scriptExecutor.ShellTypeAutodetect)
	stdOut, stdErr, err := ps.Execute("/workspaces/prometheus-shell-exporter/metrics/pse_example.sh")
	if err != nil {
		println(err.Error())
	}
	println(stdOut)
	println(stdErr)
}
