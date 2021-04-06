package main

import (
	"errors"
	"os"

	PS "github.com/dodopizza/powershell-exporter/pkg/powershell/go"
)

func main() {
	ps := PS.New()
	stdOut, stdErr, err := ps.Execute("echo Hello")
	if err != nil {
		errors.New(err.Error())
		os.Exit(1)
	}
	println(stdOut)
	println(stdErr)
}
