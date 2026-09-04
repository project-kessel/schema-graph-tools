.PHONY: test lint build-graph-mermaid build-graph-analyze build-graph-playground build-graph-wasm build-compile-schema graph graph-analyze graph-playground serve-graph-playground clean

GRAPH_WEB_PORT ?= 8000
GRAPH_OUTPUT_DIR ?= output/graph
GRAPH_PLAYGROUND_DIR ?= output/playground

# Schema lives in Repo A (starlark-unified-schema). For local dev, point to a
# sibling checkout. CI will check out Repo A and set this to the checkout path.
SCHEMA_DIR ?= ../starlark-unified-schema/schema
SCHEMA_CHECKOUT_DIR ?= ../starlark-unified-schema

# Go ships wasm_exec.js in lib/wasm (>= 1.24) or misc/wasm (older). Resolve once.
WASM_EXEC := $(shell if [ -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then echo "$$(go env GOROOT)/lib/wasm/wasm_exec.js"; else echo "$$(go env GOROOT)/misc/wasm/wasm_exec.js"; fi)

lint:
	gofmt -w ./cmd ./internal
	go vet ./...

test:
	@test -d "$(SCHEMA_DIR)" || { echo "error: schema directory $(SCHEMA_DIR) not found. Clone the sibling repo: git clone https://github.com/project-kessel/starlark-unified-schema ../" >&2; exit 1; }
	SCHEMA_DIR=$$(cd $(SCHEMA_DIR) && pwd) go test -count=1 ./...

build-graph-mermaid:
	go build -o bin/graph-mermaid ./cmd/graph-mermaid

build-graph-analyze:
	go build -o bin/graph-analyze ./cmd/graph-analyze

build-graph-playground:
	go build -o bin/graph-playground ./cmd/graph-playground

build-compile-schema:
	go build -o bin/compile-schema ./cmd/compile-schema

# Build the in-browser schema compiler (Go -> WASM). The binary is large and
# toolchain-specific, so it lands in the gitignored output dir, never in git.
build-graph-wasm:
	mkdir -p "$(GRAPH_PLAYGROUND_DIR)"
	GOOS=js GOARCH=wasm go build -o "$(GRAPH_PLAYGROUND_DIR)/graph-playground.wasm" ./cmd/graph-wasm

# Render Mermaid diagram from an existing graph.json
graph: build-graph-mermaid
	@test -f "$(GRAPH_OUTPUT_DIR)/graph.json" || { echo "error: $(GRAPH_OUTPUT_DIR)/graph.json not found" >&2; exit 1; }
	./bin/graph-mermaid -in "$(GRAPH_OUTPUT_DIR)/graph.json" -out "$(GRAPH_OUTPUT_DIR)/graph.mmd"

# Analyze graph.json for structural problems (islands / isolated resources).
graph-analyze: build-graph-analyze
	@test -f "$(GRAPH_OUTPUT_DIR)/graph.json" || { echo "error: $(GRAPH_OUTPUT_DIR)/graph.json not found" >&2; exit 1; }
	./bin/graph-analyze -in "$(GRAPH_OUTPUT_DIR)/graph.json"

# Assemble the live playground site (schema source + WASM compiler +
# wasm_exec.js) into $(GRAPH_PLAYGROUND_DIR). This is the deployable artifact —
# the CI Pages workflow uploads exactly this directory. Building it is separate
# from serving so both the workflow and serve-graph-playground share one recipe.
graph-playground: build-graph-playground build-graph-wasm
	@test -d "$(SCHEMA_DIR)" || { echo "error: schema directory $(SCHEMA_DIR) not found" >&2; exit 1; }
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
