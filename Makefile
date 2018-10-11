SOURCES = $(shell find . -name "*.go") syntax.treerack

.PHONY: boot

build:
	go build ./...

install:
	go install ./...

deps:
	go get github.com/aryszka/treerack/...
	make gen-parser
	go get ./...

recompile:
	mkdir -p build
	mml compile.mml > build/compile.1.go
	go run build/compile.1.go compile.mml > build/compile.2.go
	go run build/compile.2.go compile.mml > build/compile.3.go
	rm build/compile.1.go build/compile.2.go
	mv build/compile.3.go cmd/mml/main.go
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
