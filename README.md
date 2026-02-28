# scare: An opinionated script care utility

Tired of caring for shell scripts or YAML files by manually invoking tools like
shellcheck or yamllint? _scare_ can do it for you.

Just execute `scare .` to automatically and recursively identify all suitable
files in the current directory and have them checked, linted, and/or formatted
for you. `scare -h` will display the list of supported options, as shown
[here](USAGE.md).

Copyright © 2026 Ralph Seichter

## Purpose

_scare_ is opinionated, meaning it will deliberately choose particular
behaviour instead of letting the user pick their own. This improves
convenience. Users who crave more flexibility can simply use the underlying
external tools directly. Depending on file type, _scare_ needs the following
third party utilities to be available via your shell's PATH configuration.

* Bash & POSIX shells: [shellcheck](https://www.shellcheck.net/), [shfmt](https://github.com/mvdan/sh).
* Python: [Black](https://black.readthedocs.io/), [Flake8](https://flake8.pycqa.org/).
* YAML: [yamllint](https://yamllint.readthedocs.io/)

## Installation

[Go](https://go.dev/) version 1.21 or later needs to be available on your
machine. The details of installing Go depend on your operating system, and are
outside the scope of this document.

Once Go is set up, you can use the following command to install the latest
version of _scare_ using your shell. If you prefer a specific, fixed version,
replace `latest` with a Git tag or commit hash.

```bash
# Install latest development version
go install github.com/rseichter/scare@latest

# Alternative: Install a specific version
# go install github.com/rseichter/scare@v0.4
```

## Behaviour

_scare_ expects one or more file system paths as arguments, which will be
processed in order. If a given path represents a directory, _scare_ will
process the contents recursively, unless the maximum configured depth limit is
reached.

Use the `-r {depth}` option to specify a recursion limit. Note that "depth" is
determined by counting the number of OS-specific path separators. For example,
`dir/subdir/subsubdir` has a depth of 2 on Unix-like systems.`-r 0` disables
recursion.
