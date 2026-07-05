package main

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed config.go.tmpl
var configTemplate string

func Generate(w io.Writer) error {
	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, nil)
}
