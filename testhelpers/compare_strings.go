package testhelpers

import (
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// CompareStrings compares the given strings
func CompareStrings(actual, expected string, t *testing.T) {
	if actual != expected {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(actual, expected, false)
		fmt.Println(dmp.DiffPrettyText(diffs))
		t.FailNow()
	}
}
