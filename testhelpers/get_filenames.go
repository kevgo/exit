package testhelpers

import (
	"fmt"
	"io/ioutil"
	"testing"
)

// GetFileNames returns the names of all files in the given directory
func GetFileNames(dirname string, t *testing.T) []string {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	result := []string{}
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result
}
