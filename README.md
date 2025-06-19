# Replyzer

A Go-based CLI tool for analyzing local repositories to identify development environment requirements. Detects languages, frameworks, version requirements, and external dependencies across single projects and monorepos.

## Features

- **Language Detection**: Automatically detects programming languages used in your repository
- **Framework Analysis**: Identifies frameworks and libraries being used
- **Version Requirements**: Extracts language and runtime version constraints
- **External Dependencies**: Detects databases and services from configuration files
- **Monorepo Support**: Analyzes both single projects and monorepos
- **Multiple Output Formats**: Supports YAML and JSON output

## Installation

### Prerequisites

- Go 1.19 or later

### Install from Source

```bash
# Clone the repository
git clone https://github.com/darumersya/replyzer.git
cd replyzer

# Build the binary
go build -o bin/analyze-repo cmd/analyze-repo/main.go

# Optional: Install to system PATH
sudo cp bin/analyze-repo /usr/local/bin/
```

### Cross-Platform Builds

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o bin/analyze-repo-windows.exe cmd/analyze-repo/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/analyze-repo-darwin cmd/analyze-repo/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/analyze-repo-linux cmd/analyze-repo/main.go
```

## Usage

### Basic Usage

```bash
# Analyze current directory
./bin/analyze-repo

# Analyze specific directory
./bin/analyze-repo /path/to/your/project
```

### Options

- `--format` (yaml|json) - Output format (default: yaml)
- `--output` - Output file path (default: stdout)
- `--verbose` - Enable detailed logging
- `--component` - Analyze specific component only
- `--exclude` - Exclude patterns (glob format)

### Examples

```bash
# Output as JSON
./bin/analyze-repo --format json

# Save to file
./bin/analyze-repo --output analysis.yaml

# Analyze specific component
./bin/analyze-repo --component frontend

# Exclude certain directories
./bin/analyze-repo --exclude "*.test,temp/*"
```

## Supported Technologies

### Languages
- Go
- JavaScript/TypeScript (Node.js)
- Python
- Java
- .NET/C#
- Rust

### Configuration Files
- `package.json`, `.nvmrc` (Node.js)
- `requirements.txt`, `pyproject.toml`, `.python-version` (Python)
- `pom.xml`, `build.gradle` (Java)
- `*.csproj`, `global.json` (.NET)
- `go.mod` (Go)
- `Cargo.toml` (Rust)
- `docker-compose.yml` (Docker services)

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

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

### Dependencies

- `github.com/go-enry/go-enry/v2` - Language detection
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `gopkg.in/yaml.v3` - YAML processing

## Releases

This project uses **automated releases** based on [semantic versioning](https://semver.org/):

- **Features** (`feat:`): Minor version bump (1.0.0 → 1.1.0)
- **Bug fixes** (`fix:`): Patch version bump (1.0.0 → 1.0.1)  
- **Breaking changes**: Major version bump (1.0.0 → 2.0.0)

Releases are automatically created when changes are merged to the `master` branch and CI passes.

### Version Management

```bash
# Check current version
./scripts/version.sh show

# Preview next version
./scripts/version.sh preview

# Validate commit messages
./scripts/version.sh validate
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using [conventional commits](https://conventionalcommits.org/) (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.