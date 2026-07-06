package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"

	providerspec "github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	specschema "github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"
)

//go:embed main.go.tmpl
var mainTemplate string

//go:embed client.go.tmpl
var clientTemplate string

//go:embed provider_gen.go.tmpl
var providerImplTemplate string

type ResourceItem struct {
	Name         string
	PascalName   string
	PackageAlias string
	ImportPath   string
}

type DataSourceItem struct {
	Name         string
	PascalName   string
	PackageAlias string
	ImportPath   string
}

type ProviderAttrItem struct {
	Name     string
	TypeExpr string
}

type TemplateData struct {
	RegistryAddress string
	ModulePath      string
	ProviderPackage string
	Resources       []ResourceItem
	DataSources     []DataSourceItem
	ProviderAttrs   []ProviderAttrItem
}

func TemplateMain(w io.Writer, data TemplateData) error {
	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func TemplateClient(w io.Writer) error {
	tmpl, err := template.New("client").Parse(clientTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, nil)
}

func TemplateProviderImpl(w io.Writer, data TemplateData) error {
	tmpl, err := template.New("provider_gen").Parse(providerImplTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func GenerateMain(root *os.Root, data TemplateData) error {
	if err := root.MkdirAll("cmd/terraform-provider-pfsense", os.ModePerm); err != nil {
		return err
	}
	f, err := root.Create("cmd/terraform-provider-pfsense/main.go")
	if err != nil {
		return err
	}
	return TemplateMain(f, data)
}

func GenerateClient(root *os.Root) error {
	if err := root.MkdirAll("internal/client", os.ModePerm); err != nil {
		return fmt.Errorf("mkdir internal/client: %w", err)
	}
	f, err := root.Create("internal/client/client.go")
	if err != nil {
		return err
	}
	return TemplateClient(f)
}

func GenerateProviderImpl(root *os.Root, data TemplateData) error {
	if err := root.MkdirAll(data.ProviderPackage, os.ModePerm); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := TemplateProviderImpl(&buf, data); err != nil {
		return err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format provider_gen.go: %w\n%s", err, buf.String())
	}
	f, err := root.Create(data.ProviderPackage + "/provider_gen.go")
	if err != nil {
		return err
	}
	_, err = f.Write(formatted)
	return err
}

func Generate(root *os.Root, data TemplateData, s *spec.Specification) error {
	populateProviderData(&data, s)
	if err := GenerateMain(root, data); err != nil {
		return fmt.Errorf("generating cmd/terraform-provider-pfsense/main.go: %w", err)
	}
	if err := GenerateClient(root); err != nil {
		return fmt.Errorf("generating internal/client/client.go: %w", err)
	}
	if err := GenerateProviderImpl(root, data); err != nil {
		return fmt.Errorf("generating %s/provider_gen.go: %w", data.ProviderPackage, err)
	}
	return nil
}

func populateProviderData(data *TemplateData, s *spec.Specification) {
	for _, r := range s.Resources {
		pkg := "resource_" + r.Name + "_resource"
		data.Resources = append(data.Resources, ResourceItem{
			Name:         r.Name,
			PascalName:   toPascalCase(r.Name),
			PackageAlias: pkg,
			ImportPath:   data.ModulePath + "/" + pkg,
		})
	}
	for _, ds := range s.DataSources {
		pkg := "datasource_" + ds.Name + "_data_source"
		data.DataSources = append(data.DataSources, DataSourceItem{
			Name:         ds.Name,
			PascalName:   toPascalCase(ds.Name),
			PackageAlias: pkg,
			ImportPath:   data.ModulePath + "/" + pkg,
		})
	}
	if s.Provider != nil && s.Provider.Schema != nil {
		for _, attr := range s.Provider.Schema.Attributes {
			lit := attrLiteral(attr)
			if lit == "" {
				continue
			}
			data.ProviderAttrs = append(data.ProviderAttrs, ProviderAttrItem{
				Name:     attr.Name,
				TypeExpr: lit,
			})
		}
	}
}

func attrLiteral(attr providerspec.Attribute) string {
	switch {
	case attr.String != nil:
		return "schema.StringAttribute{" + corFields(attr.String.OptionalRequired) + "}"
	case attr.Int64 != nil:
		return "schema.Int64Attribute{" + corFields(attr.Int64.OptionalRequired) + "}"
	case attr.Bool != nil:
		return "schema.BoolAttribute{" + corFields(attr.Bool.OptionalRequired) + "}"
	case attr.Float64 != nil:
		return "schema.Float64Attribute{" + corFields(attr.Float64.OptionalRequired) + "}"
	case attr.Number != nil:
		return "schema.NumberAttribute{" + corFields(attr.Number.OptionalRequired) + "}"
	}
	return ""
}

func corFields(cor specschema.ComputedOptionalRequired) string {
	switch cor {
	case specschema.Required:
		return "Required: true"
	case specschema.Optional:
		return "Optional: true"
	case specschema.Computed:
		return "Computed: true"
	case specschema.ComputedOptional:
		return "Computed: true, Optional: true"
	default:
		return "Optional: true"
	}
}

func toPascalCase(s string) string {
	var b strings.Builder
	nextUpper := true
	for _, r := range s {
		if r == '_' || r == '-' {
			nextUpper = true
			continue
		}
		if nextUpper {
			b.WriteRune(unicode.ToUpper(r))
			nextUpper = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
