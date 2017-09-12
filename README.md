# Nifty error handling for Golang

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

## Automated conversion

A [gofix](https://blog.golang.org/introducing-gofix)-compatible tool
provided by this library
automatically converts all possible expressions in the source code.
Install it via:

```
go get github.com/Originate/exit
```

then run:

```
exitfix .
```

It can also stream source code via STDIN and STDOUT.
Full usage details including list of available fixes:

```
exitfix -h
```
