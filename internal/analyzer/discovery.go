package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/replyzer/analyze-repo/internal/types"
)

var configFiles = []string{
	"package.json",
	"requirements.txt",
	"pyproject.toml",
	"pom.xml",
	"build.gradle",
	"Cargo.toml",
	"go.mod",
	"*.csproj",
	"docker-compose.yml",
	"docker-compose.yaml",
}

func DiscoverProjectStructure(repoPath string) (*types.ProjectStructure, error) {
	components := make(map[string]*types.ComponentInfo)
	
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		fileName := info.Name()
		if isConfigFile(fileName) {
			dir := filepath.Dir(path)
			relDir, err := filepath.Rel(repoPath, dir)
			if err != nil {
				return err
			}

			if relDir == "." {
				relDir = ""
			}

			componentName := getComponentName(dir, relDir)
			
			if comp, exists := components[componentName]; exists {
				comp.ConfigFiles = append(comp.ConfigFiles, fileName)
			} else {
				components[componentName] = &types.ComponentInfo{
					Name:         componentName,
					Path:         dir,
					RelativePath: relDir,
					ConfigFiles:  []string{fileName},
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	componentList := make([]types.ComponentInfo, 0, len(components))
	for _, comp := range components {
		componentList = append(componentList, *comp)
	}

	repoType := "single"
	if len(componentList) > 1 {
		repoType = "monorepo"
	}

	return &types.ProjectStructure{
		Type:       repoType,
		Components: componentList,
	}, nil
}

func isConfigFile(fileName string) bool {
	for _, pattern := range configFiles {
		if pattern == fileName {
			return true
		}
		if strings.Contains(pattern, "*") {
			matched, err := filepath.Match(pattern, fileName)
			if err == nil && matched {
				return true
			}
		}
	}
	return false
}

func shouldSkipDir(dirName string) bool {
	skipDirs := []string{
		".git",
		"node_modules",
		".venv",
		"venv",
		"__pycache__",
		".pytest_cache",
		"target",
		"build",
		"dist",
		".next",
		".nuxt",
		"coverage",
		".coverage",
		"bin",
		"obj",
	}

	for _, skip := range skipDirs {
		if dirName == skip {
			return true
		}
	}

	return strings.HasPrefix(dirName, ".")
}

func getComponentName(fullPath, relativePath string) string {
	if relativePath == "" {
		return filepath.Base(fullPath)
	}

	parts := strings.Split(relativePath, string(filepath.Separator))
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return filepath.Base(fullPath)
}