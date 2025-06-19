package analyzer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/replyzer/analyze-repo/internal/types"
	"gopkg.in/yaml.v3"
)

func DetectExternalDependencies(componentPath string) (*types.ExternalDependencies, error) {
	deps := &types.ExternalDependencies{
		Databases: make([]string, 0),
		Services:  make([]string, 0),
	}

	if err := analyzeDockerCompose(componentPath, deps); err != nil {
		return deps, err
	}

	if err := analyzeEnvironmentFiles(componentPath, deps); err != nil {
		return deps, err
	}

	return deps, nil
}

func analyzeDockerCompose(componentPath string, deps *types.ExternalDependencies) error {
	composeFiles := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"}

	for _, filename := range composeFiles {
		composePath := filepath.Join(componentPath, filename)
		if _, err := os.Stat(composePath); err == nil {
			content, err := os.ReadFile(composePath)
			if err != nil {
				continue
			}

			var compose struct {
				Services map[string]struct {
					Image   string            `yaml:"image"`
					Build   interface{}       `yaml:"build"`
					Environment map[string]string `yaml:"environment"`
					Ports   []string          `yaml:"ports"`
				} `yaml:"services"`
			}

			if err := yaml.Unmarshal(content, &compose); err != nil {
				continue
			}

			for serviceName, service := range compose.Services {
				if service.Image != "" {
					categorizeService(serviceName, service.Image, deps)
				}
			}
		}
	}

	return nil
}

func analyzeEnvironmentFiles(componentPath string, deps *types.ExternalDependencies) error {
	envFiles := []string{".env", ".env.local", ".env.development", ".env.production"}

	for _, filename := range envFiles {
		envPath := filepath.Join(componentPath, filename)
		if _, err := os.Stat(envPath); err == nil {
			content, err := os.ReadFile(envPath)
			if err != nil {
				continue
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				if strings.Contains(line, "DATABASE_URL") || strings.Contains(line, "DB_") {
					extractDatabaseFromEnv(line, deps)
				}

				if strings.Contains(line, "REDIS") || strings.Contains(line, "CACHE") {
					extractServiceFromEnv(line, deps)
				}
			}
		}
	}

	return nil
}

func categorizeService(serviceName, image string, deps *types.ExternalDependencies) {
	imageLower := strings.ToLower(image)

	databases := map[string]string{
		"postgres":    "PostgreSQL",
		"mysql":       "MySQL",
		"mariadb":     "MariaDB", 
		"mongodb":     "MongoDB",
		"redis":       "Redis",
		"elasticsearch": "Elasticsearch",
		"cassandra":   "Cassandra",
		"couchdb":     "CouchDB",
		"neo4j":       "Neo4j",
		"influxdb":    "InfluxDB",
		"timescaledb": "TimescaleDB",
	}

	services := map[string]string{
		"nginx":       "Nginx",
		"apache":      "Apache",
		"traefik":     "Traefik",
		"rabbitmq":    "RabbitMQ",
		"kafka":       "Apache Kafka",
		"zookeeper":   "Apache Zookeeper",
		"memcached":   "Memcached",
		"consul":      "Consul",
		"vault":       "HashiCorp Vault",
		"prometheus":  "Prometheus",
		"grafana":     "Grafana",
		"jaeger":      "Jaeger",
		"zipkin":      "Zipkin",
		"minio":       "MinIO",
	}

	for key, dbName := range databases {
		if strings.Contains(imageLower, key) {
			if !contains(deps.Databases, dbName) {
				deps.Databases = append(deps.Databases, dbName)
			}
			return
		}
	}

	for key, serviceName := range services {
		if strings.Contains(imageLower, key) {
			if !contains(deps.Services, serviceName) {
				deps.Services = append(deps.Services, serviceName)
			}
			return
		}
	}

	if !strings.Contains(imageLower, "/") || strings.Contains(imageLower, "scratch") {
		return
	}

	if !contains(deps.Services, image) {
		deps.Services = append(deps.Services, image)
	}
}

func extractDatabaseFromEnv(line string, deps *types.ExternalDependencies) {
	lineLower := strings.ToLower(line)

	databases := map[string]string{
		"postgres":    "PostgreSQL",
		"mysql":       "MySQL", 
		"mariadb":     "MariaDB",
		"mongodb":     "MongoDB",
		"sqlite":      "SQLite",
		"redis":       "Redis",
		"elasticsearch": "Elasticsearch",
	}

	for key, dbName := range databases {
		if strings.Contains(lineLower, key) {
			if !contains(deps.Databases, dbName) {
				deps.Databases = append(deps.Databases, dbName)
			}
			break
		}
	}
}

func extractServiceFromEnv(line string, deps *types.ExternalDependencies) {
	lineLower := strings.ToLower(line)

	services := map[string]string{
		"redis":     "Redis",
		"memcached": "Memcached",
		"rabbitmq":  "RabbitMQ",
		"kafka":     "Apache Kafka",
	}

	for key, serviceName := range services {
		if strings.Contains(lineLower, key) {
			if !contains(deps.Services, serviceName) {
				deps.Services = append(deps.Services, serviceName)
			}
			break
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}