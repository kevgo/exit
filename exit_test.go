package exit

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/Originate/exit/test/helpers"
	"github.com/cf-guardian/guardian/kernel/fileutils"
)

var tempDirName = helpers.CreateTempDir()
var fu = fileutils.New()
var cwd = helpers.GetWD()

func TestFix(t *testing.T) {
	createTestDataDir(t)
	helpers.RunBinary(tempDirName, t)

	// compare the results
	for _, fixtureName := range helpers.FixtureNames(t) {
		t.Run(fixtureName, func(t *testing.T) {
			for _, file := range helpers.GetFileNames(path.Join(cwd, "test", "examples", fixtureName, "new"), t) {
				t.Run(file, func(t *testing.T) {
					actual := helpers.GetFileContent(path.Join(tempDirName, fixtureName, file), t)
					expected := helpers.GetFileContent(path.Join(cwd, "test", "examples", fixtureName, "new", file), t)
					helpers.CompareStrings(actual, expected, t)
				})
			}
		})
	}
}

func createTestDataDir(t *testing.T) {
	for _, fixtureName := range helpers.FixtureNames(t) {
		srcPath := path.Join(cwd, "test", "examples", fixtureName, "old")
		destPath := path.Join(tempDirName, fixtureName)
		err := fu.Copy(destPath, srcPath)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		fmt.Println(destPath)
	}
}

func TestMain(m *testing.M) {
	var err error
	tempDirName, err = ioutil.TempDir("", "exit-specs")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDirName)

	result := m.Run()
	os.Exit(result)
}
