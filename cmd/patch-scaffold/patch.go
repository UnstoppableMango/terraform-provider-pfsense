package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
)

type replacement struct {
	start, end int
	body       string
}

func Patch(providerFile string) error {
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
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Body == nil {
			continue
		}
		if !isPfsenseProviderMethod(fn) {
			continue
		}
		body := delegationBody(fn.Name.Name)
		if body == "" {
			continue
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

	formatted, err := format.Source(out)
	if err != nil {
		return fmt.Errorf("format: %w", err)
	}
	_, err = os.Stdout.Write(formatted)
	return err
}

func delegationBody(method string) string {
	switch method {
	case "Resources":
		return "\n\treturn pfsenseResources()\n"
	case "DataSources":
		return "\n\treturn pfsenseDataSources()\n"
	case "Schema":
		return "\n\tresp.Schema = pfsenseSchema()\n"
	case "Configure":
		return "\n\tpfsenseConfigure(ctx, req, resp)\n"
	case "Metadata":
		return "\n\tpfsenseMetadata(req, resp)\n"
	}
	return ""
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
