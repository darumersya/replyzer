package types

type AnalysisResult struct {
	Repository Repository  `yaml:"repository" json:"repository"`
	Components []Component `yaml:"components" json:"components"`
}

type Repository struct {
	Type string `yaml:"type" json:"type"` // "single" | "monorepo"
	Path string `yaml:"path" json:"path"`
	Name string `yaml:"name" json:"name"`
}

type Component struct {
	Name                 string                 `yaml:"name" json:"name"`
	Path                 string                 `yaml:"path" json:"path"`
	Type                 string                 `yaml:"type" json:"type"`
	PrimaryLanguage      string                 `yaml:"primary_language" json:"primary_language"`
	LanguageStats        map[string]float64     `yaml:"language_stats" json:"language_stats"`
	Framework            string                 `yaml:"framework" json:"framework"`
	VersionRequirements  map[string]string      `yaml:"version_requirements" json:"version_requirements"`
	ExternalDependencies ExternalDependencies   `yaml:"external_dependencies" json:"external_dependencies"`
	DevelopmentTools     []string               `yaml:"development_tools" json:"development_tools"`
}

type ExternalDependencies struct {
	Databases []string `yaml:"databases" json:"databases"`
	Services  []string `yaml:"services" json:"services"`
}

type AnalysisOptions struct {
	Format    string
	Output    string
	Verbose   bool
	Component string
	Exclude   []string
}

type ProjectStructure struct {
	Type       string
	Components []ComponentInfo
}

type ComponentInfo struct {
	Name         string
	Path         string
	ConfigFiles  []string
	RelativePath string
}