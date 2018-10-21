.PHONY: recompile boot

default: recompile

deps:
	go get github.com/aryszka/treerack/...
	make gen-parser
	go get ./...

boot:
	go install ./boot/mml

recompile:
	mkdir -p build
	mml main.mml > build/main.1.go
	go run build/main.1.go main.mml > build/main.2.go
	go run build/main.2.go main.mml > build/main.3.go
	diff build/main.2.go build/main.3.go
	rm build/main.1.go build/main.2.go
	mv build/main.3.go boot/mml/main.go
	# in order to avoid unnecessary diffs:
	go fmt builtin.go
	go install ./boot/mml

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

clean:
	rm -rf build
