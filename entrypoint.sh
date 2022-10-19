#!/usr/bin/env bash

set -xe

TFSEC_VERSION=""
if [ "$INPUT_TFSEC_VERSION" != "latest" ] && [ -n "$INPUT_TFSEC_VERSION" ]; then
  TFSEC_VERSION="/tags/${INPUT_TFSEC_VERSION}"
else
  TFSEC_VERSION="/latest"
fi

COMMENTER_VERSION="latest"
if [ "$INPUT_COMMENTER_VERSION" != "latest" ] && [ -n "$INPUT_COMMENTER_VERSION" ]; then
  COMMENTER_VERSION="/tags/${INPUT_COMMENTER_VERSION}"
else
  COMMENTER_VERSION="/latest"
fi

function get_release_assets {
  repo="$1"
  version="$2"
  args=(
    -sSL
    --header "Accept: application/vnd.github+json"
  )
  curl "${args[@]}" "https://api.github.com/repos/$repo/releases${version}" | jq '.assets[] | { name: .name, download_url: .browser_download_url }'
}

function install_release {
  repo="$1"
  version="$2"
  binary="$3-linux-amd64"
  checksum="$4"
  release_assets="$(get_release_assets "${repo}" "${version}")"

  curl -sLo "${binary}" "$(echo "${release_assets}" | jq -r ". | select(.name == \"${binary}\") | .download_url")"
  curl -sLo "$3-checksums.txt" "$(echo "${release_assets}" | jq -r ". | select(.name | contains(\"$checksum\")) | .download_url")"

  grep "${binary}" "$3-checksums.txt" | sha256sum -c -
  install "${binary}" "/usr/local/bin/${3}"
}

install_release aquasecurity/tfsec "${TFSEC_VERSION}" tfsec tfsec_checksums.txt
install_release aquasecurity/tfsec-pr-commenter-action "${COMMENTER_VERSION}" commenter checksums.txt

if [ -n "${GITHUB_WORKSPACE}" ]; then
  cd "${GITHUB_WORKSPACE}" || exit
fi

if [ -n "${INPUT_TFSEC_ARGS}" ]; then
  TFSEC_ARGS_OPTION="${INPUT_TFSEC_ARGS}"
fi

TFSEC_FORMAT_OPTION="json"
TFSEC_OUT_OPTION="results.json"
if [ -n "${INPUT_TFSEC_FORMATS}" ]; then
  TFSEC_FORMAT_OPTION="${TFSEC_FORMAT_OPTION},${INPUT_TFSEC_FORMATS}"
  TFSEC_OUT_OPTION="${TFSEC_OUT_OPTION%.*}"
fi

tfsec --out=${TFSEC_OUT_OPTION} --format="${TFSEC_FORMAT_OPTION}" --soft-fail "${TFSEC_ARGS_OPTION}" "${INPUT_WORKING_DIRECTORY}"
commenter
