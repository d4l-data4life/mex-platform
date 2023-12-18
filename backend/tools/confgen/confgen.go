package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/d4l-data4life/mex/mex/shared/cfg"
)

// Protoc plugin to derive the following files from a proto file that uses our config option annotations:
// - a Markdown describing the configuration in a human-readable form,
// - a K8s ConfigMap file containing all non-secret (and non-ignored) fields,
// - a K8s Secret file  containing all secret (and non-ignored) fields.
// See the global Makefile for invocation details.

func logf(f string, args ...any) {
	fmt.Fprintf(os.Stderr, f, args...)
}

func main() {
	logf("confgen\n")
	flags := flag.FlagSet{}

	// The `output` parameter determines which files are emitted.
	rawOutput := flags.String("output", "md", "which output to generate (md, k8s)")
	// Set `comments` to true and some comments will be added to the K8s files.
	rawComments := flags.Bool("comments", false, "generate comments (where applicable)")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = 1 // FEATURE_PROTO3_OPTIONAL

		logf("flags: %v, '%s'\n", *rawComments, *rawOutput)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			switch *rawOutput {
			case "md":
				generateMarkdown(gen, f)
			case "k8s":
				generateK8sMaps(gen, f, *rawComments)
			default:
				return fmt.Errorf("unsupported output: '%s'", *rawOutput)
			}
		}
		logf("done.\n")
		return nil
	})
}

const (
	symbCheckmark = "✅"
	symbLock      = "🔒"
	symbExclam    = "❗"
)

func usedBy(tag string, tags []string) bool {
	if tag == "*" {
		return true
	}

	if slices.Contains(tags, "*") {
		return true
	}

	return slices.Contains(tags, tag)
}

func defaultOrNone(def string) string {
	if def == "" {
		return "_none_"
	}
	return fmt.Sprintf("`'%s'`", def)
}

func generateMarkdown(gen *protogen.Plugin, file *protogen.File) {
	logf("- filename: %s.md\n", file.GeneratedFilenamePrefix)

	filename := file.GeneratedFilenamePrefix + ".cfg.md"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("# MEx configuration")
	g.P()

	for _, msg := range file.Messages {
		logf("  - message: %s\n", msg.Desc.FullName())

		overview, err := cfg.GetMessageOverview(msg.Desc, "MEX")
		if err != nil {
			logf("    (warn: %s)\n", err.Error())
			continue
		}

		g.P("## Overview")
		g.P("| `met` | `idx` | `qry` | `cfg` | `aut` | Go struct field | Type | Secret | Environment variable | Default value | Title |")
		g.P("| ----- | ----- | ----- | ----- | ----- | --------------- | ---- | ------ | -------------------- | ------------- | ----- |")

		for i := 0; i < overview.Len(); i++ {
			fov := overview.Get(i)
			g.P(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s%s | %s | %s `%s` | %s | %s |",
				tern(usedBy("metadata", fov.EffTags), symbCheckmark, ""),
				tern(usedBy("index", fov.EffTags), symbCheckmark, ""),
				tern(usedBy("query", fov.EffTags), symbCheckmark, ""),
				tern(usedBy("config", fov.EffTags), symbCheckmark, ""),
				tern(usedBy("auth", fov.EffTags), symbCheckmark, ""),
				fov.GoPath,
				tern(fov.Repeated, "[]", ""), fov.Kind.String(),
				tern(fov.Secret, symbLock, ""),
				tern(fov.AltEnvName != "" || fov.AliasEnvName != "", symbExclam, ""),
				fov.EffEnvName,
				defaultOrNone(fov.Default),
				fov.Title))
		}

		g.P("## Configuration details")

		for i := 0; i < overview.Len(); i++ {
			fov := overview.Get(i)

			g.P(fmt.Sprintf("### `%s`: %s", fov.EffEnvName, fov.Title))

			if fov.Summary != "" {
				g.P("#### Summary\n")
				g.P(fov.Summary)
			}

			g.P("#### Info\n")

			g.P("| Key | Value |")
			g.P("| --- | ----- |")

			g.P(fmt.Sprintf("| Go struct field: | `%s` |", fov.GoPath))
			g.P(fmt.Sprintf("| Environment variable: | `%s` %s |", fov.EffEnvName, tern(fov.AltEnvName != "", "(Note the name deviation!)", "")))
			if fov.AliasEnvName != "" {
				g.P(fmt.Sprintf("| Vault source variable: | `%s` |", fov.AliasEnvName))
			}

			if fov.Default != "" {
				g.P(fmt.Sprintf("| Default value: | `'%s'` |", fov.Default))
			}

			if fov.Secret {
				g.P("| Secret: | **yes** |")
			}

			sbt := strings.Builder{}
			sbt.WriteString("<ul>")
			for _, tag := range fov.EffTags {
				if tag == "*" {
					sbt.WriteString("<li>_all_</li>")
				} else {
					sbt.WriteString(fmt.Sprintf("<li>%s</li>", tag))
				}
			}
			sbt.WriteString("</ul>")

			g.P(fmt.Sprintf("| Used by: | %s |", sbt.String()))

			if fov.Description != "" {
				g.P("#### Description\n")
				g.P(fov.Description)
			}

			g.P("\n----")
		}
	}
}

func generateK8sMaps(gen *protogen.Plugin, file *protogen.File, genComments bool) {
	logf(">> filename: %s.md\n", file.GeneratedFilenamePrefix)

	filenameConfigMap := file.GeneratedFilenamePrefix + ".configmap.yaml"
	filenameSecret := file.GeneratedFilenamePrefix + ".secret.yaml"
	gC := gen.NewGeneratedFile(filenameConfigMap, file.GoImportPath)
	gS := gen.NewGeneratedFile(filenameSecret, file.GoImportPath)

	gC.P("apiVersion: v1")
	gC.P("kind: ConfigMap")
	gC.P("metadata:")
	gC.P("  name: mex-services-{{ .Values.tenantID }}-config")
	gC.P("  labels:")
	gC.P("    product: mex")
	gC.P("    env: {{ .Values.environment }}")
	gC.P("    tenant: {{ .Values.tenantID }}")
	gC.P("data:")

	gS.P("apiVersion: v1")
	gS.P("kind: Secret")
	gS.P("metadata:")
	gS.P("  name: mex-services-{{ .Values.tenantID }}-secret")
	gS.P("  labels:")
	gS.P("    product: mex")
	gS.P("    env: {{ .Values.environment }}")
	gS.P("    tenant: {{ .Values.tenantID }}")
	gS.P("type: Opaque")
	gS.P("data:")

	for _, msg := range file.Messages {
		logf(" - message %s\n", msg.Desc.FullName())

		overview, err := cfg.GetMessageOverview(msg.Desc, "MEX")
		if err != nil {
			logf("WARN: %s\n", err.Error())
			continue
		}

		for i := 0; i < overview.Len(); i++ {
			fov := overview.Get(i)

			if fov.IgnoreK8s {
				continue
			}

			if fov.Secret {
				if genComments {
					gS.P("  # ", fov.Title)
				}
				gS.P(fmt.Sprintf("  %s: {{ .Values.%s %s| b64enc }}",
					fov.EffEnvName, tern(fov.AliasEnvName != "", fov.AliasEnvName, fov.EffEnvName),
					tern(fov.Kind == protoreflect.BytesKind, "| b64enc ", ""),
				))
			} else {
				if genComments {
					gC.P("  # ", fov.Title)
				}
				if fov.Default == "" {
					gC.P()
					gC.P("  # no default: ", fov.EffEnvName)
				}
				gC.P(fmt.Sprintf("  %s: {{ .Values.%s %s%s| quote }}",
					fov.EffEnvName, tern(fov.AliasEnvName != "", fov.AliasEnvName, fov.EffEnvName),
					tern(fov.Default != "", fmt.Sprintf(`| default "%s" `, fov.Default), ""),
					tern(fov.Kind == protoreflect.BytesKind, "| b64enc ", ""),
				))
			}
		}
	}
}

func tern(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
