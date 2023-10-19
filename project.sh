#!/usr/bin/env bash
# -*- mode: shell-script; sh-shell: bash; sh-basic-offset: 4; sh-indentation: 4; coding: utf-8 -*-
# shellcheck shell=bash

# This is going to be the simplest possible thing that could work
# until the full version is in MyCmd itself.
set -o nounset -o errexit -o errtrace -o pipefail

if ! PROJECT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P); then
    echo >&2 "Error fetching project directory."
    exit 1
fi
readonly PROJECT_DIR

if ! PROJECT_SH_FILE="$(grealpath "${BASH_SOURCE[0]}")"; then
    echo >&2 "Error getting project.sh file path."
    exit 1
fi
readonly PROJECT_SH_FILE

declare -a PROJECT_SHELL_FILES=("${PROJECT_SH_FILE}")
readonly PROJECT_SHELL_FILES

function go-lines() {
    go run github.com/segmentio/golines@latest --write-output --max-len=100 --tab-len=8 --reformat-tags --shorten-comments .
}

function go-imports() {
    go run golang.org/x/tools/cmd/goimports@latest -w .
}

function go-fmt() {
    go fmt -x ./...
}

function format() {
    call_tasks format-shell go-imports go-lines
}

function go-mod-tidy() {
    go mod tidy
}

# shellcheck disable=SC2317
function lint-golangci-lint-auto-fix() {
    if ! command -v golangci-lint &>/dev/null; then
        echo >&2 "Please install golangci-lint"
        return 1
    fi

    local -r linter_name="${1}"

    golangci-lint run \
        --verbose \
        --disable-all --enable "${linter_name}" \
        --fix
}

function lint-golangci-lint() {
    if ! command -v golangci-lint &>/dev/null; then
        echo >&2 "Please install golangci-lint"
        return 1
    fi

    golangci-lint run --sort-results
}

function go-vet() {
    go vet
}

function lint-all() {
    call_tasks lint-shell go-vet lint-golangci-lint
}

function run-dump_doc() {
    go run github.com/travisbhartwell/bashdoc/cmd/dump_doc
}

function run-dump-script-structure() {
    go run github.com/travisbhartwell/bashdoc/cmd/dump-script-structure "${@}"
}

function go-doc() {
    if [[ -v TMUX ]]; then
        tmux -CC new-window -c "${PROJECT_DIR}" -d 'go run golang.org/x/tools/cmd/godoc@latest'

        echo >&2 "Waiting for docs site to start."
        curl --silent --head -X GET --retry 20 --retry-connrefused --retry-delay 1 http://localhost:6060

        open http://localhost:6060
    else
        go run golang.org/x/tools/cmd/godoc@latest
    fi
}

function build() {
    go build -v ./...
}

## Shell Script Support
function lint-shell() {
    if (("${#PROJECT_SHELL_FILES[@]}" == 0)); then
        echo >&2 "No shell script files defined, skipping shell lint check."
        return 0
    fi

    echo "Linting the following files:"
    list-shell-files

    cd "${PROJECT_DIR}"
    echo "Running ShellCheck:"
    shellcheck --check-sourced "${PROJECT_SHELL_FILES[@]}"
}

function format-shell() {
    if (("${#PROJECT_SHELL_FILES[@]}" == 0)); then
        echo >&2 "No shell files defined, skipping shell format."
        return 0
    fi

    echo "Formatting the following files:"
    list-shell-files

    cd "${PROJECT_DIR}"
    shfmt --language-dialect bash --indent=4 --binary-next-line --case-indent --write "${PROJECT_SHELL_FILES[@]}"
}

function list-shell-files() {
    list-files "${PROJECT_SHELL_FILES[@]}"
}

function list-files() {
    echo "${*}" | tr ' ' '\n'
}

function list-tasks() {
    declare -F | grep -v \
        -e "^declare -f call_task" \
        -e "^declare -f function_exists" \
        -e "^declare -f list-files" \
        | sed 's/declare -f //' \
        | sort
}

function function_exists() {
    declare -F "$1" >/dev/null
}

# shellcheck disable=SC2317
function call_tasks() {
    for task in "${@}"; do
        local return_code=0

        call_task "${task}" || return_code=$?

        if ((return_code != 0)); then
            return "${return_code}"
        fi
    done
}

function call_task() {
    local -r fn=$1
    shift

    cd "${PROJECT_DIR}"

    local return_code=0
    if function_exists "${fn}"; then
        echo "➡️ Executing task '${fn}'..."

        "${fn}" "$@" || return_code=$?
    else
        echo >&2 "Unknown task: '${fn}'."
        return_code=1
    fi

    if ((return_code == 0)); then
        echo "✅ Task '${fn}' succeeded."
    else
        echo "❌ Task '${fn}' failed."
    fi

    return "${return_code}"
}

if (($# == 0)); then
    echo >&2 "Expecting task to run:"
    echo >&2 "$0 <task>"
    exit 1
fi

call_task "${@}"
