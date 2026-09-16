#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

read_sdk_major() {
  local sdk_version_file="${1:-pkg/client/interfaces/v1/utils.go}"
  local sdk_major

  sdk_major="$(grep -Eo 'sdkVersion[[:space:]]+string[[:space:]]*=[[:space:]]*"v[0-9]+' "${sdk_version_file}" | grep -Eo '[0-9]+$' || true)"
  [[ -n "${sdk_major}" ]] || return 1
  printf '%s\n' "${sdk_major}"
}

main() {
  local module_path
  local sdk_major
  local unexpected_imports

  module_path="$(go list -m -f '{{.Path}}')"
  if ! sdk_major="$(read_sdk_major)"; then
    printf 'Could not read sdkVersion from pkg/client/interfaces/v1/utils.go\n'
    exit 1
  fi

  if [[ "${module_path}" != */v"${sdk_major}" ]]; then
    printf 'SDK version v%s and module path %s use different major versions.\n' "${sdk_major}" "${module_path}"
    exit 1
  fi

  unexpected_imports="$(git grep -nE "github\\.com/deepgram/deepgram-go-sdk/v[0-9]+(/[^[:space:]\\\`\"']*)?" | grep -vF "${module_path}" || true)"
  if [[ -n "${unexpected_imports}" ]]; then
    printf 'Self-imports must use the module path %s:\n%s\n' "${module_path}" "${unexpected_imports}"
    exit 1
  fi
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
