package pkg

import (
	"io"
	"io/fs"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/bundler"
	"github.com/unmango/go/world"
)

func PatchSpec(os world.Os, src, dest string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	data, err := io.ReadAll(s)
	if err != nil {
		return err
	}

	doc, err := libopenapi.NewDocument(data)
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
