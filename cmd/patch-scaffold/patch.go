package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"sort"
	"strings"
	"unicode"

	providerspec "github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	specschema "github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"
	"golang.org/x/tools/go/ast/astutil"
)

type replacement struct {
	start, end int
	body       string
}

func Patch(providerFile, schemaFile string) error {
	s, err := ParseSchema(schemaFile)
	if err != nil {
		return fmt.Errorf("parse schema: %w", err)
	}

	src, err := os.ReadFile(providerFile)
	if err != nil {
		return fmt.Errorf("read provider file: %w", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, providerFile, src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse provider file: %w", err)
	}

	var replacements []replacement
	neededImports := map[string]string{}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Body == nil {
			continue
		}
		if !isPfsenseProviderMethod(fn) {
			continue
		}
		body, imports := generateBody(fn.Name.Name, s)
		if body == "" {
			continue
		}
		maps.Copy(neededImports, imports)

		lbrace := fset.Position(fn.Body.Lbrace).Offset
		rbrace := fset.Position(fn.Body.Rbrace).Offset
		replacements = append(replacements, replacement{lbrace + 1, rbrace, body})
	}

	sort.Slice(replacements, func(i, j int) bool {
		return replacements[i].start > replacements[j].start
	})
	out := make([]byte, len(src))
	copy(out, src)
	for _, r := range replacements {
		out = append(out[:r.start], append([]byte(r.body), out[r.end:]...)...)
	}
	out = append(out, []byte("\nvar Version = \"dev\"\n")...)

	fset2 := token.NewFileSet()
	f2, err := parser.ParseFile(fset2, providerFile, out, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("re-parse after splice: %w", err)
	}
	for path, alias := range neededImports {
		if alias != "" {
			astutil.AddNamedImport(fset2, f2, alias, path)
		} else {
			astutil.AddImport(fset2, f2, path)
		}
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset2, f2); err != nil {
		return fmt.Errorf("format: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format source: %w", err)
	}
	if _, err = os.Stdout.Write(formatted); err != nil {
		return err
	}
	return nil
}

func isPfsenseProviderMethod(fn *ast.FuncDecl) bool {
	if len(fn.Recv.List) == 0 {
		return false
	}
	star, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	ident, ok := star.X.(*ast.Ident)
	return ok && ident.Name == "pfsenseProvider"
}

func generateBody(method string, s *spec.Specification) (string, map[string]string) {
	switch method {
	case "Resources":
		return generateResources(s)
	case "DataSources":
		return generateDataSources(s)
	case "Schema":
		return generateSchema(s)
	case "Configure":
		return generateConfigure()
	case "Metadata":
		return generateMetadata(), nil
	}
	return "", nil
}

func generateMetadata() string {
	return "\n\tresp.TypeName = \"pfsense\"\n\tresp.Version = Version\n"
}

func generateResources(s *spec.Specification) (string, map[string]string) {
	if len(s.Resources) == 0 {
		return "\n\treturn nil\n", nil
	}
	imports := map[string]string{}
	var b strings.Builder
	b.WriteString("\n\treturn []func() resource.Resource{\n")
	for _, r := range s.Resources {
		pkg := "resource_" + r.Name + "_resource"
		constructor := "New" + toPascalCase(r.Name) + "Resource"
		fmt.Fprintf(&b, "\t\t%s.%s,\n", pkg, constructor)
		imports["github.com/unstoppablemango/terraform-provider-pfsense/"+pkg] = ""
	}
	b.WriteString("\t}\n")
	return b.String(), imports
}

func generateDataSources(s *spec.Specification) (string, map[string]string) {
	if len(s.DataSources) == 0 {
		return "\n\treturn nil\n", nil
	}
	imports := map[string]string{}
	var b strings.Builder
	b.WriteString("\n\treturn []func() datasource.DataSource{\n")
	for _, ds := range s.DataSources {
		pkg := "datasource_" + ds.Name + "_data_source"
		constructor := "New" + toPascalCase(ds.Name) + "DataSource"
		fmt.Fprintf(&b, "\t\t%s.%s,\n", pkg, constructor)
		imports["github.com/unstoppablemango/terraform-provider-pfsense/"+pkg] = ""
	}
	b.WriteString("\t}\n")
	return b.String(), imports
}

func generateSchema(s *spec.Specification) (string, map[string]string) {
	imports := map[string]string{
		"github.com/hashicorp/terraform-plugin-framework/provider/schema": "",
	}

	var b strings.Builder
	b.WriteString("\n\tresp.Schema = schema.Schema{\n")
	b.WriteString("\t\tAttributes: map[string]schema.Attribute{\n")
	b.WriteString("\t\t\t\"host\":     schema.StringAttribute{Required: true},\n")
	b.WriteString("\t\t\t\"username\": schema.StringAttribute{Required: true},\n")
	b.WriteString("\t\t\t\"password\": schema.StringAttribute{Required: true, Sensitive: true},\n")
	if s.Provider != nil && s.Provider.Schema != nil {
		for _, attr := range s.Provider.Schema.Attributes {
			lit := attrLiteral(attr)
			if lit == "" {
				continue
			}
			fmt.Fprintf(&b, "\t\t\t%q: %s,\n", attr.Name, lit)
		}
	}
	b.WriteString("\t\t},\n")
	b.WriteString("\t}\n")
	return b.String(), imports
}

func generateConfigure() (string, map[string]string) {
	imports := map[string]string{
		"net/http": "",
		"time":     "",
		"github.com/hashicorp/terraform-plugin-framework/path":                   "",
		"github.com/hashicorp/terraform-plugin-framework/types":                  "",
		"github.com/unstoppablemango/terraform-provider-pfsense/internal/client": "",
	}
	body := `
	var host, username, password types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("host"), &host)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("username"), &username)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password"), &password)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cfg := &client.Config{
		Host:     host.ValueString(),
		Username: username.ValueString(),
		Password: password.ValueString(),
		HTTP:     &http.Client{Timeout: 30 * time.Second},
	}
	resp.DataSourceData = cfg
	resp.ResourceData = cfg
`
	return body, imports
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
