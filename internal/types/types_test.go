package types

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestAnalysisResultSerialization(t *testing.T) {
	// Create a sample AnalysisResult
	result := AnalysisResult{
		Repository: Repository{
			Type: "single",
			Path: "/path/to/repo",
			Name: "test-repo",
		},
		Components: []Component{
			{
				Name:            "frontend",
				Path:            "./frontend",
				Type:            "web-application",
				PrimaryLanguage: "JavaScript",
				LanguageStats: map[string]float64{
					"JavaScript": 80.5,
					"CSS":        15.2,
					"HTML":       4.3,
				},
				Framework: "React",
				VersionRequirements: map[string]string{
					"node": ">=18.0.0",
					"npm":  ">=8.0.0",
				},
				ExternalDependencies: ExternalDependencies{
					Databases: []string{"PostgreSQL"},
					Services:  []string{"Redis"},
				},
				DevelopmentTools: []string{"webpack", "eslint"},
			},
		},
	}

	t.Run("JSON serialization", func(t *testing.T) {
		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal to JSON: %v", err)
		}

		var unmarshaled AnalysisResult
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal from JSON: %v", err)
		}

		if unmarshaled.Repository.Name != result.Repository.Name {
			t.Errorf("Repository name mismatch. Expected: %s, Got: %s", result.Repository.Name, unmarshaled.Repository.Name)
		}

		if len(unmarshaled.Components) != len(result.Components) {
			t.Errorf("Components count mismatch. Expected: %d, Got: %d", len(result.Components), len(unmarshaled.Components))
		}

		if len(unmarshaled.Components) > 0 {
			comp := unmarshaled.Components[0]
			if comp.Framework != "React" {
				t.Errorf("Framework mismatch. Expected: React, Got: %s", comp.Framework)
			}
		}
	})

	t.Run("YAML serialization", func(t *testing.T) {
		data, err := yaml.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal to YAML: %v", err)
		}

		var unmarshaled AnalysisResult
		err = yaml.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal from YAML: %v", err)
		}

		if unmarshaled.Repository.Type != result.Repository.Type {
			t.Errorf("Repository type mismatch. Expected: %s, Got: %s", result.Repository.Type, unmarshaled.Repository.Type)
		}

		if len(unmarshaled.Components) > 0 {
			comp := unmarshaled.Components[0]
			if comp.PrimaryLanguage != "JavaScript" {
				t.Errorf("Primary language mismatch. Expected: JavaScript, Got: %s", comp.PrimaryLanguage)
			}
		}
	})
}

func TestComponentValidation(t *testing.T) {
	tests := []struct {
		name      string
		component Component
		valid     bool
	}{
		{
			name: "valid component",
			component: Component{
				Name:            "test-component",
				Path:            "./test",
				Type:            "web-application",
				PrimaryLanguage: "JavaScript",
				LanguageStats:   map[string]float64{"JavaScript": 100.0},
				Framework:       "React",
				VersionRequirements: map[string]string{
					"node": ">=16.0.0",
				},
				ExternalDependencies: ExternalDependencies{
					Databases: []string{},
					Services:  []string{},
				},
				DevelopmentTools: []string{"webpack"},
			},
			valid: true,
		},
		{
			name: "component with minimal fields",
			component: Component{
				Name:                 "minimal",
				Path:                 "./minimal",
				Type:                 "library",
				PrimaryLanguage:      "Go",
				LanguageStats:        map[string]float64{"Go": 100.0},
				Framework:            "",
				VersionRequirements:  map[string]string{},
				ExternalDependencies: ExternalDependencies{},
				DevelopmentTools:     []string{},
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check required fields are not empty
			if tt.valid {
				if tt.component.Name == "" {
					t.Error("Component name should not be empty for valid component")
				}
				if tt.component.Path == "" {
					t.Error("Component path should not be empty for valid component")
				}
				if tt.component.Type == "" {
					t.Error("Component type should not be empty for valid component")
				}
			}
		})
	}
}

func TestAnalysisOptions(t *testing.T) {
	options := AnalysisOptions{
		Format:    "yaml",
		Output:    "analysis.yaml",
		Verbose:   true,
		Component: "frontend",
		Exclude:   []string{"*.test", "node_modules"},
	}

	if options.Format != "yaml" {
		t.Errorf("Expected format 'yaml', got %s", options.Format)
	}

	if !options.Verbose {
		t.Error("Expected verbose to be true")
	}

	if len(options.Exclude) != 2 {
		t.Errorf("Expected 2 exclude patterns, got %d", len(options.Exclude))
	}
}

func TestProjectStructure(t *testing.T) {
	structure := ProjectStructure{
		Type: "monorepo",
		Components: []ComponentInfo{
			{
				Name:         "frontend",
				Path:         "/repo/frontend",
				ConfigFiles:  []string{"package.json", ".eslintrc"},
				RelativePath: "frontend",
			},
			{
				Name:         "backend",
				Path:         "/repo/backend",
				ConfigFiles:  []string{"go.mod", "go.sum"},
				RelativePath: "backend",
			},
		},
	}

	if structure.Type != "monorepo" {
		t.Errorf("Expected type 'monorepo', got %s", structure.Type)
	}

	if len(structure.Components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(structure.Components))
	}

	frontend := structure.Components[0]
	if frontend.Name != "frontend" {
		t.Errorf("Expected first component name 'frontend', got %s", frontend.Name)
	}

	if len(frontend.ConfigFiles) != 2 {
		t.Errorf("Expected 2 config files for frontend, got %d", len(frontend.ConfigFiles))
	}
}