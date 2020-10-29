#!/bin/bash

set -x

if [ -n "${GITHUB_WORKSPACE}" ]; then
  cd "${GITHUB_WORKSPACE}" || exit
fi

env

cat /github/workflow/event.json