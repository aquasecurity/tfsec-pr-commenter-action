#!/usr/bin/env bash

set -xe

TFSEC_VERSION=""
if [ "$INPUT_TFSEC_VERSION" != "latest" ]; then
  TFSEC_VERSION="/tags/${INPUT_TFSEC_VERSION}"
fi

wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec/releases${TFSEC_VERSION} -O - | grep -m 1 -o -E "https://.+?tfsec-linux-amd64" | head -n1)" > tfsec-linux-amd64
wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec/releases${TFSEC_VERSION} -O - | grep -m 1 -o -E "https://.+?tfsec_checksums.txt" | head -n1)" > tfsec.checksums

grep tfsec-linux-amd64 tfsec.checksums > tfsec-linux-amd64.checksum
sha256sum -c tfsec-linux-amd64.checksum
install tfsec-linux-amd64 /usr/local/bin/tfsec

COMMENTER_VERSION="latest"
if [ "$INPUT_COMMENTER_VERSION" != "latest" ]; then
  COMMENTER_VERSION="tags/${INPUT_COMMENTER_VERSION}"
fi

wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec-pr-commenter-action/releases/${COMMENTER_VERSION} -O - | grep -o -E "https://.+?commenter-linux-amd64")" > commenter-linux-amd64
wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec-pr-commenter-action/releases/${COMMENTER_VERSION} -O - | grep -o -E "https://.+?checksums.txt")" > commenter.checksums

grep commenter-linux-amd64 commenter.checksums > commenter-linux-amd64.checksum
sha256sum -c commenter-linux-amd64.checksum
install commenter-linux-amd64 /usr/local/bin/commenter

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

tfsec --out=${TFSEC_OUT_OPTION} --format=${TFSEC_FORMAT_OPTION} --soft-fail ${TFSEC_ARGS_OPTION} "${INPUT_WORKING_DIRECTORY}"
commenter
