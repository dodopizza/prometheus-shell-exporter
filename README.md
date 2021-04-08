# Prometheus Shell Exporter

Shell Exporter can execute `Powershell` or `Bash` scripts and transform its output to Prometheus metrics.

## Metrics

Shell Exporter exposes metrics based on your `Powershell` or `Bash` scripts names. Examples of scripts can be finded at `metrics` folder

| Script name                                      | Metric name                                  |
| -------------------------------------------------|----------------------------------------------|
| `bash_gauge.sh`                                  | `bash_gauge`                                 |
| `pse_tcp_connection_metrics.ps1`                 | `pse_tcp_connection_metrics`                 |
| `pse_tcp_dynamic_port_range_number_of_ports.ps1` | `pse_tcp_dynamic_port_range_number_of_ports` |

## Startup options

| Option    | Default value | Description           |
| ----------|---------------|-----------------------|
| --f       | `./metrics`   | scripts directory     |
| --port    | 9360          | exporter port         |
| --help    | -             | show help             |
| --version | -             | show exporter version |


## Installing as Windows Service

1. Download binary
2. Install [nssm](https://nssm.cc/) via [Chocolatey](https://community.chocolatey.org/packages/NSSM)

| Action    | Command                                                          |
| ----------|------------------------------------------------------------------|
| install   | nssm install shell_exporter C:\shell-exporter.exe -f C:\scripts  |
| remove    | nssm remove shell_exporter confirm                               |

## Development

1. Make branch from `master`
2. Reopen repo at [vscode container](https://code.visualstudio.com/docs/remote/containers)
3. Make necessary changes
4. Push it to branch and make pull-request

### Makefile targets

| Target    | Action                                                                                  |
| ----------|-----------------------------------------------------------------------------------------|
| build     | build binary for `linux` and `windows` platforms                                        |
| tidy      | remove all dependencies from the go.mod file which are not required in the source files |
| test      | run unit tests                                                                          |
| lint      | run linting                                                                             |
| run       | run exporter locally                                                                    |
