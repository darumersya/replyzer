# 環境分析ツール開発指示書

## 1. プロジェクト概要

### 1.1 目的
ローカルの既存リポジトリを分析し、開発環境に必要な項目をリストアップするCLIツールを開発する。

### 1.2 開発言語・技術スタック
- **メイン言語**: Go (1.21以上)
- **主要ライブラリ**:
  - `github.com/go-enry/go-enry/v2` (言語検出)
  - `github.com/spf13/cobra` (CLI構築)
  - `github.com/spf13/viper` (設定管理)
  - `gopkg.in/yaml.v3` (YAML処理)
  - `encoding/json` (JSON処理)

### 1.3 配布形式
- シングルバイナリ
- クロスプラットフォーム対応 (Windows, macOS, Linux)

## 2. 機能仕様

### 2.1 コマンドラインインターフェース

```bash
# 基本コマンド
analyze-repo [path]

# オプション
analyze-repo [path] [options]

Options:
  --format string     出力形式 (yaml|json) (default "yaml")
  --output string     出力ファイルパス (default: stdout)
  --verbose          詳細ログ出力
  --component string 特定コンポーネントのみ分析
  --exclude string   除外パターン (glob形式)
  --help            ヘルプ表示
  --version         バージョン表示
```

### 2.2 分析フロー

#### Stage 1: 構造発見
1. **設定ファイル検索**
   - 対象: `package.json`, `requirements.txt`, `pyproject.toml`, `pom.xml`, `build.gradle`, `Cargo.toml`, `go.mod`, `*.csproj`
   - 実装: `filepath.Walk()` で再帰検索

2. **リポジトリタイプ判定**
   - 単一プロジェクト: 設定ファイルがルートのみ
   - モノレポ: 複数の設定ファイルが異なるディレクトリに存在

3. **コンポーネント境界識別**
   - 設定ファイルの存在するディレクトリをコンポーネントとして認識

#### Stage 2: 言語・フレームワーク分析
1. **enryによる言語統計**
   - 全体分析: `enry.GetLanguagesByDirectory(repoRoot)`
   - コンポーネント別分析: 各コンポーネントディレクトリで実行

2. **フレームワーク検出**
   - `package.json`: dependencies/devDependencies解析
   - `requirements.txt`: パッケージ名解析
   - `pom.xml`: dependency解析

#### Stage 3: バージョン・依存関係分析
1. **言語バージョン検出**
   - Node.js: `package.json` engines, `.nvmrc`
   - Python: `pyproject.toml` requires-python, `.python-version`
   - Java: `pom.xml`, `build.gradle`
   - .NET: `*.csproj`, `global.json`

2. **外部サービス検出**
   - Docker Compose解析
   - 環境変数ファイル解析

## 3. データ構造定義

### 3.1 メイン構造体

```go
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
    Name                string                 `yaml:"name" json:"name"`
    Path                string                 `yaml:"path" json:"path"`
    Type                string                 `yaml:"type" json:"type"`
    PrimaryLanguage     string                 `yaml:"primary_language" json:"primary_language"`
    LanguageStats       map[string]float64     `yaml:"language_stats" json:"language_stats"`
    Framework           string                 `yaml:"framework" json:"framework"`
    VersionRequirements map[string]string      `yaml:"version_requirements" json:"version_requirements"`
    ExternalDependencies ExternalDependencies   `yaml:"external_dependencies" json:"external_dependencies"`
    DevelopmentTools    []string               `yaml:"development_tools" json:"development_tools"`
}

type ExternalDependencies struct {
    Databases []string `yaml:"databases" json:"databases"`
    Services  []string `yaml:"services" json:"services"`
}
```

## 4. 実装指針

### 4.1 パッケージ構成

```
cmd/
  analyze-repo/
    main.go                 # エントリーポイント
internal/
  analyzer/
    analyzer.go             # メイン分析ロジック
    discovery.go            # 構造発見
    language.go             # 言語・フレームワーク分析
    version.go              # バージョン検出
    dependency.go           # 依存関係分析
  config/
    config.go               # 設定管理
  output/
    formatter.go            # 出力フォーマット
  types/
    types.go                # データ構造定義
```

### 4.2 重要な実装ポイント

#### 4.2.1 エラーハンドリング
- ファイル読み取りエラーの適切な処理
- 部分的な分析失敗時も可能な限り結果を返す
- ユーザーフレンドリーなエラーメッセージ

#### 4.2.2 パフォーマンス
- goroutineによる並列処理
- 大きなファイルの適切な処理
- メモリ効率的な実装

#### 4.2.3 拡張性
- 新しい言語・フレームワークの追加が容易
- 設定による動作カスタマイズ

### 4.3 主要関数の仕様

```go
// メイン分析関数
func AnalyzeRepository(repoPath string, options *AnalysisOptions) (*AnalysisResult, error)

// 構造発見
func DiscoverProjectStructure(repoPath string) (*ProjectStructure, error)

// コンポーネント分析
func AnalyzeComponent(componentPath string) (*Component, error)

// 言語統計取得
func GetLanguageStats(path string) (map[string]float64, error)

// フレームワーク検出
func DetectFrameworks(componentPath string, primaryLang string) (string, error)

// バージョン要件抽出
func ExtractVersionRequirements(componentPath string) (map[string]string, error)

// 外部依存関係検出
func DetectExternalDependencies(componentPath string) (*ExternalDependencies, error)
```

## 5. 設定ファイル対応

### 5.1 対応ファイル一覧

| 言語/フレームワーク | 設定ファイル | 抽出情報 |
|-------------------|-------------|----------|
| Node.js | package.json | engines, dependencies, scripts |
| Node.js | .nvmrc | Node.jsバージョン |
| Python | requirements.txt | パッケージ一覧 |
| Python | pyproject.toml | requires-python, dependencies |
| Python | .python-version | Pythonバージョン |
| Java | pom.xml | maven.compiler.source, dependencies |
| Java | build.gradle | sourceCompatibility, dependencies |
| .NET | *.csproj | TargetFramework, PackageReference |
| .NET | global.json | SDK version |
| Go | go.mod | Go version, dependencies |
| Rust | Cargo.toml | rust-version, dependencies |
| Docker | docker-compose.yml | services, databases |

### 5.2 フレームワーク検出パターン

```go
var frameworkPatterns = map[string][]string{
    "React":     {"react", "@types/react"},
    "Vue":       {"vue", "@vue/cli"},
    "Angular":   {"@angular/core", "@angular/cli"},
    "Next.js":   {"next"},
    "Express":   {"express"},
    "Django":    {"django", "Django"},
    "FastAPI":   {"fastapi"},
    "Flask":     {"flask", "Flask"},
    "Spring":    {"org.springframework"},
}
```

## 6. テスト仕様

### 6.1 テストカテゴリ
1. **単体テスト**: 各機能モジュールの個別テスト
2. **統合テスト**: 実際のプロジェクトを使った結合テスト
3. **E2Eテスト**: CLI全体の動作テスト

### 6.2 テストデータ
- サンプルプロジェクト（単一、モノレポ）
- 各言語・フレームワークの典型的な構成

## 7. ビルド・配布

### 7.1 ビルド設定

```makefile
# Makefile例
.PHONY: build test clean

build:
	go build -o bin/analyze-repo cmd/analyze-repo/main.go

build-all:
	GOOS=windows GOARCH=amd64 go build -o bin/analyze-repo-windows.exe cmd/analyze-repo/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/analyze-repo-darwin cmd/analyze-repo/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/analyze-repo-linux cmd/analyze-repo/main.go

test:
	go test ./...

clean:
	rm -rf bin/
```

### 7.2 CI/CD
- GitHub Actionsでのビルド・テスト自動化
- リリース時の自動バイナリ生成

## 8. 出力例

### 8.1 単一プロジェクト（React）

```yaml
repository:
  type: "single"
  path: "/path/to/react-app"
  name: "react-app"

components:
  - name: "react-app"
    path: "."
    type: "web-application"
    primary_language: "TypeScript"
    language_stats:
      TypeScript: 85.5
      CSS: 10.2
      HTML: 4.3
    framework: "React"
    version_requirements:
      node: ">=18.0.0"
      react: "^18.2.0"
    external_dependencies:
      databases: []
      services: []
    development_tools:
      - "ESLint"
      - "Prettier"
```

### 8.2 モノレポ（フルスタック）

```yaml
repository:
  type: "monorepo"
  path: "/path/to/fullstack-app"
  name: "fullstack-app"

components:
  - name: "frontend"
    path: "./frontend"
    type: "web-application"
    primary_language: "TypeScript"
    language_stats:
      TypeScript: 90.0
      CSS: 10.0
    framework: "React"
    version_requirements:
      node: ">=18.0.0"
      react: "^18.2.0"
    external_dependencies:
      databases: []
      services: []
    development_tools:
      - "ESLint"
      - "Prettier"

  - name: "backend"
    path: "./backend"
    type: "api-service"
    primary_language: "Python"
    language_stats:
      Python: 85.0
      YAML: 15.0
    framework: "FastAPI"
    version_requirements:
      python: ">=3.11"
      fastapi: "^0.100.0"
    external_dependencies:
      databases: ["PostgreSQL"]
      services: ["Redis"]
    development_tools:
      - "Black"
      - "Flake8"
```

