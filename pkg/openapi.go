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

	bundled, err := bundler.BundleDocument(&model.Model)
	if err != nil {
		return err
	}

	// TODO: flatten allOf's

	return os.WriteFile(dest, bundled, fs.ModePerm)
}
