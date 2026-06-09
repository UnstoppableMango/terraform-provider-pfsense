package pkg

import (
	"context"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/unstoppablemango/terraform-provider-pfsense/internal/config"
)

func GenerateConfig(ctx context.Context, doc *v3.Document) (*config.Config, error) {
	cfg := &config.Config{}

	return cfg, nil
}
