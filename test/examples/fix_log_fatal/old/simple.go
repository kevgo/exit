package main

import (
	"errors"
	"log"
)

func test() {
	err := errors.New("foo")
	if err != nil {
		log.Fatal(err)
	}
}
