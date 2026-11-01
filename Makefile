GOROOT := $(shell go env GOROOT)
PORT ?= 8080

.PHONY: run test lint build-wasm serve

run:
	go run ./cmd/game

test:
	go test ./...

lint:
	go vet ./...
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi

build-wasm:
	GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/game
	cp "$(GOROOT)/lib/wasm/wasm_exec.js" web/

serve:
	python3 -m http.server $(PORT) -d web
