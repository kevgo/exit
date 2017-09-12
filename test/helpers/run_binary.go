package helpers

import (
	"fmt"
	"os/exec"
	"path"
	"testing"
)

// RunBinary runs the "exitfixes" binary in the test directory
func RunBinary(dirname string, t *testing.T) {
	cmd := exec.Cmd{
		Path: path.Join(GetGoPath(), "bin", "exitfixes"),
		Args: []string{"exitfixes", "."},
		Dir:  dirname,
	}
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		t.FailNow()
	}
}
