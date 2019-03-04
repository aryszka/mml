.PHONY: recompile boot

default: recompile

deps:
	go get github.com/aryszka/treerack/...
	make gen-parser
	go get ./...

boot:
	go install ./boot/mml

builddir:
	mkdir -p build

compile-proto: builddir
	mml main > build/main.1.go

compile-new:
	go run build/main.1.go main > build/main.2.go
	go run build/main.2.go main > build/main.3.go
	diff build/main.2.go build/main.3.go
	rm build/main.1.go build/main.2.go
	mv build/main.3.go boot/mml/main.go
	# in order to avoid unnecessary diffs:
	go fmt boot/mml/main.go
	go install ./boot/mml

recompile: compile-proto compile-new

check: check-syntax

check-syntax: parser.treerack
	treerack check-syntax parser.treerack

parser/parser.go: check-syntax
	mkdir -p parser
	treerack generate \
		-export \
		-package-name parser \
		-syntax parser.treerack \
		> parser/parser.go
	# in order to avoid unnecessary diffs:
	go fmt ./parser

gen-parser: parser/parser.go

check-syntax2: parser2.treerack
	treerack check-syntax parser2.treerack

check-syntax-test: parsertest.treerack
	treerack check-syntax parsertest.treerack

gen-parser2: check-syntax2
	mkdir -p parser
	treerack generate \
		-export \
		-package-name parser \
		-syntax parser2.treerack \
		> parser/parser.go
	# in order to avoid unnecessary diffs:
	go fmt ./parser

gen-parser-test: check-syntax-test
	mkdir -p parser
	treerack generate \
		-export \
		-package-name parser \
		-syntax parsertest.treerack \
		> parser/parser.go
	# in order to avoid unnecessary diffs:
	go fmt ./parser

fmt:
	go fmt builtin.go

clean:
	rm -rf build
