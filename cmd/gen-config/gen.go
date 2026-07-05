package main

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed config.go.tmpl
var configTemplate string

type MainData struct {
	RegistryAddress string
	ModulePath      string
	ProviderPackage string
}

func Generate(w io.Writer, data MainData) error {
	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
