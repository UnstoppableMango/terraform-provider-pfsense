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

	model, err := doc.BuildV3Model()
	if err != nil {
		return err
	}

	if m := model.Model; m.Components != nil && m.Components.Schemas != nil {
		flattenAllOfs(&m, log)
	}

	bundled, err := bundler.BundleDocument(&model.Model)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, bundled, fs.ModePerm)
}

func flattenAllOfs(doc *v3.Document, log *log.Logger) {
	for pair := doc.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		key := pair.Key()
		proxy := pair.Value()
		if proxy == nil {
			continue
		}

		schema := proxy.Schema()
		if schema == nil || len(schema.AllOf) == 0 {
			continue
		}

		log.Info("Flattening", "schema", key)
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
			mergeAllOf(ap, merged)
		}

		newProxy := highbase.CreateSchemaProxy(merged)
		doc.Components.Schemas.Set(key, newProxy)
	}
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
