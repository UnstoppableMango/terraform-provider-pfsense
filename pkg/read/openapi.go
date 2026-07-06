package read

import (
	"io"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func OpenAPIV3(r io.Reader) (*v3.Document, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	doc, err := libopenapi.NewDocument(data)
	if err != nil {
		return nil, err
	}

	m, err := doc.BuildV3Model()
	if err != nil {
		return nil, err
	}

	return &m.Model, nil
}
