SOURCES = $(shell find . -name "*.go") syntax.treerack

.PHONY: boot

build:
	go build ./...

install:
	go install ./...

deps:
	@go get ./...

boot:
	go install ./cmd/runmml
	mkdir -p boot
	runmml compile.mml < compile.mml > boot/compile.go
	go run boot/compile.go < compile.mml > boot/compile.1.go
	go run boot/compile.1.go < compile.mml > cmd/mml/main.go
	rm boot/compile.1.go
	go install ./cmd/mml

check-syntax: syntax.treerack
	treerack check-syntax syntax.treerack

parser/parser.go: check-syntax
	@mkdir -p parser
	treerack generate \
		-export \
		-package-name parser \
		-syntax syntax.treerack \
		> parser/parser.go
	go fmt ./parser

gen-parser: parser/parser.go

check: $(SOURCES) gen-parser
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf boot
