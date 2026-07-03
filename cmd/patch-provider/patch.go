package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
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
		for path, alias := range imports {
			neededImports[path] = alias
		}

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
	_, err = os.Stdout.Write(formatted)
	return err
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
		return generateConfigure(), nil
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

	if s.Provider == nil || s.Provider.Schema == nil || len(s.Provider.Schema.Attributes) == 0 {
		return "\n\tresp.Schema = schema.Schema{}\n", imports
	}

	var b strings.Builder
	b.WriteString("\n\tresp.Schema = schema.Schema{\n")
	b.WriteString("\t\tAttributes: map[string]schema.Attribute{\n")
	for _, attr := range s.Provider.Schema.Attributes {
		lit := attrLiteral(attr)
		if lit == "" {
			continue
		}
		fmt.Fprintf(&b, "\t\t\t%q: %s,\n", attr.Name, lit)
	}
	b.WriteString("\t\t},\n")
	b.WriteString("\t}\n")
	return b.String(), imports
}

func generateConfigure() string {
	return "\n\t// TODO: initialize API client\n"
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
