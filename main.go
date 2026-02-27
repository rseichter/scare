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
	version = "0.4.dev1"
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

var (
	failFast   bool
	forcedType string
	quiet      bool
)

func init() {
	flag.BoolVar(&failFast, "f", false, "Fail fast, stop at first reported issue.")
	flag.BoolVar(&quiet, "q", false, "Quieter operation, reduced output.")
	flag.StringVar(&forcedType, "t", "auto", "File type.")
}

func runCmd(cmd *exec.Cmd) (error, int) {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if !quiet {
		fmt.Printf("» %s\n", strings.Join(cmd.Args, " "))
	}
	err := cmd.Run()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			code := e.ExitCode()
			if !quiet {
				fmt.Printf("⚠️ %s returned code %v\n", strings.Join(cmd.Args, " "), code)
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
	for _, cmd := range cmds {
		if cmd != nil {
			err, rc := runCmd(cmd)
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
		return ftDunno, fmt.Errorf("Unsupported file type %q (valid choices: %s)", forcedType, ftChoices(ftMap))
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
	} else if entry.IsDir() && entry.Name() == ".git" {
		return filepath.SkipDir
	} else if strings.HasPrefix(entry.Name(), ".git") {
		// NOOP
	} else if !entry.IsDir() {
		ft, err := determineFiletype(path)
		if err != nil {
			return err
		}
		err = careFor(path, ft)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] {path} [path ...]\n", program)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\n%s %s Copyright © 2026 Ralph Seichter\n", program, version)
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	for _, arg := range flag.Args() {
		if err := filepath.WalkDir(arg, walkDirFunc); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
