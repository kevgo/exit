package helpers

import (
	"log"
	"os"
)

// GetWD returns the current working directory
func GetWD() string {
	result, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return result
}
