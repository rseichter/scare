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
	"strings"
)

var (
	failFast bool
	fileType string
	verbose  bool
)

func init() {
	flag.BoolVar(&failFast, "f", false, "fail fast at first reported issue")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&fileType, "t", "auto", "file type")
}

type Filetype int

const (
	Bash Filetype = iota
	Go
	Python
	Yaml
	Unknown = -1
	Prog    = "scare"
	Version = "0.2"
)

func runCmd(cmd *exec.Cmd) (error, int) {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if verbose {
		fmt.Printf("» %s\n", strings.Join(cmd.Args, " "))
	}
	err := cmd.Run()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			code := e.ExitCode()
			if verbose {
				fmt.Printf("⚠️ %s returned code %v\n", strings.Join(cmd.Args, " "), code)
			}
			return nil, code
		} else if e, ok := err.(*exec.Error); ok {
			fmt.Printf("❌  %v\n", e)
		}
	}
	return err, 0
}

func careFor(path string, ft Filetype) error {
	cmds := make([]*exec.Cmd, 2)
	switch ft {
	case Bash:
		cmds = append(cmds, exec.Command("shellcheck", "-s", "bash", path))
		cmds = append(cmds, exec.Command("shfmt", "-ln", "bash", "-s", "-w", path))
	case Go:
		cmds = append(cmds, exec.Command("go", "fmt", path))
	case Python:
		cmds = append(cmds, exec.Command("flake8", path))
		cmds = append(cmds, exec.Command("black", path))
	case Yaml:
		cmds = append(cmds, exec.Command("yamllint", path))
	}
	for _, cmd := range cmds {
		if cmd != nil {
			err, rc := runCmd(cmd)
			if err != nil {
				return err
			} else if rc != 0 && failFast {
				return errors.New("")
			}
		}
	}
	return nil
}

var (
	reBash   = regexp.MustCompile("#!/.+bash")
	rePython = regexp.MustCompile("#!/.+python")
	reYaml   = regexp.MustCompile("ya?ml$")
	ftMap    = map[string]Filetype{
		"bash":   Bash,
		"py":     Python,
		"python": Python,
		"sh":     Bash,
		"yaml":   Yaml,
		"yml":    Yaml,
	}
)

func determineFiletype(path string) (Filetype, error) {
	if ft, found := ftMap[fileType]; found {
		return ft, nil
	} else if fileType != "auto" {
		return Unknown, fmt.Errorf("Unsupported file type %q", fileType)
	}
	ext := filepath.Ext(path)
	if strings.EqualFold(".sh", ext) {
		return Bash, nil
	} else if strings.EqualFold(".go", ext) {
		return Go, nil
	} else if strings.EqualFold(".py", ext) {
		return Python, nil
	} else if nil != reYaml.FindStringIndex(ext) {
		return Yaml, nil
	}
	// Shebang check
	file, err := os.Open(path)
	if err != nil {
		return Unknown, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if nil != reBash.FindStringIndex(line) {
			return Bash, nil
		} else if nil != rePython.FindStringIndex(line) {
			return Python, nil
		}
	}
	return Unknown, nil
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
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] {path} [path ...]\n", Prog)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "%s %s Copyright © 2026 Ralph Seichter\n", Prog, Version)
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
