# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based CLI tool called "analyze-repo" (replyzer) that analyzes local repositories to identify development environment requirements. The tool detects languages, frameworks, version requirements, and external dependencies across single projects and monorepos.

## Key Dependencies

- `github.com/go-enry/go-enry/v2` - Language detection
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `gopkg.in/yaml.v3` - YAML processing

## Architecture

The project follows a standard Go CLI structure:

```
cmd/analyze-repo/          # CLI entry point
internal/
  analyzer/                # Core analysis logic
    analyzer.go           # Main analysis orchestration
    discovery.go          # Project structure discovery
    language.go           # Language and framework detection
    version.go            # Version requirement extraction
    dependency.go         # External dependency analysis
  config/                 # Configuration management
  output/                 # Output formatting (YAML/JSON)
  types/                  # Data structure definitions
```

## Core Data Structures

The main analysis result follows this structure:
- `AnalysisResult` containing `Repository` and `Components[]`
- Each `Component` includes language stats, framework detection, version requirements, and external dependencies
- Supports both single projects and monorepo detection

## Analysis Flow

1. **Structure Discovery**: Find configuration files (package.json, go.mod, etc.)
2. **Language Analysis**: Use go-enry for language statistics
3. **Framework Detection**: Parse configuration files for framework patterns
4. **Version Requirements**: Extract language/runtime version constraints
5. **External Dependencies**: Detect databases and services from Docker Compose, etc.

## Supported Configuration Files

- Node.js: package.json, .nvmrc
- Python: requirements.txt, pyproject.toml, .python-version
- Java: pom.xml, build.gradle
- .NET: *.csproj, global.json
- Go: go.mod
- Rust: Cargo.toml
- Docker: docker-compose.yml

## Common Commands

```bash
# Build the CLI tool
go build -o bin/analyze-repo cmd/analyze-repo/main.go

# Cross-platform builds
GOOS=windows GOARCH=amd64 go build -o bin/analyze-repo-windows.exe cmd/analyze-repo/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/analyze-repo-darwin cmd/analyze-repo/main.go
GOOS=linux GOARCH=amd64 go build -o bin/analyze-repo-linux cmd/analyze-repo/main.go

# Run tests
go test ./...

# Clean build artifacts
rm -rf bin/
```

## CLI Usage

The tool supports various output formats and filtering options:
- `--format` (yaml|json) - Output format
- `--output` - Output file path
- `--verbose` - Detailed logging
- `--component` - Analyze specific component only
- `--exclude` - Exclude patterns (glob format)