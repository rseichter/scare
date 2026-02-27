# scare: An opinionated script care utility

Tired of caring for shell scripts or YAML files by manually invoking tools like
shellcheck or yamllint? _scare_ can do it for you.

Just execute `scare .` to automatically and recursively identify all suitable
files in the current directory and have them checked, linted, and/or formatted
for you. `scare -h` will display the list of supported options.

Copyright © 2026 Ralph Seichter

## Purpose

_scare_ is opinionated, meaning it will deliberately choose particular
behaviour instead of letting the user pick their own. This improves
convenience. Users who crave more flexibility can simply use the underlying
external tools directly. Depending on file type, _scare_ needs the following
third party utilities to be available via your shell's PATH configuration.

* Bash: [shellcheck](https://www.shellcheck.net/), [shfmt](https://github.com/mvdan/sh).
* Python: [Black](https://black.readthedocs.io/), [Flake8](https://flake8.pycqa.org/).
* YAML: [yamllint](https://yamllint.readthedocs.io/)

## Installation

[Go](https://go.dev/) version 1.22 or later needs to be available on your
machine. The details of installing Go depend on your operating system, and are
outside the scope of this document.

Once Go is set up, you can use the following command to install the latest
version of _scare_ using your shell. If you prefer a specific, fixed version,
replace `latest` with a Git tag or commit hash.

```bash
# Install latest development version
go install github.com/rseichter/scare@latest

# Alternative: Install a specific version
# go install github.com/rseichter/scare@v0.3
```
