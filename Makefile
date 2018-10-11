SOURCES = $(shell find . -name "*.go") syntax.treerack

.PHONY: recompile

build: $(SOURCES)
	go build ./...

install: $(SOURCES)
	go install ./...

deps:
	go get github.com/aryszka/treerack/...
	make gen-parser
	go get ./...

recompile:
	mkdir -p build
	mml main.mml > build/main.1.go
	go run build/main.1.go main.mml > build/main.2.go
	go run build/main.2.go main.mml > build/main.3.go
	rm build/main.1.go build/main.2.go
	mv build/main.3.go cmd/mml/main.go
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
