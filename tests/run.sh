#!/usr/bin/env bash
docker run --rm -it -w /src -v $(PWD):/src:ro -e GITHUB_TOKEN alpine sh -c "./tests/setup.sh"
