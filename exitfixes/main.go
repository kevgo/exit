// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

var (
	fset     = token.NewFileSet()
	exitCode = 0
)

var allowedRewrites = flag.String("r", "",
	"restrict the rewrites to this comma-separated list")

var forceRewrites = flag.String("force", "",
	"force these fixes to run even if the code looks updated")

var allowed, force map[string]bool

var doDiff = flag.Bool("diff", false, "display diffs instead of rewriting files")

// enable for debugging fix failures
const debug = false // display incorrectly reformatted source and exit

func usage() {
	fmt.Fprintf(os.Stderr, "usage: exitfixes [-diff] [-r fixname,...] [-force fixname,...] [path ...]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nAvailable rewrites are:\n")
	sort.Sort(byName(fixes))
	for _, f := range fixes {
		if f.disabled {
			fmt.Fprintf(os.Stderr, "\n%s (disabled)\n", f.name)
		} else {
			fmt.Fprintf(os.Stderr, "\n%s\n", f.name)
		}
		desc := strings.TrimSpace(f.desc)
		desc = strings.Replace(desc, "\n", "\n\t", -1)
		fmt.Fprintf(os.Stderr, "\t%s\n", desc)
	}
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	sort.Sort(byDate(fixes))

	if *allowedRewrites != "" {
		allowed = make(map[string]bool)
		for _, f := range strings.Split(*allowedRewrites, ",") {
			allowed[f] = true
		}
	}

	if *forceRewrites != "" {
		force = make(map[string]bool)
		for _, f := range strings.Split(*forceRewrites, ",") {
			force[f] = true
		}
	}

	if flag.NArg() == 0 {
		if err := processFile("standard input", true); err != nil {
			report(err)
		}
		os.Exit(exitCode)
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, false); err != nil {
				report(err)
			}
		}
	}

	os.Exit(exitCode)
}

const parserMode = parser.ParseComments

func gofmtFile(f *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func processFile(filename string, useStdin bool) error {
	var f *os.File
	var err error
	var fixlog bytes.Buffer

	if useStdin {
		f = os.Stdin
	} else {
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err != nil {
		return err
	}

	// Apply all fixes to file.
	newFile := file
	fixed := false
	for _, fix := range fixes {
		if allowed != nil && !allowed[fix.name] {
			continue
		}
		if fix.disabled && !force[fix.name] {
			continue
		}
		if fix.f(newFile) {
			fixed = true
			fmt.Fprintf(&fixlog, " %s", fix.name)

			// AST changed.
			// Print and parse, to update any missing scoping
			// or position information for subsequent fixers.
			newSrc, err := gofmtFile(newFile)
			if err != nil {
				return err
			}
			newFile, err = parser.ParseFile(fset, filename, newSrc, parserMode)
			if err != nil {
				if debug {
					fmt.Printf("%s", newSrc)
					report(err)
					os.Exit(exitCode)
				}
				return err
			}
		}
	}
	if !fixed {
		return nil
	}
	fmt.Fprintf(os.Stderr, "%s: fixed %s\n", filename, fixlog.String()[1:])

	// Print AST.  We did that after each fix, so this appears
	// redundant, but it is necessary to generate gofmt-compatible
	// source code in a few cases. The official gofmt style is the
	// output of the printer run on a standard AST generated by the parser,
	// but the source we generated inside the loop above is the
	// output of the printer run on a mangled AST generated by a fixer.
	newSrc, err := gofmtFile(newFile)
	if err != nil {
		return err
	}

	// Remove empty lines after the fix
	re := regexp.MustCompile("(exit.If\\(.*?\\))\n")
	newSrc = re.ReplaceAll(newSrc, []byte("$1"))

	if *doDiff {
		data, err := diff(src, newSrc)
		if err != nil {
			return fmt.Errorf("computing diff: %s", err)
		}
		fmt.Printf("diff %s fixed/%s\n", filename, filename)
		os.Stdout.Write(data)
		return nil
	}

	if useStdin {
		os.Stdout.Write(newSrc)
		return nil
	}

	return ioutil.WriteFile(f.Name(), newSrc, 0)
}

var gofmtBuf bytes.Buffer

func gofmt(n interface{}) string {
	gofmtBuf.Reset()
	if err := format.Node(&gofmtBuf, fset, n); err != nil {
		return "<" + err.Error() + ">"
	}
	return gofmtBuf.String()
}

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = processFile(path, false)
	}
	if err != nil {
		report(err)
	}
	return nil
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func writeTempFile(dir, prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := writeTempFile("", "go-fix", b1)
	if err != nil {
		return
	}
	defer os.Remove(f1)

	f2, err := writeTempFile("", "go-fix", b2)
	if err != nil {
		return
	}
	defer os.Remove(f2)

	cmd := "diff"
	if runtime.GOOS == "plan9" {
		cmd = "/bin/ape/diff"
	}

	data, err = exec.Command(cmd, "-u", f1, f2).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}
