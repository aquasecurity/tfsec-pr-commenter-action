#!/usr/bin/env sh
echo "* Running in a container"

export GITHUB_REPOSITORY=aquasecurity/tfsec-pr-commenter-action
export GITHUB_WORKSPACE=/github/workspace
export INPUT_GITHUB_TOKEN=${GITHUB_TOKEN}

mkdir -p /github/workflow
mkdir -p /github/workspace/vendor

# Here we pretend to be a real PR where a file is changed in a directory:
# https://github.com/aquasecurity/tfsec-pr-commenter-action/pull/43
cat >/github/workflow/event.json<<EOF
{"number":43}
EOF

# Mock the file changed in this PR (vendor/modules.txt)
for i in seq 1 10; do
  echo this is line $i >> ${GITHUB_WORKSPACE}/vendor/modules.txt
done

# Mock another file in this PR, this time in / (go.mod)
for i in seq 1 10; do
  echo this is line $i >> ${GITHUB_WORKSPACE}/go.mod
done

cp ./bin/commenter-linux-amd64 ${GITHUB_WORKSPACE}/
cp ./tests/results.json ${GITHUB_WORKSPACE}/
cd $GITHUB_WORKSPACE

# This will run through the PR and WRITE A COMMENT!
./commenter-linux-amd64
