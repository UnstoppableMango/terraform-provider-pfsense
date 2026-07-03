package main

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"
)

func ParseSchema(path string) (*spec.Specification, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s, err := spec.Parse(context.Background(), data)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
