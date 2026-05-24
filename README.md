# PRism

> AIが量産するPRを、人間がサクッと判断できるトリアージツール

## 概要

PRismは、GitHub App として動作し、PRが開かれた瞬間にdiffを分析してChecksタブに結果を表示するトリアージツールです。

- **優先度スコア**: 1〜5のスコアでPRの緊急度を示す
- **リスクレベル**: high / medium / low の3段階
- **ファイル優先順位**: どのファイルから見るべきかを順位付け
- **レビュー推定時間**: このPRに何分かかるかの目安

## アーキテクチャ

クリーンアーキテクチャを採用し、責任を明確に分離しています。

```
cmd/server/           # Entry point
internal/
  domain/            # Business logic & entities (no external dependencies)
    entity/          # Core entities
    repository/      # Repository interfaces
    service/         # Domain services
  usecase/           # Application business rules
  infrastructure/    # External systems implementation
    github/          # GitHub API client
    claude/          # Claude API client
    supabase/        # Database client
    config/          # Config file reader
  handler/           # HTTP handlers
```

### 依存関係の流れ

```
handler → usecase → domain ← infrastructure
```

- **Domain層**: ビジネスロジックのみ、外部依存なし
- **Infrastructure層**: 外部API実装、Domain層のインターフェースを実装
- **UseCase層**: アプリケーションロジック、Domainとinfrastructureを組み合わせる
- **Handler層**: HTTPリクエストの受付とレスポンス

## セットアップ

### 前提条件

- Go 1.23以上
- GitHub App の作成
- Gemini API キー
- Supabase プロジェクト

### 1. 依存関係のインストール

```bash
# 依存パッケージを自動ダウンロード
go mod download
```

### 2. 環境変数の設定

```bash
cp .env.example .env
# .envファイルを編集して、各種APIキーを設定
```

### 3. ローカルサーバーの起動

```bash
go run cmd/server/main.go
```

サーバーが起動したら `http://localhost:8080/health` でヘルスチェックができます。

### 4. GitHub Webhookの設定

ngrokなどを使ってローカルサーバーを公開し、GitHub AppのWebhook URLに設定します。

```bash
ngrok http 8080
```

## 設定ファイル（.prism.yml）

リポジトリのルートに`.prism.yml`を配置することで、出力をカスタマイズできます。

```yaml
version: 1

language: ja

output:
  triage:
    priority_score: true
    risk_level: true
    estimated_review_time: true
    file_priority_list: true

  support:
    review_focus: true
    breaking_changes: true
    coverage_drop: true

  extra:
    api_changes: false
    dependency_changes: false

thresholds:
  large_pr_lines: 500
  high_risk_files:
    - "auth/**"
    - "payment/**"
```

詳細は `configs/.prism.yml.example` を参照してください。

## 開発

### テストの実行

```bash
# ドメイン層のテスト
go test ./internal/domain/...

# 全テスト
go test ./...
```

### ビルド

```bash
go build -o bin/prism cmd/server/main.go
```

## デプロイ

Renderへのデプロイを想定しています。

1. Renderで新しいWeb Serviceを作成
2. GitHub リポジトリを接続
3. 環境変数を設定
4. ビルドコマンド: `go build -o bin/prism cmd/server/main.go`
5. 起動コマンド: `./bin/prism`

## ライセンス

MIT
