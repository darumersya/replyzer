package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/replyzer/analyze-repo/internal/types"
)

func AnalyzeRepository(repoPath string, options *types.AnalysisOptions) (*types.AnalysisResult, error) {
	if options.Verbose {
		fmt.Printf("Starting analysis of repository: %s\n", repoPath)
	}

	structure, err := DiscoverProjectStructure(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover project structure: %w", err)
	}

	if options.Verbose {
		fmt.Printf("Discovered %d components\n", len(structure.Components))
	}

	repoName := filepath.Base(repoPath)
	result := &types.AnalysisResult{
		Repository: types.Repository{
			Type: structure.Type,
			Path: repoPath,
			Name: repoName,
		},
		Components: make([]types.Component, 0),
	}

	for _, compInfo := range structure.Components {
		if options.Component != "" && compInfo.Name != options.Component {
			continue
		}

		if shouldExcludeComponent(compInfo.RelativePath, options.Exclude) {
			continue
		}

		if options.Verbose {
			fmt.Printf("Analyzing component: %s\n", compInfo.Name)
		}

		component, err := AnalyzeComponent(compInfo)
		if err != nil {
			if options.Verbose {
				fmt.Printf("Warning: failed to analyze component %s: %v\n", compInfo.Name, err)
			}
			continue
		}

		result.Components = append(result.Components, *component)
	}

	return result, nil
}

func AnalyzeComponent(compInfo types.ComponentInfo) (*types.Component, error) {
	langStats, err := GetLanguageStats(compInfo.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get language stats: %w", err)
	}

	primaryLang := GetPrimaryLanguage(langStats)

	framework, err := DetectFrameworks(compInfo.Path, primaryLang)
	if err != nil {
		return nil, fmt.Errorf("failed to detect framework: %w", err)
	}

	versionReqs, err := ExtractVersionRequirements(compInfo.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to extract version requirements: %w", err)
	}

	externalDeps, err := DetectExternalDependencies(compInfo.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to detect external dependencies: %w", err)
	}

	devTools, err := DetectDevelopmentTools(compInfo.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to detect development tools: %w", err)
	}

	componentType := inferComponentType(primaryLang, framework, compInfo.ConfigFiles)

	component := &types.Component{
		Name:                 compInfo.Name,
		Path:                 compInfo.RelativePath,
		Type:                 componentType,
		PrimaryLanguage:      primaryLang,
		LanguageStats:        langStats,
		Framework:            framework,
		VersionRequirements:  versionReqs,
		ExternalDependencies: *externalDeps,
		DevelopmentTools:     devTools,
	}

	return component, nil
}

func inferComponentType(primaryLang, framework string, configFiles []string) string {
	langLower := strings.ToLower(primaryLang)
	frameworkLower := strings.ToLower(framework)

	webFrameworks := []string{"react", "vue", "angular", "next.js", "nuxt"}
	apiFrameworks := []string{"express", "fastify", "nest", "django", "fastapi", "flask", "spring", "springboot", "gin", "echo", "fiber"}
	
	for _, fw := range webFrameworks {
		if strings.Contains(frameworkLower, fw) {
			return "web-application"
		}
	}

	for _, fw := range apiFrameworks {
		if strings.Contains(frameworkLower, fw) {
			return "api-service"
		}
	}

	if hasConfigFile(configFiles, "docker-compose.yml") || hasConfigFile(configFiles, "docker-compose.yaml") {
		return "service"
	}

	switch langLower {
	case "javascript", "typescript":
		if hasConfigFile(configFiles, "package.json") {
			return "web-application"
		}
	case "python":
		if hasConfigFile(configFiles, "requirements.txt") || hasConfigFile(configFiles, "pyproject.toml") {
			return "api-service"
		}
	case "java":
		if hasConfigFile(configFiles, "pom.xml") || hasConfigFile(configFiles, "build.gradle") {
			return "api-service"
		}
	case "go":
		if hasConfigFile(configFiles, "go.mod") {
			return "api-service"
		}
	case "rust":
		if hasConfigFile(configFiles, "Cargo.toml") {
			return "application"
		}
	case "c#":
		for _, file := range configFiles {
			if strings.HasSuffix(file, ".csproj") {
				return "api-service"
			}
		}
	}

	if primaryLang == "" {
		return "configuration"
	}

	return "library"
}

func hasConfigFile(configFiles []string, targetFile string) bool {
	for _, file := range configFiles {
		if file == targetFile {
			return true
		}
	}
	return false
}

func shouldExcludeComponent(relativePath string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		matched, err := filepath.Match(pattern, relativePath)
		if err == nil && matched {
			return true
		}
		
		if strings.Contains(relativePath, pattern) {
			return true
		}
	}
	return false
}