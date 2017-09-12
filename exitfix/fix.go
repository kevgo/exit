// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "go/ast"

type fix struct {
	name     string
	date     string // date that fix was introduced, in YYYY-MM-DD format
	f        func(*ast.File) bool
	desc     string
	disabled bool // whether this fix should be disabled by default
}

// main runs sort.Sort(byName(fixes)) before printing list of fixes.
type byName []fix

func (f byName) Len() int           { return len(f) }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f byName) Less(i, j int) bool { return f[i].name < f[j].name }

// main runs sort.Sort(byDate(fixes)) before applying fixes.
type byDate []fix

func (f byDate) Len() int           { return len(f) }
func (f byDate) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f byDate) Less(i, j int) bool { return f[i].date < f[j].date }

var fixes []fix

func register(f fix) {
	fixes = append(fixes, f)
}
