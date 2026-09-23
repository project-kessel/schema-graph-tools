package compiler_test

// These tests exercise the public compile API of starlark-unified-schema
// (Repo A) from the consumer side. They live here rather than in Repo A so that
// repo's review surface stays minimal; the exported contract is still covered
// end-to-end, and a breaking change there fails this suite.

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/project-kessel/starlark-unified-schema/compile"
)

// realSchemaDir resolves the schema directory. SCHEMA_DIR is set by make test
// (pointing to the downloaded cache at .cache/starlark-unified-schema/schema).
func realSchemaDir(t *testing.T) string {
	t.Helper()

	dir := os.Getenv("SCHEMA_DIR")
	if dir == "" {
		dir = "../../../.cache/starlark-unified-schema/schema"
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("schema directory %q not found; run 'make fetch-schema'", dir)
	}

	return dir
}

// readStarFiles loads every .star file under dir, keyed by path relative to dir.
func readStarFiles(t *testing.T, dir string) map[string][]byte {
	t.Helper()

	files := map[string][]byte{}
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
		t.Fatalf("Failed to read schema files: %v", err)
	}

	return files
}

// readKesselStar loads the prelude, skipping the test when it is unavailable.
func readKesselStar(t *testing.T) []byte {
	t.Helper()

	contents, err := os.ReadFile(filepath.Join(realSchemaDir(t), "kessel.star"))
	if err != nil {
		t.Skip("kessel.star not available")
	}

	return contents
}

// mockVisitor implements compile.SchemaVisitor for testing
type mockVisitor struct {
	beginTypeCalls    []string
	visitResourceCall bool
	resultEntries     []compile.OutputEntry
	resultError       error
}

func (m *mockVisitor) BeginType(name string) {
	m.beginTypeCalls = append(m.beginTypeCalls, name)
}

func (m *mockVisitor) VisitResource(typeName string, reporter string, commonMembers *compile.Members, reporterMembers *compile.Members, extendsResource *compile.ResourceTypeReference) error {
	m.visitResourceCall = true
	return nil
}

func (m *mockVisitor) VisitDataField(name string, required bool, description *string, dataType any) any {
	return map[string]any{"name": name, "required": required}
}

func (m *mockVisitor) VisitTextDataType(minLength *int, maxLength *int, regex *string) any {
	return map[string]any{"kind": "text"}
}

func (m *mockVisitor) VisitUUIDDataType() any {
	return map[string]any{"kind": "uuid"}
}

func (m *mockVisitor) VisitNumericIDDataType(min *int, max *int) any {
	return map[string]any{"kind": "numeric_id"}
}

func (m *mockVisitor) VisitBooleanDataType() any {
	return map[string]any{"kind": "boolean"}
}

func (m *mockVisitor) VisitDateTimeDataType() any {
	return map[string]any{"kind": "date_time"}
}

func (m *mockVisitor) VisitEnumDataType(values []string) any {
	return map[string]any{"kind": "enum", "values": values}
}

func (m *mockVisitor) VisitNullableDataType(inner any) any {
	return map[string]any{"kind": "nullable", "inner": inner}
}

func (m *mockVisitor) VisitCompositeDataType(dataTypes []any) any {
	return map[string]any{"kind": "composite", "types": dataTypes}
}

func (m *mockVisitor) VisitArrayDataType(items any) any {
	return map[string]any{"kind": "array", "items": items}
}

func (m *mockVisitor) VisitObjectDataType(properties []any, required []string) any {
	return map[string]any{"kind": "object", "properties": properties}
}

func (m *mockVisitor) VisitAnd(left any, right any) any {
	return map[string]any{"kind": "and", "left": left, "right": right}
}

func (m *mockVisitor) VisitOr(left any, right any) any {
	return map[string]any{"kind": "or", "left": left, "right": right}
}

func (m *mockVisitor) VisitUnless(left any, right any) any {
	return map[string]any{"kind": "unless", "left": left, "right": right}
}

func (m *mockVisitor) VisitReferenceExpression(name string) any {
	return map[string]any{"kind": "reference", "name": name}
}

func (m *mockVisitor) VisitSubReferenceExpression(name string, sub string) any {
	return map[string]any{"kind": "subreference", "name": name, "sub": sub}
}

func (m *mockVisitor) VisitRelation(name string, reporter string, typeName string, cardinality string, idType any) any {
	return map[string]any{"name": name, "typeName": typeName}
}

func (m *mockVisitor) BeginPermission(name string) {}

func (m *mockVisitor) VisitPermission(name string, body any) any {
	return map[string]any{"name": name, "body": body}
}

func (m *mockVisitor) Results() ([]compile.OutputEntry, error) {
	return m.resultEntries, m.resultError
}

// TestCompileRealSchema compiles the actual committed schema to verify the public
// API works end-to-end with real world complexity
func TestCompileRealSchema(t *testing.T) {
	files := readStarFiles(t, realSchemaDir(t))

	visitor := &mockVisitor{}
	if err := compile.Compile(files, visitor); err != nil {
		t.Fatalf("Compile failed on real schema: %v", err)
	}

	// Real schema has multiple resource types
	if len(visitor.beginTypeCalls) < 3 {
		t.Errorf("Expected multiple resource types in real schema, got %d", len(visitor.beginTypeCalls))
	}

	if !visitor.visitResourceCall {
		t.Error("Expected VisitResource to be called")
	}
}

// TestCompileMissingPrelude verifies error handling when kessel.star is missing
func TestCompileMissingPrelude(t *testing.T) {
	files := map[string][]byte{
		"simple.star": []byte(`
load("kessel.star", "resource")
simple_resource = resource(reporter="test")
`),
	}

	visitor := &mockVisitor{}
	if err := compile.Compile(files, visitor); err == nil {
		t.Fatal("Expected error when kessel.star is missing")
	}
}

// TestCompileInvalidStarlark verifies error handling for syntax errors
func TestCompileInvalidStarlark(t *testing.T) {
	files := map[string][]byte{
		"kessel.star":  []byte(`def resource(): pass`),
		"invalid.star": []byte(`this is not valid starlark syntax ][{`),
	}

	visitor := &mockVisitor{}
	if err := compile.Compile(files, visitor); err == nil {
		t.Fatal("Expected error for invalid Starlark syntax")
	}
}

// TestCompileIgnoresNonStarFiles verifies that only .star files are processed
func TestCompileIgnoresNonStarFiles(t *testing.T) {
	files := map[string][]byte{
		"kessel.star": readKesselStar(t),
		"test.star": []byte(`
load("kessel.star", "resource", "uuid")
test_resource = resource(reporter="test", id_type=uuid())
`),
		"README.md":   []byte(`# This should be ignored`),
		"config.json": []byte(`{"ignored": true}`),
		"helper.py":   []byte(`print("not starlark")`),
	}

	visitor := &mockVisitor{}
	if err := compile.Compile(files, visitor); err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Only test.star should produce a resource
	if len(visitor.beginTypeCalls) != 1 {
		t.Errorf("Expected 1 resource, got %d (non-.star files should be ignored)", len(visitor.beginTypeCalls))
	}

	if visitor.beginTypeCalls[0] != "test_resource" {
		t.Errorf("Expected resource 'test_resource', got %s", visitor.beginTypeCalls[0])
	}
}

// TestCompileDeterministicOrder verifies files are processed in sorted order
func TestCompileDeterministicOrder(t *testing.T) {
	files := map[string][]byte{
		"kessel.star": readKesselStar(t),
		"z_last.star": []byte(`
load("kessel.star", "resource", "uuid")
z_resource = resource(reporter="z", id_type=uuid())
`),
		"a_first.star": []byte(`
load("kessel.star", "resource", "uuid")
a_resource = resource(reporter="a", id_type=uuid())
`),
		"m_middle.star": []byte(`
load("kessel.star", "resource", "uuid")
m_resource = resource(reporter="m", id_type=uuid())
`),
	}

	visitor := &mockVisitor{}
	if err := compile.Compile(files, visitor); err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Files should be processed in sorted order
	if len(visitor.beginTypeCalls) != 3 {
		t.Fatalf("Expected 3 resources, got %d", len(visitor.beginTypeCalls))
	}

	expected := []string{"a_resource", "m_resource", "z_resource"}
	for i, name := range expected {
		if visitor.beginTypeCalls[i] != name {
			t.Errorf("Position %d: expected %s, got %s", i, name, visitor.beginTypeCalls[i])
		}
	}
}
