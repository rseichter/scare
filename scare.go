package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var fileType string

func init() {
	const (
		auto  = "auto"
		usage = "file type"
	)
	flag.StringVar(&fileType, "type", auto, usage)
	flag.StringVar(&fileType, "t", auto, usage)
}

type Filetype int

const (
	Bash Filetype = iota
	Go
	Yaml
	Unknown = -1
)

func careFor(path string, ft Filetype) error {
	var cmd *exec.Cmd
	switch ft {
	case Bash:
		cmd = exec.Command("shellcheck", "--shell=bash", path)
	case Go:
		cmd = exec.Command("go", "fmt", path)
	case Yaml:
		cmd = exec.Command("yamllint", path)
	default:
		return nil
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	fmt.Printf("\n>> %v\n", cmd.Args[0:])
	err := cmd.Run()
	// if err != nil {
	// 	fmt.Printf("%v returned %v\n", cmd.Path, err)
	// }
	fmt.Println("<<")
	return err
}

var reBash = regexp.MustCompile("#!/.+bash")
var reYaml = regexp.MustCompile("ya?ml$")

func determineFiletype(path string) (Filetype, error) {
	ext := filepath.Ext(path)
	if strings.EqualFold(".sh", ext) {
		return Bash, nil
	} else if strings.EqualFold(".go", ext) {
		return Go, nil
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
		}
	}
	return Unknown, nil
}

func walkDirFunc(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if entry.IsDir() && entry.Name() == ".git" {
		return filepath.SkipDir
	} else if strings.HasPrefix(entry.Name(), ".git") {
		// NOOP
	} else if !entry.IsDir() {
		ft, err := determineFiletype(path)
		if err != nil {
			return err
		}
		careFor(path, ft)
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] {path} [path ...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}
	// TODO Use fileType variable
	for _, arg := range flag.Args() {
		err := filepath.WalkDir(arg, walkDirFunc)
		if err != nil {
			log.Printf("Error walking the path %q: %v\n", arg, err)
			return
		}
	}
}
