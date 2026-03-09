// Copyright © 2026 Ralph Seichter
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const (
	program = "scare"
	version = "0.5"
)

type ftype int

const (
	ftDunno ftype = iota
	ftBash
	ftGo
	ftPosixSh
	ftPython
	ftYaml
)

func runCmd(c *exec.Cmd) (error, int) {
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	if !quiet {
		fmt.Printf("» %s\n", strings.Join(c.Args, " "))
	}
	err := c.Run()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			code := e.ExitCode()
			if !quiet {
				fmt.Printf("⚠️ %s returned code %v\n", strings.Join(c.Args, " "), code)
			}
			return nil, code
		} else if e, ok := err.(*exec.Error); ok {
			fmt.Printf("❌  %v\n", e)
		}
	}
	return err, 0
}

func careFor(path string, ft ftype) error {
	cmds := make([]*exec.Cmd, 2)
	switch ft {
	case ftBash:
		cmds = append(cmds, exec.Command("shellcheck", "-s", "bash", path))
		cmds = append(cmds, exec.Command("shfmt", "-ln", "bash", "-s", "-w", path))
	case ftGo:
		cmds = append(cmds, exec.Command("go", "fmt", path))
	case ftPosixSh:
		cmds = append(cmds, exec.Command("shellcheck", "-s", "sh", path))
		cmds = append(cmds, exec.Command("shfmt", "-ln", "posix", "-s", "-w", path))
	case ftPython:
		cmds = append(cmds, exec.Command("flake8", path))
		cmds = append(cmds, exec.Command("black", path))
	case ftYaml:
		cmds = append(cmds, exec.Command("yamllint", path))
	}
	for _, c := range cmds {
		if c != nil {
			err, rc := runCmd(c)
			if err != nil {
				return err
			} else if rc != 0 && failFast {
				return errors.New("Fail-fast requested, exiting.")
			}
		}
	}
	return nil
}

var (
	reBash    = regexp.MustCompile("#!/.+bash")
	rePosixSh = regexp.MustCompile("#!/.+[ /]sh")
	rePython  = regexp.MustCompile("#!/.+python")
	reYaml    = regexp.MustCompile("ya?ml$")
	ftMap     = map[string]ftype{
		"bash":   ftBash,
		"py":     ftPython,
		"python": ftPython,
		"sh":     ftPosixSh,
		"yaml":   ftYaml,
		"yml":    ftYaml,
	}
)

func ftChoices(m map[string]ftype) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return strings.Join(keys, ", ")
}

func determineFiletype(path string) (ftype, error) {
	if ft, found := ftMap[forcedType]; found {
		return ft, nil
	} else if forcedType != "auto" {
		return ftDunno, fmt.Errorf("Unsupported file type %q (choices: %s)", forcedType, ftChoices(ftMap))
	}
	ext := filepath.Ext(path)
	if strings.EqualFold(".sh", ext) {
		return ftBash, nil
	} else if strings.EqualFold(".go", ext) {
		return ftGo, nil
	} else if strings.EqualFold(".py", ext) {
		return ftPython, nil
	} else if nil != reYaml.FindStringIndex(ext) {
		return ftYaml, nil
	}
	// Shebang check
	file, err := os.Open(path)
	if err != nil {
		return ftDunno, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if nil != reBash.FindStringIndex(line) {
			return ftBash, nil
		} else if nil != rePosixSh.FindStringIndex(line) {
			return ftPosixSh, nil
		} else if nil != rePython.FindStringIndex(line) {
			return ftPython, nil
		}
	}
	return ftDunno, nil
}

func walkDirFunc(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	} else if n := strings.Count(path, string(os.PathSeparator)); entry.IsDir() && n >= maxDepth {
		if !quiet {
			fmt.Fprintf(os.Stderr, "⛔ Maximum depth %d reached, skipping %q\n", maxDepth, path)
		}
		return filepath.SkipDir
	} else if entry.IsDir() && entry.Name() == ".git" {
		// Skip local Git repository.
		return filepath.SkipDir
	} else if strings.HasPrefix(entry.Name(), ".git") {
		// Do nothing for files with .git prefix.
	} else if !entry.IsDir() {
		ft, err := determineFiletype(path)
		if err != nil {
			return err
		}
		if err = careFor(path, ft); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if verPrint {
		fmt.Println(version)
	} else if flag.NArg() < 1 {
		// Missing path spec.
		flag.Usage()
	} else {
		for _, arg := range flag.Args() {
			// Trailing path separators interfere with depth counting, strip them.
			a := strings.TrimRight(arg, string(os.PathSeparator))
			if err := filepath.WalkDir(a, walkDirFunc); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	}
}
