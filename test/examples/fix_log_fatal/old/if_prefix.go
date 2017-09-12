package main

import (
	"errors"
	"log"
)

func test() {
	if err := errors.New("foo"); err != nil {
		log.Fatal(err)
	}
}
