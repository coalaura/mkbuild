package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/coalaura/plain"
)

var Version = "dev"

var log = plain.New()

type CLI struct {
	Name     string `arg:"" help:"Name of the project (used for binary and artifact naming)."`
	MainPath string `name:"main" default:"." help:"Path to the main package (e.g., ./cmd/app)."`

	CGO       bool `name:"cgo" help:"Enable CGO and cross-compilation via Zig."`
	Checksums bool `name:"checksums" help:"Generate SHA256 checksums file."`
	Draft     bool `name:"draft" help:"Create releases as draft (not published immediately)."`
	Archive   bool `name:"archive" help:"Create .tar.gz/.zip archives instead of raw binaries."`

	Interactive bool `short:"i" name:"interactive" help:"Configure options interactively."`
	DryRun      bool `short:"n" name:"dry-run" help:"Preview generated file without writing."`
	Help        bool `short:"h" help:"Show detailed help."`
	Version     bool `short:"v" help:"Print version."`
}

func main() {
	var cli CLI

	kong.Parse(&cli,
		kong.Name("mkbuild"),
		kong.Description("GitHub Actions release workflow generator for Go"),
		kong.NoDefaultHelp(),
	)

	if cli.Version {
		version()

		return
	}

	if cli.Help {
		help()

		return
	}

	cfg := &BuildConfig{
		Name:      cli.Name,
		CGO:       cli.CGO,
		MainPath:  cli.MainPath,
		Checksums: cli.Checksums,
		Draft:     cli.Draft,
		Archive:   cli.Archive,
	}

	if cli.Interactive {
		runInteractive(cfg)
	}

	var buf bytes.Buffer

	err := ReleaseTmpl.Execute(&buf, cfg)
	log.MustFail(err)

	if cli.DryRun {
		dryRun(cfg, buf.String())

		return
	}

	outPath := filepath.Join(".github", "workflows", "release.yml")

	err = os.MkdirAll(filepath.Dir(outPath), 0755)
	log.MustFail(err)

	err = os.WriteFile(outPath, buf.Bytes(), 0644)
	log.MustFail(err)

	log.Printf("Wrote %s\n", outPath)
}

func version() {
	log.Printf("mkbuild %s\n", Version)
}

func dryRun(cfg *BuildConfig, content string) {
	log.Println("Dry run - no files written.")
	log.Println()
	log.Println("Configuration:")
	log.Printf("  Name:       %s\n", cfg.Name)
	log.Printf("  Main Path:  %s\n", cfg.MainPath)
	log.Println()
	log.Println("Build Options:")
	log.Printf("  CGO:        %v\n", cfg.CGO)
	log.Println()
	log.Println("Release Options:")
	log.Printf("  Checksums:  %v\n", cfg.Checksums)
	log.Printf("  Draft:      %v\n", cfg.Draft)
	log.Printf("  Archive:    %v\n", cfg.Archive)
	log.Println()
	log.Println("Generated Workflow:")
	log.Println(content)
}

func help() {
	log.Println("mkbuild - GitHub Actions release workflow generator for Go")
	log.Println()
	log.Println("Usage: mkbuild <name> [options]")
	log.Println()
	log.Println("Options:")
	log.Println("  --cgo          Enable CGO and cross-compilation via Zig.")
	log.Println("  --main <path>  Path to the main package (default: '.').")
	log.Println("  --checksums    Generate SHA256 checksums (default: true).")
	log.Println("  --draft        Create releases as draft.")
	log.Println("  --archive      Create archives instead of raw binaries (default: true).")
	log.Println("  -i, --interactive  Configure options interactively.")
	log.Println("  -n, --dry-run  Preview generated file without writing.")
	log.Println("  -v, --version  Print version.")
	log.Println("  -h, --help     Show this help.")
}

func runInteractive(cfg *BuildConfig) {
	log.Println("Interactive Configuration")
	log.Println("Press Enter to accept defaults.")
	log.Println()

	// Project
	cfg.Name = askString("Project Name", cfg.Name)
	cfg.MainPath = askString("Main Path", cfg.MainPath)

	// Build
	cfg.CGO = ask(
		"CGO",
		"Enable CGO and cross-compilation via Zig.",
		cfg.CGO,
	)

	// Release
	cfg.Checksums = ask(
		"Checksums",
		"Generate SHA256 checksums file for release artifacts.",
		cfg.Checksums,
	)

	cfg.Draft = ask(
		"Draft Release",
		"Create releases as draft (not published immediately).",
		cfg.Draft,
	)

	cfg.Archive = ask(
		"Archives",
		"Create .tar.gz/.zip archives instead of raw binaries.",
		cfg.Archive,
	)
}

func ask(title, desc string, def bool) bool {
	log.Println()

	log.Println(title)
	log.Printf("  %s\n", desc)

	val, err := log.ConfirmWithEcho("  Enable", def, " ")
	log.MustExit(err)

	return val
}

func askString(prompt, def string) string {
	display := def

	if display == "" {
		display = "none"
	}

	log.Printf("%s [%s]: ", prompt, display)

	val, err := log.Read("", 12)
	log.MustFail(err)

	val = strings.TrimSpace(val)

	if val == "" {
		return def
	}

	return val
}
