package helpers

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

// GetFileContent returns the content of the given file
func GetFileContent(filepath string, t *testing.T) string {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	return strings.TrimSpace(string(content))
}
