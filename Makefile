PREFIX ?= $(HOME)/.local

BIN = notify

SRCS = $(shell find cmd/ internal/ -type f -name '*go')

$(BIN): $(SRCS)
	go build ./cmd/notify

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build: $(BIN)

.PHONY: install
install: $(BIN)
	install -D $(BIN) $(PREFIX)/bin
