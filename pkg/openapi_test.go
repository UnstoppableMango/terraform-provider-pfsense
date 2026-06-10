package pkg

import (
	"io"
	"testing"

	"charm.land/log/v2"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	orderedmap "github.com/pb33f/libopenapi/orderedmap"
)

func TestFlattenAllOfs(t *testing.T) {
	// property schemas
	aSchema := &highbase.Schema{}
	aProxy := highbase.CreateSchemaProxy(aSchema)
	bSchema := &highbase.Schema{}
	bProxy := highbase.CreateSchemaProxy(bSchema)

	// base schema with property "a"
	base := &highbase.Schema{}
	base.Properties = orderedmap.New[string, *highbase.SchemaProxy]()
	base.Properties.Set("a", aProxy)
	base.Required = []string{"a"}

	// other schema with property "b"
	other := &highbase.Schema{}
	other.Properties = orderedmap.New[string, *highbase.SchemaProxy]()
	other.Properties.Set("b", bProxy)
	other.Required = []string{"b"}

	// parent schema that allOfs base and other
	parent := &highbase.Schema{}
	parent.AllOf = []*highbase.SchemaProxy{highbase.CreateSchemaProxy(base), highbase.CreateSchemaProxy(other)}

	// build document
	doc := &v3.Document{}
	doc.Components = &v3.Components{Schemas: orderedmap.New[string, *highbase.SchemaProxy]()}
	doc.Components.Schemas.Set("Parent", highbase.CreateSchemaProxy(parent))

	flattenAllOfs(doc, log.New(io.Discard))

	// find merged schema
	var merged *highbase.Schema
	for pair := doc.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		if pair.Key() == "Parent" {
			if pair.Value() == nil {
				t.Fatalf("schema proxy is nil")
			}
			merged = pair.Value().Schema()
			break
		}
	}
	if merged == nil {
		t.Fatalf("merged schema not found")
	}

	// check properties
	count := 0
	if merged.Properties != nil {
		for p := merged.Properties.First(); p != nil; p = p.Next() {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("expected 2 properties after flattening, got %d", count)
	}

	// check required
	existsA, existsB := false, false
	for _, r := range merged.Required {
		if r == "a" {
			existsA = true
		}
		if r == "b" {
			existsB = true
		}
	}
	if !existsA || !existsB {
		t.Fatalf("expected required to contain both 'a' and 'b', got %v", merged.Required)
	}
}

func TestFlattenNestedAllOfs(t *testing.T) {
	// property schemas
	aSchema := &highbase.Schema{}
	aProxy := highbase.CreateSchemaProxy(aSchema)
	bSchema := &highbase.Schema{}
	bProxy := highbase.CreateSchemaProxy(bSchema)
	cSchema := &highbase.Schema{}
	cProxy := highbase.CreateSchemaProxy(cSchema)

	// base schema with property "a"
	base := &highbase.Schema{}
	base.Properties = orderedmap.New[string, *highbase.SchemaProxy]()
	base.Properties.Set("a", aProxy)
	base.Required = []string{"a"}

	// other schema with property "b"
	other := &highbase.Schema{}
	other.Properties = orderedmap.New[string, *highbase.SchemaProxy]()
	other.Properties.Set("b", bProxy)
	other.Required = []string{"b"}

	// inner schema that allOfs base and other
	inner := &highbase.Schema{}
	inner.AllOf = []*highbase.SchemaProxy{highbase.CreateSchemaProxy(base), highbase.CreateSchemaProxy(other)}

	// schema with property "c"
	third := &highbase.Schema{}
	third.Properties = orderedmap.New[string, *highbase.SchemaProxy]()
	third.Properties.Set("c", cProxy)
	third.Required = []string{"c"}

	// parent schema that allOfs inner and third (nested allOf)
	parent := &highbase.Schema{}
	parent.AllOf = []*highbase.SchemaProxy{highbase.CreateSchemaProxy(inner), highbase.CreateSchemaProxy(third)}

	// build document
	doc := &v3.Document{}
	doc.Components = &v3.Components{Schemas: orderedmap.New[string, *highbase.SchemaProxy]()}
	doc.Components.Schemas.Set("ParentNested", highbase.CreateSchemaProxy(parent))

	flattenAllOfs(doc, log.New(io.Discard))

	// find merged schema
	var merged *highbase.Schema
	for pair := doc.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		if pair.Key() == "ParentNested" {
			if pair.Value() == nil {
				t.Fatalf("schema proxy is nil")
			}
			merged = pair.Value().Schema()
			break
		}
	}
	if merged == nil {
		t.Fatalf("merged schema not found")
	}

	// check properties (expect a,b,c)
	count := 0
	if merged.Properties != nil {
		for p := merged.Properties.First(); p != nil; p = p.Next() {
			count++
		}
	}
	if count != 3 {
		t.Fatalf("expected 3 properties after flattening nested allOfs, got %d", count)
	}

	// check required
	existsA, existsB, existsC := false, false, false
	for _, r := range merged.Required {
		if r == "a" {
			existsA = true
		}
		if r == "b" {
			existsB = true
		}
		if r == "c" {
			existsC = true
		}
	}
	if !existsA || !existsB || !existsC {
		t.Fatalf("expected required to contain 'a','b','c', got %v", merged.Required)
	}
}
