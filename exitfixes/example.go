package main

import (
	"errors"
	"fmt"

	"github.com/Originate/exit"
)

func test() {
	err := errors.New("foo")
	exit.If(err)

	fmt.Println("done")
	exit.If(err)

	fmt.Println("done")
}
