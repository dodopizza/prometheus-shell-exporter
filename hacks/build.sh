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
if [ -n "${APP_VERSION:-}" ]; then
    ldflags+=" -X main.appVersion=${APP_VERSION}"
fi

echo "[~] Cleanup bin/ dir"
rm -rf "${project_dir}/bin/"

for arc in amd64; do
    for os in linux windows; do
        for app in 'shell-exporter'; do

            app_folder="${app}-${APP_VERSION:-draft}.${os}.${arc}"
            app_folder_abs="${project_dir}/bin/${app_folder}"
            app_extension=""

            case ${os} in
                windows) app_extension=".exe" ;;
            esac

            echo "[~] Build ${app} for ${os} ${arc}"
            GOOS=${os} GOARCH=${arc} CGO_ENABLED=0 go build \
                -ldflags="${ldflags}" \
                -o "${app_folder_abs}/${app}${app_extension}" \
                "${project_dir}/cmd/${app}"


            for artifact in 'metrics'; do
                echo "[~] Add ${artifact} artifacts"
                cp -r "${project_dir}/${artifact}" "${app_folder_abs}/"
            done

            echo "[~] Archive to ${app_folder}.zip"
            (
                cd "${project_dir}/bin/"
                zip -r "${app_folder}.zip" "${app_folder}"
            )
        done
    done
done

echo "[.] Done"