# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o crx2md .
./crx2md convert <url|file|extension-id>
./crx2md convert <source> -o output.md
./crx2md convert <source> --include-binary
```

Go proxy can be flaky in this environment. If `go get` or `go mod tidy` fails, use:
```bash
GONOSUMCHECK='*' GONOSUMDB='*' GOPROXY=direct go mod tidy
```

No tests exist yet. No linter is configured.

## Architecture

CLI tool that converts Chrome extensions (.crx/ZIP) into single LLM-optimized markdown documents.

**Pipeline:** Input detection (`cmd/convert.go`) -> CRX parsing (`internal/crx/`) -> ZIP extraction (`internal/extension/`) -> Markdown rendering (`internal/markdown/`)

### Package responsibilities

- **`cmd/`** - Cobra CLI layer. `root.go` defines the root command; `convert.go` handles input type auto-detection (URL vs file vs extension ID) and orchestrates the pipeline.
- **`internal/crx/`** - Binary format handling. `crx.go` parses CRX2/CRX3 headers to locate the ZIP payload (also passes through plain ZIPs). `download.go` fetches CRX files from Chrome Web Store via Google's update2 endpoint.
- **`internal/extension/`** - Format-agnostic extension model. `extension.go` defines types (`Extension`, `File`, `Manifest`) and file classification (code vs binary by extension). `extract.go` reads ZIP contents into the model. Designed so unpacked directory reading could be added without changing downstream code.
- **`internal/markdown/`** - Output renderer. `renderer.go` produces markdown with metadata, file tree, then fenced code blocks with language hints. manifest.json is always rendered first.

### Key design decisions

- CRX files are read entirely into memory (extensions are typically small)
- No CRX library dependency; header parsing is ~20 lines per format version
- Unknown file extensions default to "text" and are included rather than silently dropped
- Version string lives in `cmd/root.go` as a package var
- Status/progress messages go to stderr; markdown output goes to stdout

### Planned expansion points

The architecture supports future subcommands (`analyze`, `diff`, `serve`) and additional output formats (JSON renderer alongside markdown). `internal/extension/` is deliberately format-agnostic to support unpacked directories later.
