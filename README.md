**[WIP]**

# MML

Alternative frontend for the Go and JS environments.

- [Documentation](https://github.com/aryszka/mml/blob/master/manual.md)
- [Example](https://github.com/aryszka/mml/blob/master/main.mml)

## Boot the compiler

Prerequisits: Go installed, $GOPATH set and $GOPATH/bin added to $PATH.

```
make deps boot
```

Test:

```
mkdir -p hello
echo 'stdout("Hello, world!\n")' > hello/hello.mml
mml hello/hello > hello/hello.go
go run hello/hello.go
```
