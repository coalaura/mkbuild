package main

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed workflow.tmpl
	releaseStr string

	ReleaseTmpl = template.Must(template.New("release").Delims("<<", ">>").Parse(releaseStr))
)

type BuildConfig struct {
	Name      string
	MainPath  string
	CGO       bool
	Checksums bool
	Draft     bool
	Archive   bool
}
