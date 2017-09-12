package testhelpers

import (
	"io/ioutil"
	"testing"
)

// FixtureNames returns the names of all testdata fixtures
func FixtureNames(t *testing.T) []string {
	testDirs, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	result := []string{}
	for _, testDir := range testDirs {
		if !testDir.IsDir() {
			continue
		}
		result = append(result, testDir.Name())
	}
	return result
}
