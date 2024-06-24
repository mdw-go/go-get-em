# go-get-em

A tool for finding and updating outdated dependencies in a Go project.

## Installation

`make install`

## Usage

1. Runs `go list -m -u -json all` and scans the output for outdated dependencies.
2. Emits a bitbucket/github URL to allow review of changes since current version.
3. Emits `go get -u` for each outdated dependency.
4. Execute emitted commands as desired.


----

Relevant documentation concerning the `go list` command being used:

```
$ go help list
usage: go list [-f format] [-json] [-m] [list flags] [build flags] [packages]

List lists the named packages, one per line.
...
The -m flag causes list to list modules instead of packages.
...
The arguments to list -m are interpreted as a list of modules, not packages.
The main module is the module containing the current directory.
The active modules are the main module and its dependencies.
With no arguments, list -m shows the main module.
With arguments, list -m shows the modules specified by the arguments.
Any of the active modules can be specified by its module path.
The special pattern "all" specifies all the active modules, first the main
module and then dependencies sorted by module path.
...
The -u flag adds information about available upgrades.
For example, 'go list -m -u all' might print:

    my/main/module
    golang.org/x/text v0.3.0 [v0.4.0] => /tmp/text
    rsc.io/pdf v0.1.1 (retracted) [v0.1.2]

(For tools, 'go list -m -u -json all' may be more convenient to parse.)
```