# github-org-repos-sync [![CircleCI](https://circleci.com/gh/xbglowx/github-org-repos-sync.svg?style=svg)](https://circleci.com/gh/xbglowx/github-org-repos-sync)
Sync Github Org's repos

## Build
1. `git clone git@github.com:xbglowx/github-org-repos-sync.git`
1. `go get -d .`
1. `go build .`

## Auth
1. Authenticated access to the GitHub API via personal access token with scope **repo**
   * `export GITHUB_TOKEN=<token>`
1. `git` cli authenticated to GitHub

## Usage
1. `./github-org-repos-sync --help`
