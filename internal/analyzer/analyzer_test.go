package analyzer

import (
	"testing"

	"github.com/replyzer/analyze-repo/internal/types"
)

func TestInferComponentType(t *testing.T) {
	tests := []struct {
		name        string
		primaryLang string
		framework   string
		configFiles []string
		expected    string
	}{
		{
			name:        "React web application",
			primaryLang: "JavaScript",
			framework:   "React",
			configFiles: []string{"package.json"},
			expected:    "web-application",
		},
		{
			name:        "Vue web application",
			primaryLang: "JavaScript",
			framework:   "Vue.js",
			configFiles: []string{"package.json"},
			expected:    "web-application",
		},
		{
			name:        "Express API service",
			primaryLang: "JavaScript",
			framework:   "Express",
			configFiles: []string{"package.json"},
			expected:    "api-service",
		},
		{
			name:        "Go API service",
			primaryLang: "Go",
			framework:   "Gin",
			configFiles: []string{"go.mod"},
			expected:    "api-service",
		},
		{
			name:        "Python API service",
			primaryLang: "Python",
			framework:   "Django",
			configFiles: []string{"requirements.txt"},
			expected:    "api-service",
		},
		{
			name:        "Rust application",
			primaryLang: "Rust",
			framework:   "",
			configFiles: []string{"Cargo.toml"},
			expected:    "application",
		},
		{
			name:        "Docker service",
			primaryLang: "Go",
			framework:   "",
			configFiles: []string{"docker-compose.yml"},
			expected:    "service",
		},
		{
			name:        "Configuration only",
			primaryLang: "",
			framework:   "",
			configFiles: []string{"config.yaml"},
			expected:    "configuration",
		},
		{
			name:        "Library",
			primaryLang: "JavaScript",
			framework:   "",
			configFiles: []string{"package.json"},
			expected:    "web-application",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferComponentType(tt.primaryLang, tt.framework, tt.configFiles)
			if result != tt.expected {
				t.Errorf("inferComponentType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasConfigFile(t *testing.T) {
	configFiles := []string{"package.json", "go.mod", "docker-compose.yml"}

	tests := []struct {
		name       string
		targetFile string
		expected   bool
	}{
		{"existing file", "package.json", true},
		{"existing file", "go.mod", true},
		{"non-existing file", "requirements.txt", false},
		{"case sensitive", "Package.json", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasConfigFile(configFiles, tt.targetFile)
			if result != tt.expected {
				t.Errorf("hasConfigFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldExcludeComponent(t *testing.T) {
	tests := []struct {
		name            string
		relativePath    string
		excludePatterns []string
		expected        bool
	}{
		{
			name:            "exact match",
			relativePath:    "node_modules",
			excludePatterns: []string{"node_modules"},
			expected:        true,
		},
		{
			name:            "glob pattern match",
			relativePath:    "test.txt",
			excludePatterns: []string{"*.txt"},
			expected:        true,
		},
		{
			name:            "substring match",
			relativePath:    "src/node_modules/lib",
			excludePatterns: []string{"node_modules"},
			expected:        true,
		},
		{
			name:            "no match",
			relativePath:    "src/main.go",
			excludePatterns: []string{"*.txt", "node_modules"},
			expected:        false,
		},
		{
			name:            "empty patterns",
			relativePath:    "src/main.go",
			excludePatterns: []string{},
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldExcludeComponent(tt.relativePath, tt.excludePatterns)
			if result != tt.expected {
				t.Errorf("shouldExcludeComponent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAnalysisOptions(t *testing.T) {
	options := &types.AnalysisOptions{
		Format:    "yaml",
		Output:    "output.yaml",
		Verbose:   true,
		Component: "frontend",
		Exclude:   []string{"*.test", "node_modules"},
	}

	if options.Format != "yaml" {
		t.Errorf("Expected format to be 'yaml', got %v", options.Format)
	}

	if options.Verbose != true {
		t.Errorf("Expected verbose to be true, got %v", options.Verbose)
	}

	if len(options.Exclude) != 2 {
		t.Errorf("Expected 2 exclude patterns, got %v", len(options.Exclude))
	}
}