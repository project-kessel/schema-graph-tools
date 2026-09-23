module github.com/project-kessel/schema-graph-tools

go 1.26.6

require (
	github.com/project-kessel/starlark-unified-schema v0.0.0
	github.com/stretchr/testify v1.11.1
	gonum.org/v1/gonum v0.17.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.starlark.net v0.0.0-20260522144826-ec58d4b459e2 // indirect
	golang.org/x/sys v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Nested module requires replace directive
// Points to downloaded cache (via fetch-schema) or local sibling checkout
replace github.com/project-kessel/starlark-unified-schema => ./.cache/starlark-unified-schema/interpreter
