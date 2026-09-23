.PHONY: test lint build-graph-mermaid build-graph-analyze build-graph-playground build-graph-wasm build-compile-schema graph-playground serve-graph-playground fetch-schema refresh-schema clean-schema clean

GRAPH_WEB_PORT ?= 8000
GRAPH_OUTPUT_DIR ?= output/graph
GRAPH_PLAYGROUND_DIR ?= output/playground

# Schema downloaded from GitHub and cached locally (override with SCHEMA_DIR if needed)
SCHEMA_DIR ?= .cache/starlark-unified-schema/schema
SCHEMA_REF ?= main

# Go ships wasm_exec.js in lib/wasm (>= 1.24) or misc/wasm (older). Resolve once.
WASM_EXEC := $(shell if [ -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then echo "$$(go env GOROOT)/lib/wasm/wasm_exec.js"; else echo "$$(go env GOROOT)/misc/wasm/wasm_exec.js"; fi)

lint:
	gofmt -w ./cmd ./internal
	go vet ./...

# Download and cache schema files from GitHub (for integration tests)
fetch-schema:
	@./scripts/fetch-schema.sh $(SCHEMA_REF)

# Force refresh: download latest schema from GitHub
refresh-schema:
	@rm -rf .cache/starlark-unified-schema
	@./scripts/fetch-schema.sh $(SCHEMA_REF)
	@echo "Schema refreshed to latest from GitHub"

# Run tests including integration tests (requires downloaded schema)
test-integration: fetch-schema
	SCHEMA_DIR=$$(cd $(SCHEMA_DIR) && pwd) go test -count=1 ./...

# Remove cached schema
clean-schema:
	rm -rf .cache/starlark-unified-schema

test: fetch-schema
	go test -count=1 ./...

build-graph-mermaid: fetch-schema
	go build -o bin/graph-mermaid ./cmd/graph-mermaid

build-graph-analyze: fetch-schema
	go build -o bin/graph-analyze ./cmd/graph-analyze

build-graph-playground: fetch-schema
	go build -o bin/graph-playground ./cmd/graph-playground

build-compile-schema: fetch-schema
	go build -o bin/compile-schema ./cmd/compile-schema

# Build the in-browser schema compiler (Go -> WASM). The binary is large and
# toolchain-specific, so it lands in the gitignored output dir, never in git.
build-graph-wasm: fetch-schema
	mkdir -p "$(GRAPH_PLAYGROUND_DIR)"
	GOOS=js GOARCH=wasm go build -o "$(GRAPH_PLAYGROUND_DIR)/graph-playground.wasm" ./cmd/graph-wasm

# Assemble the live playground site (schema source + WASM compiler +
# wasm_exec.js) into $(GRAPH_PLAYGROUND_DIR). This is the deployable artifact —
# the CI Pages workflow uploads exactly this directory. Building it is separate
# from serving so both the workflow and serve-graph-playground share one recipe.
graph-playground: build-graph-playground build-graph-wasm fetch-schema
	cp "$(WASM_EXEC)" "$(GRAPH_PLAYGROUND_DIR)/wasm_exec.js"
	./bin/graph-playground -src "$(SCHEMA_DIR)" -out "$(GRAPH_PLAYGROUND_DIR)/index.html"

# Serve the built playground over http. The page compiles Starlark to the graph
# entirely in the browser; the .wasm sidecar must be fetched (blocked over
# file://), so this must be served over http rather than opened as a file.
serve-graph-playground: graph-playground
	@echo "Serving $(GRAPH_PLAYGROUND_DIR) at http://localhost:$(GRAPH_WEB_PORT)/ (Ctrl-C to stop)"
	cd "$(GRAPH_PLAYGROUND_DIR)" && python3 -m http.server $(GRAPH_WEB_PORT) --bind 127.0.0.1

clean:
	rm -rf bin/
	rm -rf output/
	rm -rf .cache/
