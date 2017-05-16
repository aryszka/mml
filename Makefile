SOURCES = $(shell find . -name '*.go')
NEXT_SOURCES = $(shell find next -name '*.go')

default: build

next: build-next

imports:
	@goimports -w $(SOURCES)

build: $(SOURCES)
	go build ./...

build-next: $(NEXT_SOURCES)
	go build ./next

install: $(SOURCES)
	go install ./cmd/mml

check: build shortcheck
	# no race here
	# go test ./... -race

check-next: next shortcheck-next
	# no race here
	# go test ./... -race

shortcheck: build
	go test ./... -test.short -run ^Test

shortcheck-next: next
	go test ./next -test.short -run ^Test

fmt: $(SOURCES)
	@gofmt -w -s $(SOURCES)

precommit: build shortcheck fmt
