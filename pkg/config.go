package pkg

import (
	"context"
	"io"
	"io/fs"

	"charm.land/log/v2"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/unmango/go/world"
	"github.com/unstoppablemango/terraform-provider-pfsense/internal/config"
	"gopkg.in/yaml.v3"
)

func GenerateConfig(ctx context.Context, src, dest string) error {
	log := log.FromContext(ctx)
	os := world.FromContext(ctx).Os()

	log.Info("Opening source", "src", src)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
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

	log.Info("Generating tf generator config")
	cfg, err := ConfigFor(ctx, &model.Model)
	if err != nil {
		return err
	}

	log.Info("Validating config")
	if err = cfg.Validate(); err != nil {
		return err
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	log.Info("Writing generated config", "dest", dest)
	return os.WriteFile(dest, out, 0o644)
}

func ConfigFor(ctx context.Context, doc *v3.Document) (*config.Config, error) {
	cfg := &config.Config{
		Provider: config.Provider{
			Name: "pfsense",
		},
		Resources: map[string]config.Resource{
			"firewall_rule": {
				Create: &config.OpenApiSpecLocation{
					Path:   "/api/v2/firewall/rule",
					Method: "POST",
				},
				Read: &config.OpenApiSpecLocation{
					Path:   "/api/v2/firewall/rule",
					Method: "GET",
				},
			},
		},
	}

	return cfg, nil
}
