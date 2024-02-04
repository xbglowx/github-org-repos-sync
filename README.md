# github-org-repos-sync [![Build and Test](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/build-test.yaml/badge.svg)](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/build-test.yaml) [![golangci-lint](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/golangci-lint.yml) [![CodeQL](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/codeql-analysis.yml)
Sync GitHub organization repos 

## Build
1. `git clone git@github.com:xbglowx/github-org-repos-sync.git`
1. `go get -d .`
1. `go build .`

## Auth
1. Authenticated access to the GitHub API via personal access token with scope **repo**
   * `export GITHUB_TOKEN=<token>`
1. `git` cli authenticated to GitHub

## How It Works
1. Generates a list of repos the caller has access to within the specified organization
1. Clones each repo if it doesn't exist locally in the destionation path
1. If the repo already exists:
   1. Switches to the default branch, only if the current branch is clean; Stashes if dirty
   1. Updates

## Usage
1. `./github-org-repos-sync --help`
