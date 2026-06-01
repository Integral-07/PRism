# PRism

> AIが量産するPRを、人間がサクッと判断できるトリアージツール

---

## 背景・課題

AI coding toolsの普及により、PRの量が爆発的に増加している。

- AIが書いたPRは人間の **1.7倍** 多く問題を含む
- クリティカルな問題は **1.4倍**、メジャーな問題は **1.7倍** 多い
- ロジック・正確性の問題が **75%増**、可読性問題が **3倍以上**
- 開発者の **59%** が、完全に理解していないAI生成コードをそのまま使っている
- AI導入後、PRあたりのインシデントが **23.5%増**、変更失敗率が **30%増**

レビュアーが直面している本質的な問いは3つ。

```
1. このPRを今見るべきか？   → 優先度の判断
2. どこから見ればいいか？   → 着手順序の判断
3. 何に気をつければいいか？ → リスクの判断
```

既存ツール（CodeRabbit・Qodo・GitHub Copilot Review）は「バグを指摘する」ことに特化している。PRismはその上位レイヤーとして、**レビュアーの認知コストを下げる判断支援**に特化する。

---

## コンセプト

**バグを見つけるのではなく、レビュアーが判断するために必要な情報を構造化する。**

CodeRabbitなどの既存ツールと競合せず、補完する立ち位置。

### 競合との差別化

| | CodeRabbit | Qodo | GitHub Copilot Review | **PRism** |
|---|---|---|---|---|
| 目的 | バグ指摘 | レビュー+テスト生成 | 自動レビュー | **トリアージ・判断支援** |
| 出力 | PRコメント | PRコメント | PRコメント | **Checksタブ** |
| GitHubを離れる？ | No | No | No | **No** |
| 既存ツールと | 競合 | 競合 | 競合 | **補完** |

---

## ターゲット

- **主なユーザー**: AIを使って開発しているチームのレビュアー
- **主なペイン**: 大量のPRをさばけない・どこから見ればいいかわからない

---

## 機能

### トリアージ系

- **優先度スコア**: 1〜5のスコアでPRの緊急度を示す
- **リスクレベル**: high / medium / low の3段階
- **レビュー推定時間**: このPRに何分かかるかの目安
- **ファイル優先順位リスト**: どのファイルから見るべきかを順位付け

### 判断支援系

- **重点レビュー箇所**: 特に注意すべきポイントを自然言語で提示
- **破壊的変更の検出**: 後方互換性がない変更を検出
- **カバレッジ低下の検出**: テストカバレッジの変化を検出

### カスタム分析

リポジトリの `.prism.yml` に追加指示を書くと、そのPRに特化した分析が **💬 カスタム分析** セクションに表示される。

---

## 動作イメージ

PRが開くと、数秒後に Checks タブに結果が表示される。

```
🔴 high · 優先度 5/5 · 推定35min

JWTによる認証基盤を全面刷新。既存セッションとの互換性なし。

📁 ファイル優先順位

| ファイル       | リスク     | カテゴリ | メモ                              |
|----------------|------------|----------|-----------------------------------|
| auth/jwt.go    | 🔴 high    | logic    | セッション管理をJWT検証に切り替え。移行戦略の確認が必要。 |

📋 重点レビュー箇所

- 既存セッションとの後方互換性がない
- トークンのリフレッシュロジックに競合状態の可能性

⚠️ 破壊的変更

- auth/jwt.go: セッション管理の全面刷新
```

---

## .prism.yml による出力カスタマイズ

リポジトリのルートに `.prism.yml` を置くと、表示する項目を制御できる。ファイルがない場合はすべての項目が有効になる。

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

---

## アーキテクチャ

クリーンアーキテクチャを採用し、責任を明確に分離。

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

---

## 開発

```bash
# テスト
go test ./...

# ビルド
go build -o bin/prism ./cmd/server
```

CI は GitHub Actions で自動実行（push / PR のたびに `go test ./...` が走る）。

モックを使ったローカル確認は `infrastructure/mock` を参照。

セットアップ手順は [SETUP.md](SETUP.md) を参照。

---

## ライセンス

MIT
