package compiler_test

// Worked example of driving the starlark-unified-schema public compile API with
// a custom visitor — the same pattern internal/compiler uses to build the graph.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/project-kessel/starlark-unified-schema/compile"
)

// simpleVisitor collects resource names
type simpleVisitor struct {
	resources []string
	entries   []compile.OutputEntry
}

func (v *simpleVisitor) BeginType(name string) {
	v.resources = append(v.resources, name)
}

func (v *simpleVisitor) VisitResource(typeName string, reporter string, commonMembers *compile.Members, reporterMembers *compile.Members, extendsResource *compile.ResourceTypeReference) error {
	return nil
}

func (v *simpleVisitor) VisitDataField(name string, required bool, description *string, dataType any) any {
	return nil
}

func (v *simpleVisitor) VisitTextDataType(minLength *int, maxLength *int, regex *string) any {
	return nil
}
func (v *simpleVisitor) VisitUUIDDataType() any                                      { return nil }
func (v *simpleVisitor) VisitNumericIDDataType(min *int, max *int) any               { return nil }
func (v *simpleVisitor) VisitBooleanDataType() any                                   { return nil }
func (v *simpleVisitor) VisitDateTimeDataType() any                                  { return nil }
func (v *simpleVisitor) VisitEnumDataType(values []string) any                       { return nil }
func (v *simpleVisitor) VisitNullableDataType(inner any) any                         { return nil }
func (v *simpleVisitor) VisitCompositeDataType(dataTypes []any) any                  { return nil }
func (v *simpleVisitor) VisitArrayDataType(items any) any                            { return nil }
func (v *simpleVisitor) VisitObjectDataType(properties []any, required []string) any { return nil }
func (v *simpleVisitor) VisitAnd(left any, right any) any                            { return nil }
func (v *simpleVisitor) VisitOr(left any, right any) any                             { return nil }
func (v *simpleVisitor) VisitUnless(left any, right any) any                         { return nil }
func (v *simpleVisitor) VisitReferenceExpression(name string) any                    { return nil }
func (v *simpleVisitor) VisitSubReferenceExpression(name string, sub string) any     { return nil }
func (v *simpleVisitor) VisitRelation(name string, reporter string, typeName string, cardinality string, idType any) any {
	return nil
}
func (v *simpleVisitor) BeginPermission(name string)               {}
func (v *simpleVisitor) VisitPermission(name string, body any) any { return nil }
func (v *simpleVisitor) Results() ([]compile.OutputEntry, error)   { return v.entries, nil }

func Example() {
	// Read schema files from disk (in production, these could come from anywhere).
	// SCHEMA_DIR is set by make test (points to downloaded cache).
	schemaDir := os.Getenv("SCHEMA_DIR")
	if schemaDir == "" {
		schemaDir = "../../../.cache/starlark-unified-schema/schema"
	}

	files := make(map[string][]byte)

	err := filepath.WalkDir(schemaDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".star" {
			return err
		}
		rel, _ := filepath.Rel(schemaDir, path)
		data, _ := os.ReadFile(path)
		files[filepath.ToSlash(rel)] = data
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a custom visitor
	visitor := &simpleVisitor{}

	// Compile the schema
	if err := compile.Compile(files, visitor); err != nil {
		log.Fatal(err)
	}

	// Access visitor results
	fmt.Println("Schema compilation successful")
	// Output: Schema compilation successful
}
