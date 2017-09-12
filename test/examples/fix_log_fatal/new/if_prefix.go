package main

import (
	"errors"
	"github.com/Originate/exit"
	"log"
)

func test() {
	err := errors.New("foo")
	exit.If(err)
}
