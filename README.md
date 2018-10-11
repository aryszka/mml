# MML

An alternative frontend for Go and JS.

[Example](https://github.com/aryszka/mml/blob/master/compile.mml)

## Boot the compiler

Prerequisits: Go installed, $GOPATH set and $GOPATH/bin added to $PATH.

```
make deps install
```

Test:

```
mkdir -p hello
echo 'stdout("Hello, world!\n")' > hello/hello.mml
mml < hello/hello.mml > hello/hello.go
go run hello/hello.go
```
