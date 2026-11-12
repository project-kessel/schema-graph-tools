# Schema Graph Tools

Graph-theory tools for the Kessel Starlark schema: visualization, analysis, and an in-browser playground.

This repo consumes the [starlark-unified-schema](https://github.com/project-kessel/starlark-unified-schema) interpreter's public API (`interpreter/compile`) to build canonical graph representations and tooling around them.

**📖 [Visit the interactive playground](https://project-kessel.github.io/schema-graph-tools/)**

## Tools

- **graph-mermaid**: Renders Mermaid flowcharts from `graph.json`
- **graph-analyze**: Structural analysis (islands, check cost, reachability)
- **graph-playground**: In-browser Starlark → graph WASM playground

## Quick Start

```bash
# Run tests (auto-downloads schema on first run)
make test

# Build all graph tools
make build-graph-mermaid build-graph-analyze build-graph-playground

# Build and serve the interactive playground
make serve-graph-playground
```

## Development

All make targets automatically download the schema from [starlark-unified-schema](https://github.com/project-kessel/starlark-unified-schema) if needed. The schema is cached to `.cache/starlark-unified-schema/` (gitignored).

```bash
make test                # Run all tests (downloads schema if missing)
make test-integration    # Run full suite with integration tests
make refresh-schema      # Update cached schema from GitHub main
make sync-deps         # Synchronize go.mod/go.sum with the cached schema
make clean-schema        # Remove cached schema
```

Run `make refresh-schema` followed by `make sync-deps` to synchronize with the
latest upstream schema. Synchronization runs `go mod tidy`, raising our Go minimum
when required and updating dependency requirements and checksums. Use a sufficiently
recent Go installation or enable automatic toolchain downloads with
`GOTOOLCHAIN=auto make sync-deps`.

The scheduled playground rebuild always fetches upstream `main` and uses the latest
stable Go release. It runs `make sync-deps` in its disposable checkout before
building, so upstream Go minimum and dependency changes do not require a
committed version bump here. Upstream API or schema incompatibilities can still
require code changes.

For local development against a custom schema version:
```bash
go mod edit -replace=github.com/project-kessel/starlark-unified-schema=../starlark-unified-schema/interpreter
```

## Architecture

- `internal/graphvisitor/` — `compile.SchemaVisitor` implementation that builds `graph.json`
- `internal/compiler/` — Public API wrapper around the unified schema compiler
- `internal/graphdoc/` — `graph.json` read model
- `internal/render/` — Mermaid renderer
- `internal/analyze/` — Graph analysis (islands, check cost, reachability)
- `internal/web/` — Browser playground view model
- `cmd/graph-*` — CLI entry points

See [GRAPH.md](GRAPH.md) for the `graph.json` contract.
