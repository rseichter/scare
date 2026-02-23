# scare: An opinionated script care utility

Tired of caring for shell scripts or YAML files by manually invoking tools like
shellcheck or yamllint? _scare_ can do it for you. Just execute `scare .` to
automatically and recursively identify all suitable files in the current
directory and have them checked, linted, and/or formatted for you. `scare -h`
will display the list of supported options.

_scare_ is opinionated, meaning it will deliberately choose particular
behaviour instead of letting the user pick their own. This improves
convenience. Users who crave more flexibility can simply use the underlying
external tools directly. Depending on file type, _scare_ needs the following
third party utilities to be available via your shell's PATH configuration.

* Bash: [shellcheck](https://www.shellcheck.net/), [shfmt](https://github.com/mvdan/sh).
* Python: [Black](https://black.readthedocs.io/), [Flake8](https://flake8.pycqa.org/).
* YAML: [yamllint](https://yamllint.readthedocs.io/)

Copyright © 2026 Ralph Seichter
