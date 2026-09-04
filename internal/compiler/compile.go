package compiler

import (
	"fmt"

	"github.com/project-kessel/schema-graph-tools/internal/graphvisitor"
	"github.com/project-kessel/starlark-unified-schema/compile"
)

// CompileGraph builds the canonical graph.json from an in-memory set of
// Starlark source files. This is the Repo B equivalent of the old
// lang.CompileGraph function - it uses the public compile API to drive
// the GraphVisitor.
//
// The files map should be keyed by path relative to the schema root (e.g.
// "kessel.star", "workspace/reporters/features/workspace.star").
// The caller must include kessel.star in the files map.
func CompileGraph(files map[string][]byte) ([]byte, error) {
	visitor := graphvisitor.NewGraphVisitor()

	if err := compile.Compile(files, visitor); err != nil {
		return nil, err
	}

	results, err := visitor.Results()
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		if r.Path == "graph.json" {
			return r.Contents, nil
		}
	}

	return nil, fmt.Errorf("graph visitor produced no graph.json entry")
}
