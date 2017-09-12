package exit

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/Originate/exit/testhelpers"

	"github.com/cf-guardian/guardian/kernel/fileutils"
)

var tempDirName = testhelpers.CreateTempDir()
var fu = fileutils.New()
var cwd = testhelpers.GetWD()

func TestFix(t *testing.T) {
	createTestDataDir(t)
	testhelpers.RunBinary(tempDirName, t)

	// compare the results
	for _, fixtureName := range testhelpers.FixtureNames(t) {
		t.Run(fixtureName, func(t *testing.T) {
			for _, file := range testhelpers.GetFileNames(path.Join(cwd, "testdata", fixtureName, "new"), t) {
				t.Run(file, func(t *testing.T) {
					actual := testhelpers.GetFileContent(path.Join(tempDirName, fixtureName, file), t)
					expected := testhelpers.GetFileContent(path.Join(cwd, "testdata", fixtureName, "new", file), t)
					testhelpers.CompareStrings(actual, expected, t)
				})
			}
		})
	}
}

func createTestDataDir(t *testing.T) {
	for _, fixtureName := range testhelpers.FixtureNames(t) {
		srcPath := path.Join(cwd, "testdata", fixtureName, "old")
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
