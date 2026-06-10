package pkg

import (
	"context"
	"io"
	"io/fs"
	"log/slog"

	"charm.land/log/v2"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/bundler"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
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

	model, err := doc.BuildV3Model()
	if err != nil {
		return err
	}

	if err := flattenAllOfs(&model.Model); err != nil {
		return err
	}

	bundled, err := bundler.BundleDocument(&model.Model)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, bundled, fs.ModePerm)
}

// flattenAllOfs merges allOf entries for component schemas into single inline schemas
func flattenAllOfs(doc *v3.Document) error {
	if doc == nil || doc.Components == nil || doc.Components.Schemas == nil {
		return nil
	}
	for pair := doc.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		key := pair.Key()
		sp := pair.Value()
		if sp == nil {
			continue
		}
		schema := sp.Schema()
		if schema == nil || len(schema.AllOf) == 0 {
			continue
		}
		// create a shallow copy of the parent schema and clear AllOf
		merged := &highbase.Schema{}
		*merged = *schema
		merged.AllOf = nil
		// ensure properties map exists
		if merged.Properties == nil {
			merged.Properties = orderedmap.New[string, *highbase.SchemaProxy]()
		}
		// merge each allOf schema into merged
		for _, ap := range schema.AllOf {
			if ap == nil {
				continue
			}
			as := ap.Schema()
			if as == nil {
				continue
			}
			if as.Properties != nil {
				for p := as.Properties.First(); p != nil; p = p.Next() {
					merged.Properties.Set(p.Key(), p.Value())
				}
			}
			// merge required
			if len(as.Required) > 0 {
				exists := map[string]struct{}{}
				for _, r := range merged.Required {
					exists[r] = struct{}{}
				}
				for _, r := range as.Required {
					if _, ok := exists[r]; !ok {
						merged.Required = append(merged.Required, r)
					}
				}
			}
		}
		newProxy := highbase.CreateSchemaProxy(merged)
		doc.Components.Schemas.Set(key, newProxy)
	}
	return nil
}
