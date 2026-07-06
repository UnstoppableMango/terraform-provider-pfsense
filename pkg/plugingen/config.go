package plugingen

import (
	"io"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/unstoppablemango/terraform-provider-pfsense/internal/config"
	"go.yaml.in/yaml/v4"
)

func ConfigFor(doc *v3.Document) (*config.Config, error) {
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

func WriteConfig(cfg *config.Config, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()
	return enc.Encode(cfg)
}
