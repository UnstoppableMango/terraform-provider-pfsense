package main

import (
	"context"
	"io"
	"log/slog"
	"math"
	"strconv"

	"charm.land/log/v2"
	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/bundler"
	"github.com/pb33f/libopenapi/datamodel"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/unmango/go/world"
)

func PatchSpec(ctx context.Context, src, dest string) error {
	log := log.FromContext(ctx)
	os := world.FromContext(ctx).Os()

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	data, err := io.ReadAll(s)
	if err != nil {
		return err
	}

	cfg := &datamodel.DocumentConfiguration{
		Logger: slog.New(log),
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(data, cfg)
	if err != nil {
		return err
	}

	dropUnrepresentableIntBounds(doc.GetSpecInfo().RootNode)

	model, err := doc.BuildV3Model()
	if err != nil {
		return err
	}

	// flatten top-level component schemas and all other schema locations in the document
	flattenAllOfs(&model.Model)

	bundled, err := bundler.BundleDocument(&model.Model)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, bundled, 0o644)
}

// dropUnrepresentableIntBounds removes the minimum and maximum of integer schemas
// when the bound, read as a float64, is outside the int64 range. The generator
// converts bounds with int64(float64), whose result for those values differs by
// architecture: pfrest's maximum of 2^63-1 rounds up to 2^63 and becomes
// MaxInt64 on arm64 but MinInt64 on amd64. A maximum of 2^63-1 constrains
// nothing an int64 attribute could hold, so dropping it loses no validation.
// Maximums from 2^63-512 to 2^63-2 also round to 2^63 and are dropped, losing
// a bound they would have enforced; pfrest uses none of them.
func dropUnrepresentableIntBounds(node *yaml.Node) {
	if node == nil {
		return
	}
	if node.Kind == yaml.MappingNode && isIntegerSchema(node) {
		kept := node.Content[:0]
		for i := 0; i+1 < len(node.Content); i += 2 {
			key, value := node.Content[i], node.Content[i+1]
			if (key.Value == "minimum" || key.Value == "maximum") && !fitsInt64(value.Value) {
				continue
			}
			kept = append(kept, key, value)
		}
		node.Content = kept
	}
	for _, child := range node.Content {
		dropUnrepresentableIntBounds(child)
	}
}

func isIntegerSchema(node *yaml.Node) bool {
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != "type" {
			continue
		}
		// OpenAPI 3.1 allows a list of types, such as [integer, "null"]
		types := node.Content[i+1]
		if types.Value == "integer" {
			return true
		}
		for _, t := range types.Content {
			if t.Value == "integer" {
				return true
			}
		}
	}
	return false
}

func fitsInt64(value string) bool {
	// A literal beyond float64's range parses to ±Inf alongside ErrRange, and a
	// value that is not a number parses to 0, which fits.
	f, _ := strconv.ParseFloat(value, 64)
	return f >= math.MinInt64 && f < math.MaxInt64
}

func mergeAllOf(proxy *highbase.SchemaProxy, target *highbase.Schema) {
	if proxy == nil {
		return
	}

	schema := proxy.Schema()
	if schema == nil {
		return
	}

	// If this schema itself has allOf entries, merge them first (recursive)
	if len(schema.AllOf) > 0 {
		for _, sp := range schema.AllOf {
			mergeAllOf(sp, target)
		}
	}

	if schema.Properties != nil {
		for p := schema.Properties.First(); p != nil; p = p.Next() {
			target.Properties.Set(p.Key(), p.Value())
		}
	}

	// merge required
	if len(schema.Required) > 0 {
		exists := map[string]struct{}{}
		for _, r := range target.Required {
			exists[r] = struct{}{}
		}
		for _, r := range schema.Required {
			if _, ok := exists[r]; !ok {
				target.Required = append(target.Required, r)
			}
		}
	}
}

// flattenProxy checks a SchemaProxy and returns a new SchemaProxy if it had allOf entries and was flattened.
func flattenProxy(proxy *highbase.SchemaProxy) *highbase.SchemaProxy {
	if proxy == nil {
		return nil
	}
	schema := proxy.Schema()
	if schema == nil || len(schema.AllOf) == 0 {
		return nil
	}
	merged := &highbase.Schema{}
	*merged = *schema
	merged.AllOf = nil
	if merged.Properties == nil {
		merged.Properties = orderedmap.New[string, *highbase.SchemaProxy]()
	}
	for _, ap := range schema.AllOf {
		mergeAllOf(ap, merged)
	}
	return highbase.CreateSchemaProxy(merged)
}

func flattenMediaType(m *v3.MediaType) {
	if m == nil {
		return
	}
	if newp := flattenProxy(m.Schema); newp != nil {
		m.Schema = newp
	}
	if newp := flattenProxy(m.ItemSchema); newp != nil {
		m.ItemSchema = newp
	}
}

func flattenHeader(h *v3.Header) {
	if h == nil {
		return
	}
	if newp := flattenProxy(h.Schema); newp != nil {
		h.Schema = newp
	}
	if h.Content != nil {
		for pair := h.Content.First(); pair != nil; pair = pair.Next() {
			flattenMediaType(pair.Value())
		}
	}
}

func flattenRequestBody(rb *v3.RequestBody) {
	if rb == nil || rb.Content == nil {
		return
	}
	for pair := rb.Content.First(); pair != nil; pair = pair.Next() {
		flattenMediaType(pair.Value())
	}
}

func flattenParameter(p *v3.Parameter) {
	if p == nil {
		return
	}
	if newp := flattenProxy(p.Schema); newp != nil {
		p.Schema = newp
	}
	if p.Content != nil {
		for pair := p.Content.First(); pair != nil; pair = pair.Next() {
			flattenMediaType(pair.Value())
		}
	}
}

func flattenResponse(r *v3.Response) {
	if r == nil {
		return
	}
	if r.Headers != nil {
		for pair := r.Headers.First(); pair != nil; pair = pair.Next() {
			flattenHeader(pair.Value())
		}
	}
	if r.Content != nil {
		for pair := r.Content.First(); pair != nil; pair = pair.Next() {
			flattenMediaType(pair.Value())
		}
	}
}

// flattenAllOfs walks the document and flattens allOf entries across components and paths.
func flattenAllOfs(doc *v3.Document) {
	if doc == nil {
		return
	}
	// components
	if doc.Components != nil {
		if doc.Components.Schemas != nil {
			for pair := doc.Components.Schemas.First(); pair != nil; pair = pair.Next() {
				if pair.Value() == nil {
					continue
				}
				if newp := flattenProxy(pair.Value()); newp != nil {
					doc.Components.Schemas.Set(pair.Key(), newp)
				}
			}
		}
		if doc.Components.RequestBodies != nil {
			for pair := doc.Components.RequestBodies.First(); pair != nil; pair = pair.Next() {
				flattenRequestBody(pair.Value())
			}
		}
		if doc.Components.Responses != nil {
			for pair := doc.Components.Responses.First(); pair != nil; pair = pair.Next() {
				flattenResponse(pair.Value())
			}
		}
		if doc.Components.Parameters != nil {
			for pair := doc.Components.Parameters.First(); pair != nil; pair = pair.Next() {
				flattenParameter(pair.Value())
			}
		}
		if doc.Components.Headers != nil {
			for pair := doc.Components.Headers.First(); pair != nil; pair = pair.Next() {
				flattenHeader(pair.Value())
			}
		}
	}

	// paths
	if doc.Paths != nil && doc.Paths.PathItems != nil {
		for pair := doc.Paths.PathItems.First(); pair != nil; pair = pair.Next() {
			pi := pair.Value()
			if pi == nil {
				continue
			}
			ops := []*v3.Operation{pi.Get, pi.Put, pi.Post, pi.Delete, pi.Options, pi.Head, pi.Patch, pi.Trace, pi.Query}
			for _, op := range ops {
				if op == nil {
					continue
				}
				// parameters
				for _, p := range op.Parameters {
					flattenParameter(p)
				}
				// request body
				flattenRequestBody(op.RequestBody)
				// responses
				if op.Responses != nil {
					if op.Responses.Default != nil {
						flattenResponse(op.Responses.Default)
					}
					if op.Responses.Codes != nil {
						for pair := op.Responses.Codes.First(); pair != nil; pair = pair.Next() {
							flattenResponse(pair.Value())
						}
					}
				}
			}
		}
	}
}
