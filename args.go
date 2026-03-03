// Copyright © 2026 Ralph Seichter
package main

import (
	"flag"
	"fmt"
)

var (
	failFast   bool
	forcedType string
	maxDepth   int
	quiet      bool
	verPrint   bool
)

func init() {
	flag.BoolVar(&failFast, "f", false, "Fail fast, stop at first reported issue.")
	flag.BoolVar(&quiet, "q", false, "Quieter operation, reduced output.")
	flag.BoolVar(&verPrint, "V", false, "Print version number to stdout.")
	flag.IntVar(&maxDepth, "r", 2, "Recursion depth limit.")
	flag.StringVar(&forcedType, "t", "auto", "File type.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] {path} [path ...]\n", program)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\n%s %s Copyright © 2026 Ralph Seichter\n", program, version)
	}
}
