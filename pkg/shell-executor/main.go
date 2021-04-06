package main

import (
	PS "github.com/dodopizza/prometheus-shell-exporter/pkg/shell-executor/go"
)

func main() {
	ps := PS.NewShellExecutor("/workspaces/powershell_exporter/examples/pse_example.sh")
	stdOut, stdErr, err := ps.Execute()
	if err != nil {
		println(err.Error())
	}
	println(stdOut)
	println(stdErr)
}
