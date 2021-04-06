package main

import (
	PS "github.com/dodopizza/powershell-exporter/pkg/shell-executor/go"
)

func main() {
	ps := PS.NewShellExecutor("bash")
	stdOut, stdErr, err := ps.Execute("-c", "/workspaces/powershell_exporter/examples/pse_example.sh")
	if err != nil {
		println(err.Error())
	}
	println(stdOut)
	println(stdErr)
}
