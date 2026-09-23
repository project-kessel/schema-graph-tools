# Schema Graph Tools

Graph-theory tools for the Kessel Starlark schema: visualization, analysis, and an in-browser playground.

This repo consumes the [starlark-unified-schema](https://github.com/project-kessel/starlark-unified-schema) interpreter's public API (`interpreter/compile`) to build canonical graph representations and tooling around them.

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
make clean-schema        # Remove cached schema
```

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
