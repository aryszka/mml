SOURCES = $(shell find . -name "*.go") syntax.treerack

deps:
	go get ./...

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
