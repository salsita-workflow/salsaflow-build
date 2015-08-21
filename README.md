# salsaflow-build

This repository contains scripts and resources to compile a custom build of
[SalsaFlow](https://github.com/salsaflow/salsaflow).

The extra modules linked in the resulting `salsaflow` executable are following:

* `jira` - JIRA (issue tracker module)
* `reviewboard` - Review Board (code review tool module)

The module sources can be found in `workspace/src/modules`.

## Build

Run `go run tasks/install.go` to build `salsaflow` binaries. The script
uses `go install` by default, so the executables will appear in `bin`
of your current workspace.
