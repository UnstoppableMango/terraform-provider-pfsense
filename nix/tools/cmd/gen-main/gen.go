package main

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed main.go.tmpl
var mainTemplate string

type MainData struct {
	RegistryAddress string
	ModulePath      string
	ProviderPackage string
}

func GenerateMain(w io.Writer, data MainData) error {
	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
