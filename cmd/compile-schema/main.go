// Command compile-schema compiles a Starlark schema directory to graph.json.
// This is a helper for CI workflows that need to build graph.json from schema source.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/project-kessel/schema-graph-tools/internal/compiler"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <schema-directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCompiles Starlark schema files to graph.json (written to stdout).\n")
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	schemaDir := flag.Arg(0)
	files, err := readSchemaFiles(schemaDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading schema files: %v\n", err)
		os.Exit(1)
	}

	graph, err := compiler.CompileGraph(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error compiling schema: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(graph)
}

func readSchemaFiles(dir string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".star" {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = contents
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no .star files found in %s", dir)
	}

	return files, nil
}
