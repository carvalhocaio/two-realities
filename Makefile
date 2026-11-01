GOROOT := $(shell go env GOROOT)
PORT ?= 8080

.PHONY: run run-dev test lint build-wasm serve

run:
	go run ./cmd/game

run-dev:
	go run -tags dev ./cmd/game

test:
	go test ./...

lint:
	go vet ./...
	go vet -tags dev ./...
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi

build-wasm:
	GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/game
	cp "$(GOROOT)/lib/wasm/wasm_exec.js" web/

serve:
	python3 -m http.server $(PORT) -d web
