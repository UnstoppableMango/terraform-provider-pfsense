package config

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

var attributeLocationRegex = regexp.MustCompile(`^[\w]+(?:\.[\w]+)*$`)

type Config struct {
	Provider    Provider              `yaml:"provider"`
	Resources   map[string]Resource   `yaml:"resources"`
	DataSources map[string]DataSource `yaml:"data_sources"`
}

type Provider struct {
	Name      string `yaml:"name"`
	SchemaRef string `yaml:"schema_ref"`

	Ignores []string `yaml:"ignores"`
}

type Resource struct {
	Create        *OpenApiSpecLocation `yaml:"create"`
	Read          *OpenApiSpecLocation `yaml:"read"`
	Update        *OpenApiSpecLocation `yaml:"update"`
	Delete        *OpenApiSpecLocation `yaml:"delete"`
	SchemaOptions SchemaOptions        `yaml:"schema"`
}

type DataSource struct {
	Read          *OpenApiSpecLocation `yaml:"read"`
	SchemaOptions SchemaOptions        `yaml:"schema"`
}

type OpenApiSpecLocation struct {
	Path string `yaml:"path"`

	Method string `yaml:"method"`
}

type SchemaOptions struct {
	Ignores          []string         `yaml:"ignores"`
	AttributeOptions AttributeOptions `yaml:"attributes"`
}

type AttributeOptions struct {
	Aliases map[string]string `yaml:"aliases"`

	Overrides map[string]Override `yaml:"overrides"`
}

type Override struct {
	Description string `yaml:"description"`
}

func ParseConfig(bytes []byte) (*Config, error) {
	var result Config
	err := yaml.Unmarshal(bytes, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err = result.Validate(); err != nil {
		return nil, fmt.Errorf("config validation error(s):\n%w", err)
	}

	return &result, nil
}

func (c Config) Validate() error {
	var result error

	if len(c.DataSources) == 0 && len(c.Resources) == 0 {
		result = errors.Join(result, errors.New("\tat least one object is required in either 'resources' or 'data_sources'"))
	}

	err := c.Provider.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("\tprovider %w", err))
	}

	for name, resource := range c.Resources {
		err := resource.Validate()
		if err != nil {
			result = errors.Join(result, fmt.Errorf("\tresource '%s' %w", name, err))
		}
	}

	for name, dataSource := range c.DataSources {
		err := dataSource.Validate()
		if err != nil {
			result = errors.Join(result, fmt.Errorf("\tdata_source '%s' %w", name, err))
		}
	}

	return result
}

func (p Provider) Validate() error {
	var result error

	if p.Name == "" {
		result = errors.Join(result, errors.New("must have a 'name' property"))
	}

	for _, ignore := range p.Ignores {
		if !attributeLocationRegex.MatchString(ignore) {
			result = errors.Join(result, fmt.Errorf("invalid item for ignores: %q - must be dot-separated string", ignore))
		}
	}

	return result
}

func (r Resource) Validate() error {
	var result error

	if r.Create == nil {
		result = errors.Join(result, errors.New("must have a create object"))
	}
	if r.Read == nil {
		result = errors.Join(result, errors.New("must have a read object"))
	}

	err := r.Create.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid create: %w", err))
	}

	err = r.Read.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid read: %w", err))
	}

	err = r.Update.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid update: %w", err))
	}

	err = r.Delete.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid delete: %w", err))
	}

	err = r.SchemaOptions.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid schema: %w", err))
	}

	return result
}

func (d DataSource) Validate() error {
	var result error

	if d.Read == nil {
		result = errors.Join(result, errors.New("must have a read object"))
	}

	err := d.Read.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid read: %w", err))
	}

	err = d.SchemaOptions.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid schema: %w", err))
	}

	return result
}

func (o *OpenApiSpecLocation) Validate() error {
	var result error
	if o == nil {
		return nil
	}

	if o.Path == "" {
		result = errors.Join(result, errors.New("'path' property is required"))
	}

	if o.Method == "" {
		result = errors.Join(result, errors.New("'method' property is required"))
	}

	return result
}

func (s *SchemaOptions) Validate() error {
	var result error

	err := s.AttributeOptions.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid attributes: %w", err))
	}

	for _, ignore := range s.Ignores {
		if !attributeLocationRegex.MatchString(ignore) {
			result = errors.Join(result, fmt.Errorf("invalid item for ignores: %q - must be dot-separated string", ignore))
		}
	}

	return result
}

func (s *AttributeOptions) Validate() error {
	var result error

	for path := range s.Overrides {
		if !attributeLocationRegex.MatchString(path) {
			result = errors.Join(result, fmt.Errorf("invalid key for override: %q - must be dot-separated string", path))
		}
	}

	return result
}
