#!/bin/bash

set -eu

script_dir="$(
    cd "$(dirname "${BASH_SOURCE[0]}")" || exit 1
    pwd
)"
project_dir="$(
    cd "${script_dir}/.." || exit 1
    pwd
)"

ldflags="-s -w"

for os in linux windows; do
    GOOS=${os} GOARCH='amd64' CGO_ENABLED=0 go build \
        -ldflags="${ldflags}" \
        -o "${project_dir}/bin/shell-exporter-${os}" \
        "${project_dir}/cmd/shell-exporter"
done