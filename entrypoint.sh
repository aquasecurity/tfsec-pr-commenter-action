#!/bin/bash

set -x

TFSEC_VERSION="latest"
if [ "$INPUT_TFSEC_VERSION" != "latest" ]; then
  TFSEC_VERSION="tags/${INPUT_TFSEC_VERSION}"
fi

wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec/releases/${TFSEC_VERSION} -O - | grep -o -E "https://.+?tfsec-linux-amd64" | head -n1)" > tfsec 
install tfsec /usr/local/bin/

COMMENTER_VERSION="latest"
if [ "$INPUT_COMMENTER_VERSION" != "latest" ]; then
  COMMENTER_VERSION="tags/${INPUT_COMMENTER_VERSION}"
fi

wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec-pr-commenter-action/releases/${COMMENTER_VERSION} -O - | grep -o -E "https://.+?commenter-linux-amd64")" > commenter
install commenter /usr/local/bin/

if [ -n "${GITHUB_WORKSPACE}" ]; then
  cd "${GITHUB_WORKSPACE}" || exit
fi

if ! tfsec --format=json --force-all-dirs "${INPUT_WORKING_DIRECTORY}" 2>/dev/null >results.json; then
  echo "tfsec violations were identified, running commenter..."
  commenter
fi
