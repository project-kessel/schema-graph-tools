# Schema Graph Tools — Developer Guide

For how to...
- build, test, and run graph tools, see [README.md](README.md)
- understand the graph.json contract and notation, see [GRAPH.md](GRAPH.md)
- understand the architecture and internal packages, see Architecture section below

## Architecture

- **`internal/compiler/`** — Public API wrapper around `starlark-unified-schema/interpreter/compile`. Handles schema loading and compilation into the vistor pipeline.
- **`internal/graphvisitor/`** — `SchemaVisitor` implementation. Walks the compiled schema and builds the canonical `graph.json` representation.
- **`internal/graphdoc/`** — Read model for `graph.json`. Defines the in-memory graph structure (nodes, edges, permissions, types).
- **`internal/render/`** — Mermaid diagram renderer. Converts `graph.json` into flowcharts.
- **`internal/analyze/`** — Graph analysis. Computes permission check cost, reachability, and structural problems (islands, cycles).
- **`internal/web/`** — Browser playground view model. Powers the in-browser Starlark → graph compiler (WASM).
- **`cmd/graph-*`** — CLI entry points for each tool (mermaid, analyze, playground, compile-schema, graph-wasm).

## Build & Test

```bash
# Ensure starlark-unified-schema is checked out as a sibling
git clone https://github.com/project-kessel/starlark-unified-schema ../

# Build all tools
make build-graph-mermaid build-graph-analyze build-graph-playground

# Run tests
make test

# Serve the interactive playground
make serve-graph-playground
```

## Repo Rules

- Commit source code only — not `bin/`, `output/`, or generated artifacts.
- `internal/compiler` wraps a public API (`starlark-unified-schema/interpreter/compile`). Changes to this wrapper (especially the `SchemaVisitor` interface or how the compiler is invoked) may require coordination with downstream consumers.
- Always run `make test` before submitting changes.
- If changes affect the `graph.json` contract, update [GRAPH.md](GRAPH.md) in the same PR.
- The `Makefile` requires `SCHEMA_DIR` to point to a valid `starlark-unified-schema/schema` directory (checked out as a sibling by default). Tests will fail with a clear error if this is missing.
