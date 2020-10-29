#!/bin/bash

set -x

if [ -n "${GITHUB_WORKSPACE}" ]; then
  cd "${GITHUB_WORKSPACE}" || exit
fi



if ! tfsec --format=json "${INPUT_WORKING_DIRECTORY}" > results.json
then
  echo "tfsec errors occurred, running commenter "
  /commenter
fi
