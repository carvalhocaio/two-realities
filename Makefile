GOROOT := $(shell go env GOROOT)
PORT ?= 8080
DIST := dist/html5

.PHONY: run run-dev test test-cover lint build-wasm package serve shots

run:
	go run ./cmd/game

run-dev:
	go run -tags dev ./cmd/game

test:
	go test ./...
	node --test web/*.test.mjs

test-cover:
	go test -cover ./...
	node --test --experimental-test-coverage web/*.test.mjs

lint:
	go vet ./...
	go vet -tags dev ./...
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi

build-wasm:
	GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/game
	cp "$(GOROOT)/lib/wasm/wasm_exec.js" web/

package: build-wasm
	rm -rf $(DIST)
	mkdir -p $(DIST)
	cp web/index.html web/loader.js web/game.wasm web/wasm_exec.js $(DIST)/
	sed -i "s/name=\"wasm-size\" content=\"0\"/name=\"wasm-size\" content=\"$$(wc -c < $(DIST)/game.wasm)\"/" $(DIST)/index.html
	@grep -q 'name="wasm-size" content="[1-9]' $(DIST)/index.html || { echo "wasm size not injected"; exit 1; }
	@raw=$$(wc -c < $(DIST)/game.wasm); gz=$$(gzip -9 -c $(DIST)/game.wasm | wc -c); \
	echo "packaged $(DIST): game.wasm $$((raw / 1048576)).$$((raw % 1048576 * 10 / 1048576)) MB (gzip $$((gz / 1048576)).$$((gz % 1048576 * 10 / 1048576)) MB)"

shots:
	go run -tags dev ./cmd/shots -out $(or $(OUT),shots)

serve:
	python3 -m http.server $(PORT) -d web
