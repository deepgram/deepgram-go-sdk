#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

module_path="$(go list -m -f '{{.Path}}')"
sdk_major="$(grep -Eo 'sdkVersion string = "v[0-9]+' pkg/client/interfaces/v1/utils.go | grep -Eo '[0-9]+$')"

if [[ "${module_path}" != */v"${sdk_major}" ]]; then
  printf 'SDK version v%s and module path %s use different major versions.\n' "${sdk_major}" "${module_path}"
  exit 1
fi

unexpected_imports="$(git grep -nE '"github\.com/deepgram/deepgram-go-sdk/v[0-9]+(/[^" ]*)?"' -- '*.go' | grep -v "\"${module_path}" || true)"
if [[ -n "${unexpected_imports}" ]]; then
  printf 'Self-imports must use the module path %s:\n%s\n' "${module_path}" "${unexpected_imports}"
  exit 1
fi
