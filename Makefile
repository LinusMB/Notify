PREFIX ?= $(HOME)/.local

BIN = notify

$(BIN): $(wildcard cmd/notify/*go)
	go build ./cmd/notify

.PHONY: build
build: $(BIN)

.PHONY: install
install: $(BIN)
	install -D $(BIN) $(PREFIX)/bin
