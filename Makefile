SOURCES = $(shell find . -name '*.go')

default: build

imports:
	@goimports -w $(SOURCES)

build: $(SOURCES)
	go build ./...

install: $(SOURCES)
	go install ./cmd/mml

check: build shortcheck
	# no race here
	# go test ./... -race

shortcheck: build
	go test ./... -test.short -run ^Test

fmt: $(SOURCES)
	@gofmt -w -s $(SOURCES)

precommit: build shortcheck fmt
