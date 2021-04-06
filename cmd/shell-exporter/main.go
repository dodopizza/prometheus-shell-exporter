package main

import (
	"errors"
	"os"
)

func main() {

	p1 := PromMetrics{}
	p2 := PromMetrics{}
	p3 := PromMetrics{}
	p1.ReadFromFile("/workspaces/powershell_exporter/examples/pse_tcp_connection_metrics.example.json")
	p2.ReadFromFile("/workspaces/powershell_exporter/examples/pse_tcp_dynamic_port_range_number_of_ports.example.json")
	p3.ReadFromFile("/workspaces/powershell_exporter/examples/pse_example.json")

	pe := NewPromExporter()
	pe.NewGaugeVecFromPromMetrics("HelloKitty", p1)
	pe.NewGaugeVecFromPromMetrics("HelloPuppy", p2)
	pe.NewGaugeVecFromPromMetrics("HelloGleb", p3)
	err := pe.Serve()
	if err != nil {
		errors.New(err.Error())
	}

	os.Exit(0)
}
