# Schema Graph Tools

Graph-theory tools for the Kessel Starlark schema: visualization, analysis, and an in-browser playground.

This repo consumes the [starlark-unified-schema](https://github.com/project-kessel/starlark-unified-schema) interpreter's public API (`interpreter/compile`) to build canonical graph representations and tooling around them.

## Tools

- **graph-mermaid**: Renders Mermaid flowcharts from `graph.json`
- **graph-analyze**: Structural analysis (islands, check cost, reachability)
- **graph-playground**: In-browser Starlark → graph WASM playground

## Quick Start

```bash
# Build all graph tools
make build-graph-mermaid build-graph-analyze build-graph-playground

# Run tests (requires starlark-unified-schema checked out as a sibling)
make test

# Build and serve the interactive playground
make serve-graph-playground
```

## Development

This repo tracks `starlark-unified-schema@main` via a `go.mod` replace directive pointing to a local checkout. Tests and builds expect the schema source to be available at `../starlark-unified-schema/schema` by default (override with `SCHEMA_DIR`).

## Architecture

- `internal/graphvisitor/` — `compile.SchemaVisitor` implementation that builds `graph.json`
- `internal/compiler/` — Public API wrapper around the unified schema compiler
- `internal/graphdoc/` — `graph.json` read model
- `internal/render/` — Mermaid renderer
- `internal/analyze/` — Graph analysis (islands, check cost, reachability)
- `internal/web/` — Browser playground view model
- `cmd/graph-*` — CLI entry points

See [GRAPH.md](GRAPH.md) for the `graph.json` contract.
