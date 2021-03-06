# Nifty error handling for Golang

[![Go Report Card](https://goreportcard.com/badge/github.com/Originate/exit)](https://goreportcard.com/report/github.com/Originate/exit)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com)

This library provides helper methods to dry up repetitive boilerplate around error checking in Golang.

Instead of:

```go
import "log"

if err != nil {
  log.Fatal(err)
}
```

you can write:

```go
import "github.com/Originate/exit"

exit.If(err)
```

The `IfWrap` and `IfWrapf` functions wrap the given error
into the given error message using [errors.Wrap](https://godoc.org/github.com/pkg/errors#Wrap)
and [errors.Wrapf](https://godoc.org/github.com/pkg/errors#Wrapf):

```go
exit.IfWrap(err, "something went wrong")
exit.IfWrapf(err, "%s", message)
```

This makes the most sense for critical errors in Go-based CLI tools,
but could be useful elsewhere.

## Installation

```
go get github.com/Originate/exit
```

A gofix tool to change all compatible usages is available at [exitfix](https://github.com/kevgo/exitfix).
