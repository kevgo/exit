package helpers

import (
	"io/ioutil"
	"log"
)

// CreateTempDir creates a temporary directory and returns its path
func CreateTempDir() string {
	tempDirName, err := ioutil.TempDir("", "exit-specs")
	if err != nil {
		log.Fatal(err)
	}
	return tempDirName
}
