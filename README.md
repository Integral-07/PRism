# PRism

> AIが量産するPRを、人間がサクッと判断できるトリアージツール

## 概要

PRismは、GitHub App として動作し、PRが開かれた瞬間にdiffを分析してChecksタブに結果を表示するトリアージツールです。

- **優先度スコア**: 1〜5のスコアでPRの緊急度を示す
- **リスクレベル**: high / medium / low の3段階
- **レビュー推定時間**: このPRに何分かかるかの目安
- **ファイル優先順位**: 変更ファイルごとのリスクとカテゴリ
- **破壊的変更 / カバレッジ低下**: 影響範囲の検出
- **カスタム分析**: リポジトリごとの追加指示に対応

出力内容はリポジトリに `.prism.yml` を置くことで制御できます。

## アーキテクチャ

クリーンアーキテクチャを採用し、責任を明確に分離しています。

```
cmd/server/           # エントリーポイント
internal/
  domain/             # エンティティ (外部依存なし)
  usecase/            # アプリケーションロジック
    webhook/          # Webhook受信 → 分析 → Check結果投稿
    analyze_pr/       # LLM でdiffを分析
  infrastructure/     # 外部システムの実装
    github/           # GitHub API (diff取得 / Check Run作成 / ラベル付与)
    gemini/           # Gemini API (LLM)
    config/           # 環境変数ロード
    mock/             # ローカル開発・デモ用モック
  handler/            # HTTPハンドラー + 署名検証ミドルウェア
```

### リクエストの流れ

```
GitHub Webhook
    ↓ X-Hub-Signature-256 検証
handler/webhook
    ↓
usecase/webhook
    ├─ PRRepository.GetDiff        (GitHub API)
    ├─ ConfigRepository.Get        (.prism.yml 取得)
    ├─ AnalyzerUseCase.Execute     (Gemini API)
    ├─ CheckRepository.PostResult  (GitHub Checks API)
    └─ LabelRepository.SyncLabels  (GitHub Labels API)
```

## .prism.yml による出力カスタマイズ

リポジトリのルートに `.prism.yml` を置くと、表示する項目を制御できます。

```yaml
output:
  triage:
    priority_score: true        # 優先度スコア (1〜5)
    risk_level: true            # リスクレベル (high / medium / low)
    estimated_review_time: true # レビュー推定時間
    file_priority_list: true    # ファイル優先順位テーブル
  support:
    review_focus: true          # 重点レビュー箇所
    breaking_changes: true      # 破壊的変更の検出
    coverage_drop: true         # カバレッジ低下の懸念
custom: |
  追加の分析指示をここに書く。結果は「カスタム分析」セクションに表示される。
```

`.prism.yml` が存在しない場合はすべての項目が有効になります。

## セットアップ

### 前提条件

- Go 1.23以上
- [ngrok](https://ngrok.com/) (ローカル動作確認用)
- GitHub App の作成 (下記手順参照)
- [Gemini API キー](https://aistudio.google.com/app/apikey)

---

### 1. GitHub App を作成する

1. GitHub → Settings → Developer settings → **GitHub Apps** → **New GitHub App**
2. 以下を設定する

| 項目 | 値 |
|---|---|
| GitHub App name | 任意 (例: `prism-local`) |
| Homepage URL | 任意 |
| Webhook URL | 後で設定するので仮に `https://example.com` |
| Webhook secret | 任意の文字列 (後で `.env` に使用) |

3. **Permissions** を設定する

| 権限 | レベル |
|---|---|
| Pull requests | Read |
| Checks | Write |
| Contents | Read |
| Issues | Write |

4. **Where can this GitHub App be installed?** → `Only on this account`

5. 作成後:
   - **App ID** をメモする
   - **Generate a private key** をクリックして PEM ファイルをダウンロード

6. 作成した GitHub App を、テスト用リポジトリに **Install** する

---

### 2. 環境変数を設定する

```bash
cp .env.example .env
```

`.env` を編集する:

```env
GITHUB_APP_ID=<App ID>
GITHUB_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----\n<PEMの中身を1行にして\nで区切る>\n-----END RSA PRIVATE KEY-----"
GITHUB_WEBHOOK_SECRET=<手順1で設定した Webhook secret>
GEMINI_API_KEY=<Gemini API キー>
```

**PEM の整形方法:**

```bash
# ダウンロードした .pem ファイルの改行を \n に変換する
awk 'NF {printf "%s\\n", $0}' your-app.pem
```

出力された1行の文字列を `GITHUB_PRIVATE_KEY=` に貼り付ける。

---

### 3. ngrok でトンネルを開く

```bash
ngrok http 8080
```

表示された URL (例: `https://xxxx-xx-xx-xx-xx.ngrok-free.app`) をコピーする。

GitHub App の設定ページ → **Webhook URL** を `https://xxxx.ngrok-free.app/webhook` に更新して保存する。

---

### 4. サーバーを起動する

```bash
go run ./cmd/server
```

```
server starting on :8080
```

ヘルスチェック:

```bash
curl http://localhost:8080/health
# ok
```

---

### 5. 動作確認

GitHub App をインストールしたリポジトリで **PR を作成または更新** する。

数秒後、PR の **Checks タブ** に `PRism` の結果が表示される:

```
🟡 medium · 優先度 3/5 · 推定15min

READMEの挨拶文を変更。影響範囲は限定的。

📁 ファイル優先順位

| ファイル    | リスク   | カテゴリ | メモ                         |
|-------------|----------|----------|------------------------------|
| README.md   | 🟢 low   | docs     | 挨拶文の変更のみ。ロジックへの影響なし。 |

📋 重点レビュー箇所

- 変更内容は軽微で、レビュー優先度は低い
- 他ファイルへの影響がないか確認

💬 カスタム分析

はい、このPRはREADMEのみの変更です。
```

---

## 開発

### テストの実行

```bash
go test ./...
```

CI は GitHub Actions で自動実行されます（push / PR のたびに `go test ./...` が走ります）。

### ビルド

```bash
go build -o bin/prism ./cmd/server
```

### モックを使ったローカル確認

`infrastructure/mock` に固定レスポンスを返すモック実装があります。Gemini API キーなしでローカル動作確認したい場合に `main.go` で差し替えて使用してください。

## デプロイ

Render へのデプロイを想定しています。

1. Render で新しい Web Service を作成
2. GitHub リポジトリを接続
3. 環境変数を設定 (`.env.example` 参照)
4. ビルドコマンド: `go build -o bin/prism ./cmd/server`
5. 起動コマンド: `./bin/prism`

## ライセンス

MIT
