# go-get-em

A tool for finding and updating outdated dependencies in a Go project.

## Installation

`make install`

## Usage

1. Runs `go list -m -u -json all` and scans the output for outdated dependencies.
2. Emits a bitbucket/github URL to allow review of changes since current version.
3. Emits `go get -u` for each outdated dependency.
4. Use executes emitted commands as desired.
