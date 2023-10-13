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

function _function_exists() {
    declare -F "$1" >/dev/null
}

function _call_task() {
    local -r fn=$1
    shift

    if _function_exists "${fn}"; then
        echo >&2 "Executing task: ${fn}..."
        "${fn}" "$@"
    else
        echo >&2 "Unknown task: '${fn}'."
    fi
}

function go-lines() {
    go run github.com/segmentio/golines@latest --write-output --max-len=100 --tab-len=8 --reformat-tags --shorten-comments .
}

function go-imports() {
    go run golang.org/x/tools/cmd/goimports@latest -w .
}

function go-fmt() {
    go fmt -x ./...
}

function fmt() {
    _call_task go-imports
    _call_task go-lines
}

function go-mod-tidy() {
    go mod tidy
}

function lint-golangci-lint() {
    if ! command -v golangci-lint &> /dev/null; then
        echo >&2 "Please install golangci-lint"
        return 1
    fi

    golangci-lint run --sort-results
}

function go-vet() {
    go vet
}

function lint-all() {
    _call_task go-vet
    _call_task lint-golangci-lint
}

function run-dump_doc() {
    go run github.com/travisbhartwell/bashdoc/cmd/dump_doc
}

function build() {
    go build -v ./...
}

if (($# == 0)); then
    echo >&2 "Expecting task to run:"
    echo >&2 "$0 <task>"
    exit 1
fi

_call_task "${@}"
