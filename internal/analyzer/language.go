package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-enry/go-enry/v2"
)

var frameworkPatterns = map[string][]string{
	"React":      {"react", "@types/react"},
	"Vue":        {"vue", "@vue/cli", "vue-cli"},
	"Angular":    {"@angular/core", "@angular/cli"},
	"Next.js":    {"next"},
	"Nuxt":       {"nuxt"},
	"Express":    {"express"},
	"Fastify":    {"fastify"},
	"Nest":       {"@nestjs/core"},
	"Django":     {"django", "Django"},
	"FastAPI":    {"fastapi"},
	"Flask":      {"flask", "Flask"},
	"Spring":     {"org.springframework"},
	"SpringBoot": {"org.springframework.boot"},
	"Laravel":    {"laravel/framework"},
	"Symfony":    {"symfony/framework-bundle"},
	"Rails":      {"rails"},
	"Gin":        {"github.com/gin-gonic/gin"},
	"Echo":       {"github.com/labstack/echo"},
	"Fiber":      {"github.com/gofiber/fiber"},
	"Axum":       {"axum"},
	"Actix":      {"actix-web"},
	"Rocket":     {"rocket"},
	"ASP.NET":    {"Microsoft.AspNetCore"},
}

func GetLanguageStats(path string) (map[string]float64, error) {
	var languages map[string]int
	var err error
	
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return make(map[string]float64), nil
	}
	
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		
		if info.IsDir() {
			return nil
		}
		
		if shouldSkipFile(info.Name(), filePath) {
			return nil
		}
		
		content, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return nil
		}
		
		language := enry.GetLanguage(filePath, content)
		if language == "" {
			return nil
		}
		
		if languages == nil {
			languages = make(map[string]int)
		}
		
		languages[language] += len(content)
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(languages) == 0 {
		return make(map[string]float64), nil
	}

	total := 0
	for _, bytes := range languages {
		total += bytes
	}

	stats := make(map[string]float64)
	for language, bytes := range languages {
		if total > 0 {
			percentage := (float64(bytes) / float64(total)) * 100
			stats[language] = percentage
		}
	}

	return stats, nil
}

func GetPrimaryLanguage(stats map[string]float64) string {
	if len(stats) == 0 {
		return ""
	}

	var maxLang string
	var maxPercentage float64

	for lang, percentage := range stats {
		if percentage > maxPercentage {
			maxPercentage = percentage
			maxLang = lang
		}
	}

	return maxLang
}

func DetectFrameworks(componentPath string, primaryLang string) (string, error) {
	switch strings.ToLower(primaryLang) {
	case "javascript", "typescript":
		return detectJSFramework(componentPath)
	case "python":
		return detectPythonFramework(componentPath)
	case "java":
		return detectJavaFramework(componentPath)
	case "go":
		return detectGoFramework(componentPath)
	case "rust":
		return detectRustFramework(componentPath)
	case "c#":
		return detectDotNetFramework(componentPath)
	case "php":
		return detectPHPFramework(componentPath)
	case "ruby":
		return detectRubyFramework(componentPath)
	}

	return "", nil
}

func detectJSFramework(componentPath string) (string, error) {
	packageJsonPath := filepath.Join(componentPath, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return "", nil
	}

	data, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return "", err
	}

	var packageJson struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &packageJson); err != nil {
		return "", err
	}

	allDeps := make(map[string]string)
	for pkg, version := range packageJson.Dependencies {
		allDeps[pkg] = version
	}
	for pkg, version := range packageJson.DevDependencies {
		allDeps[pkg] = version
	}

	return detectFrameworkFromDependencies(allDeps), nil
}

func shouldSkipFile(fileName, filePath string) bool {
	if strings.HasPrefix(fileName, ".") && fileName != ".env" && !strings.HasPrefix(fileName, ".env.") {
		return true
	}
	
	skipExtensions := []string{
		".exe", ".dll", ".so", ".dylib", ".a", ".lib",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico",
		".mp4", ".mp3", ".wav", ".avi", ".mov",
		".zip", ".tar", ".gz", ".7z", ".rar",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	}
	
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, skipExt := range skipExtensions {
		if ext == skipExt {
			return true
		}
	}
	
	if strings.Contains(filePath, "node_modules") ||
		strings.Contains(filePath, ".git") ||
		strings.Contains(filePath, "__pycache__") ||
		strings.Contains(filePath, ".venv") ||
		strings.Contains(filePath, "venv") ||
		strings.Contains(filePath, "target") ||
		strings.Contains(filePath, "build") ||
		strings.Contains(filePath, "dist") {
		return true
	}
	
	return false
}

func detectPythonFramework(componentPath string) (string, error) {
	requirementsPath := filepath.Join(componentPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err != nil {
			return "", err
		}

		deps := make(map[string]string)
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, "==")
			if len(parts) > 0 {
				deps[strings.TrimSpace(parts[0])] = ""
			}
		}

		return detectFrameworkFromDependencies(deps), nil
	}

	pyprojectPath := filepath.Join(componentPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		return "", nil
	}

	return "", nil
}

func detectJavaFramework(componentPath string) (string, error) {
	pomPath := filepath.Join(componentPath, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		content, err := os.ReadFile(pomPath)
		if err != nil {
			return "", err
		}

		contentStr := string(content)
		for framework, patterns := range frameworkPatterns {
			for _, pattern := range patterns {
				if strings.Contains(contentStr, pattern) {
					return framework, nil
				}
			}
		}
	}

	return "", nil
}

func detectGoFramework(componentPath string) (string, error) {
	goModPath := filepath.Join(componentPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err != nil {
			return "", err
		}

		contentStr := string(content)
		for framework, patterns := range frameworkPatterns {
			for _, pattern := range patterns {
				if strings.Contains(contentStr, pattern) {
					return framework, nil
				}
			}
		}
	}

	return "", nil
}

func detectRustFramework(componentPath string) (string, error) {
	cargoPath := filepath.Join(componentPath, "Cargo.toml")
	if _, err := os.Stat(cargoPath); err == nil {
		content, err := os.ReadFile(cargoPath)
		if err != nil {
			return "", err
		}

		contentStr := string(content)
		for framework, patterns := range frameworkPatterns {
			for _, pattern := range patterns {
				if strings.Contains(contentStr, pattern) {
					return framework, nil
				}
			}
		}
	}

	return "", nil
}

func detectDotNetFramework(componentPath string) (string, error) {
	pattern := filepath.Join(componentPath, "*.csproj")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	for _, match := range matches {
		content, err := os.ReadFile(match)
		if err != nil {
			continue
		}

		contentStr := string(content)
		for framework, patterns := range frameworkPatterns {
			for _, pattern := range patterns {
				if strings.Contains(contentStr, pattern) {
					return framework, nil
				}
			}
		}
	}

	return "", nil
}

func detectPHPFramework(componentPath string) (string, error) {
	composerPath := filepath.Join(componentPath, "composer.json")
	if _, err := os.Stat(composerPath); err == nil {
		data, err := os.ReadFile(composerPath)
		if err != nil {
			return "", err
		}

		var composer struct {
			Require    map[string]string `json:"require"`
			RequireDev map[string]string `json:"require-dev"`
		}

		if err := json.Unmarshal(data, &composer); err != nil {
			return "", err
		}

		allDeps := make(map[string]string)
		for pkg, version := range composer.Require {
			allDeps[pkg] = version
		}
		for pkg, version := range composer.RequireDev {
			allDeps[pkg] = version
		}

		return detectFrameworkFromDependencies(allDeps), nil
	}

	return "", nil
}

func detectRubyFramework(componentPath string) (string, error) {
	gemfilePath := filepath.Join(componentPath, "Gemfile")
	if _, err := os.Stat(gemfilePath); err == nil {
		content, err := os.ReadFile(gemfilePath)
		if err != nil {
			return "", err
		}

		contentStr := string(content)
		for framework, patterns := range frameworkPatterns {
			for _, pattern := range patterns {
				if strings.Contains(contentStr, pattern) {
					return framework, nil
				}
			}
		}
	}

	return "", nil
}

func detectFrameworkFromDependencies(deps map[string]string) string {
	scores := make(map[string]int)

	for framework, patterns := range frameworkPatterns {
		for _, pattern := range patterns {
			if _, exists := deps[pattern]; exists {
				scores[framework]++
			}
		}
	}

	if len(scores) == 0 {
		return ""
	}

	var frameworks []string
	for framework := range scores {
		frameworks = append(frameworks, framework)
	}

	sort.Slice(frameworks, func(i, j int) bool {
		return scores[frameworks[i]] > scores[frameworks[j]]
	})

	return frameworks[0]
}

func DetectDevelopmentTools(componentPath string) ([]string, error) {
	var tools []string

	packageJsonPath := filepath.Join(componentPath, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		data, err := os.ReadFile(packageJsonPath)
		if err == nil {
			var packageJson struct {
				DevDependencies map[string]string `json:"devDependencies"`
				Scripts         map[string]string `json:"scripts"`
			}

			if json.Unmarshal(data, &packageJson) == nil {
				if _, exists := packageJson.DevDependencies["eslint"]; exists {
					tools = append(tools, "ESLint")
				}
				if _, exists := packageJson.DevDependencies["prettier"]; exists {
					tools = append(tools, "Prettier")
				}
				if _, exists := packageJson.DevDependencies["jest"]; exists {
					tools = append(tools, "Jest")
				}
				if _, exists := packageJson.DevDependencies["typescript"]; exists {
					tools = append(tools, "TypeScript")
				}
			}
		}
	}

	if _, err := os.Stat(filepath.Join(componentPath, "requirements-dev.txt")); err == nil {
		tools = append(tools, "pip-tools")
	}

	if _, err := os.Stat(filepath.Join(componentPath, "pyproject.toml")); err == nil {
		content, err := os.ReadFile(filepath.Join(componentPath, "pyproject.toml"))
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "black") {
				tools = append(tools, "Black")
			}
			if strings.Contains(contentStr, "flake8") {
				tools = append(tools, "Flake8")
			}
			if strings.Contains(contentStr, "pytest") {
				tools = append(tools, "pytest")
			}
		}
	}

	return tools, nil
}