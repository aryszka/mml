SOURCES = $(shell find . -name '*.go')

default: build

build: $(SOURCES)
	go build

check: build
	go test -race

shortcheck: build
	go test -test.short -run ^Test

fmt: $(SOURCES)
	gofmt -w -s ./*.go

precommit: build check fmt
