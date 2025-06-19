# 開発ワークフロー

## ブランチ戦略

このプロジェクトではGitHub Flowを採用しています：

- `master`ブランチ：本番環境対応の安定版
- 機能ブランチ：各機能開発用の一時的なブランチ

## 開発手順

### 1. 新機能開発の開始

```bash
# masterブランチから最新を取得
git checkout master
git pull origin master

# 機能ブランチを作成
git checkout -b feature/機能名
# 例: git checkout -b feature/add-docker-support
```

### 2. 開発作業

```bash
# コードを変更・テスト
# ...

# 変更をコミット
git add .
git commit -m "feat: 機能の説明"

# 定期的にpush
git push origin feature/機能名
```

### 3. プルリクエスト作成

```bash
# GitHubでプルリクエストを作成
gh pr create --title "機能タイトル" --body "機能の説明"

# または、GitHubのWebインターフェースで作成
```

### 4. コードレビュー・マージ

- CIが通ることを確認
- コードレビューを受ける
- 承認後、GitHubでマージ

### 5. 後処理

```bash
# masterに戻って最新を取得
git checkout master
git pull origin master

# 機能ブランチを削除
git branch -d feature/機能名
```

## コミットメッセージ規約

[Conventional Commits](https://www.conventionalcommits.org/)に従います：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 主要なtype：
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント更新
- `style`: コードフォーマット変更
- `refactor`: リファクタリング
- `test`: テスト追加・修正
- `chore`: その他の変更

### 例：
```bash
git commit -m "feat: add Docker Compose support for development environment"
git commit -m "fix: resolve language detection issue for Go modules"
git commit -m "docs: update installation instructions in README"
```

## CIチェック

プルリクエストでは以下がチェックされます：

- **テスト**: `go test ./...`
- **ビルド**: `go build ./...`
- **静的解析**: `golangci-lint`

すべてのチェックがパスしてからマージしてください。

## ローカル開発環境

### 必要なツール

```bash
# 開発依存関係のインストール
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# テスト実行
go test ./...

# リンタ実行
golangci-lint run

# ビルド
go build -o bin/analyze-repo cmd/analyze-repo/main.go
```

### プリコミットチェック

コミット前に以下を実行することを推奨：

```bash
# すべてのテストを実行
go test ./...

# 静的解析
golangci-lint run

# フォーマット
go fmt ./...

# ビルド確認
go build ./...
```

## リリースプロセス

1. `master`ブランチで`v1.0.0`形式のタグを作成
2. GitHub Actionsが自動でリリースビルドを実行
3. GitHub Releasesに成果物が公開される

```bash
git tag v1.0.0
git push origin v1.0.0
```