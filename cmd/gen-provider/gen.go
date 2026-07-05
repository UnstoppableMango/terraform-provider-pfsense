package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"
)

//go:embed main.go.tmpl
var mainTemplate string

//go:embed client.go.tmpl
var clientTemplate string

type TemplateData struct {
	RegistryAddress string
	ModulePath      string
	ProviderPackage string
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
		return fmt.Errorf("")
	}

	f, err := root.Create("internal/client/client.go")
	if err != nil {
		return err
	}
	return TemplateClient(f)
}

func Generate(root *os.Root, data TemplateData) error {
	if err := GenerateMain(root, data); err != nil {
		return fmt.Errorf("generating cmd/terraform-provider-pfsense/main.go: %w", err)
	}
	if err := GenerateClient(root); err != nil {
		return fmt.Errorf("generating internal/client/client.go: %w", err)
	}

	return nil
}
