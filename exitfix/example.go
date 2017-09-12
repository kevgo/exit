package main

import (
	"errors"
	"fmt"
	"log"
)

func test() {
	err := errors.New("foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}
