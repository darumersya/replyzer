package analyzer

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ExtractVersionRequirements(componentPath string) (map[string]string, error) {
	requirements := make(map[string]string)

	if err := extractNodeVersions(componentPath, requirements); err != nil {
		return requirements, err
	}

	if err := extractPythonVersions(componentPath, requirements); err != nil {
		return requirements, err
	}

	if err := extractJavaVersions(componentPath, requirements); err != nil {
		return requirements, err
	}

	if err := extractGoVersions(componentPath, requirements); err != nil {
		return requirements, err
	}

	if err := extractRustVersions(componentPath, requirements); err != nil {
		return requirements, err
	}

	if err := extractDotNetVersions(componentPath, requirements); err != nil {
		return requirements, err
	}

	return requirements, nil
}

func extractNodeVersions(componentPath string, requirements map[string]string) error {
	packageJsonPath := filepath.Join(componentPath, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		data, err := os.ReadFile(packageJsonPath)
		if err != nil {
			return err
		}

		var packageJson struct {
			Engines map[string]string `json:"engines"`
		}

		if err := json.Unmarshal(data, &packageJson); err == nil {
			if nodeVersion, exists := packageJson.Engines["node"]; exists {
				requirements["node"] = nodeVersion
			}
			if npmVersion, exists := packageJson.Engines["npm"]; exists {
				requirements["npm"] = npmVersion
			}
		}
	}

	nvmrcPath := filepath.Join(componentPath, ".nvmrc")
	if _, err := os.Stat(nvmrcPath); err == nil {
		content, err := os.ReadFile(nvmrcPath)
		if err == nil {
			version := strings.TrimSpace(string(content))
			if version != "" {
				requirements["node"] = version
			}
		}
	}

	return nil
}

func extractPythonVersions(componentPath string, requirements map[string]string) error {
	pyprojectPath := filepath.Join(componentPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`requires-python\s*=\s*"([^"]+)"`)
		matches := re.FindStringSubmatch(string(content))
		if len(matches) > 1 {
			requirements["python"] = matches[1]
		}
	}

	pythonVersionPath := filepath.Join(componentPath, ".python-version")
	if _, err := os.Stat(pythonVersionPath); err == nil {
		content, err := os.ReadFile(pythonVersionPath)
		if err == nil {
			version := strings.TrimSpace(string(content))
			if version != "" {
				requirements["python"] = version
			}
		}
	}

	return nil
}

func extractJavaVersions(componentPath string, requirements map[string]string) error {
	pomPath := filepath.Join(componentPath, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		content, err := os.ReadFile(pomPath)
		if err != nil {
			return err
		}

		var pom struct {
			Properties struct {
				MavenCompilerSource string `xml:"maven.compiler.source"`
				MavenCompilerTarget string `xml:"maven.compiler.target"`
				JavaVersion         string `xml:"java.version"`
			} `xml:"properties"`
		}

		if err := xml.Unmarshal(content, &pom); err == nil {
			if pom.Properties.MavenCompilerSource != "" {
				requirements["java"] = pom.Properties.MavenCompilerSource
			} else if pom.Properties.JavaVersion != "" {
				requirements["java"] = pom.Properties.JavaVersion
			}
		}
	}

	gradlePath := filepath.Join(componentPath, "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		content, err := os.ReadFile(gradlePath)
		if err == nil {
			contentStr := string(content)
			
			re := regexp.MustCompile(`sourceCompatibility\s*=\s*['"]([\d.]+)['"]`)
			matches := re.FindStringSubmatch(contentStr)
			if len(matches) > 1 {
				requirements["java"] = matches[1]
			}
		}
	}

	return nil
}

func extractGoVersions(componentPath string, requirements map[string]string) error {
	goModPath := filepath.Join(componentPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "go ") {
				version := strings.TrimSpace(strings.TrimPrefix(line, "go"))
				if version != "" {
					requirements["go"] = version
				}
				break
			}
		}
	}

	return nil
}

func extractRustVersions(componentPath string, requirements map[string]string) error {
	cargoPath := filepath.Join(componentPath, "Cargo.toml")
	if _, err := os.Stat(cargoPath); err == nil {
		content, err := os.ReadFile(cargoPath)
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`rust-version\s*=\s*"([^"]+)"`)
		matches := re.FindStringSubmatch(string(content))
		if len(matches) > 1 {
			requirements["rust"] = matches[1]
		}
	}

	return nil
}

func extractDotNetVersions(componentPath string, requirements map[string]string) error {
	pattern := filepath.Join(componentPath, "*.csproj")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		content, err := os.ReadFile(match)
		if err != nil {
			continue
		}

		var project struct {
			PropertyGroup struct {
				TargetFramework string `xml:"TargetFramework"`
			} `xml:"PropertyGroup"`
		}

		if err := xml.Unmarshal(content, &project); err == nil {
			if project.PropertyGroup.TargetFramework != "" {
				requirements["dotnet"] = project.PropertyGroup.TargetFramework
				break
			}
		}
	}

	globalJsonPath := filepath.Join(componentPath, "global.json")
	if _, err := os.Stat(globalJsonPath); err == nil {
		data, err := os.ReadFile(globalJsonPath)
		if err == nil {
			var globalJson struct {
				SDK struct {
					Version string `json:"version"`
				} `json:"sdk"`
			}

			if err := json.Unmarshal(data, &globalJson); err == nil {
				if globalJson.SDK.Version != "" {
					requirements["dotnet-sdk"] = globalJson.SDK.Version
				}
			}
		}
	}

	return nil
}