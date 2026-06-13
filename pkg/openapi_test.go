package pkg

import (
	"testing"

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

	// parent schema that allOfs base and other (inline in path)
	parent := &highbase.Schema{}
	parent.AllOf = []*highbase.SchemaProxy{
		highbase.CreateSchemaProxy(base),
		highbase.CreateSchemaProxy(other),
	}

	// build operation with request body that uses parent as schema
	op := &v3.Operation{}
	rb := &v3.RequestBody{}
	rb.Content = orderedmap.New[string, *v3.MediaType]()
	media := &v3.MediaType{Schema: highbase.CreateSchemaProxy(parent)}
	rb.Content.Set("application/json", media)
	op.RequestBody = rb

	pi := &v3.PathItem{}
	pi.Post = op

	// build document with path
	doc := &v3.Document{}
	doc.Paths = &v3.Paths{PathItems: orderedmap.New[string, *v3.PathItem]()}
	doc.Paths.PathItems.Set("/test", pi)

	// run flattening
	flattenAllOfs(doc)

	// retrieve merged schema
	p := doc.Paths.PathItems.GetOrZero("/test")
	if p == nil || p.Post == nil || p.Post.RequestBody == nil {
		t.Fatalf("constructed path or operation missing")
	}
	m := p.Post.RequestBody.Content.GetOrZero("application/json")
	if m == nil || m.Schema == nil {
		t.Fatalf("media type/schema missing")
	}
	merged := m.Schema.Schema()
	if merged == nil {
		t.Fatalf("merged schema not found")
	}

	// check properties and required
	count := 0
	if merged.Properties != nil {
		for p := merged.Properties.First(); p != nil; p = p.Next() {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("expected 2 properties after flattening, got %d", count)
	}
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
