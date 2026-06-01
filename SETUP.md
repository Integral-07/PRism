# セットアップ

## 前提条件

- Go 1.23以上
- [ngrok](https://ngrok.com/) (ローカル動作確認用)
- GitHub App の作成 (下記手順参照)
- [Gemini API キー](https://aistudio.google.com/app/apikey)

---

## 1. GitHub App を作成する

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

## 2. 環境変数を設定する

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

## 3. ngrok でトンネルを開く

```bash
ngrok http 8080
```

表示された URL (例: `https://xxxx-xx-xx-xx-xx.ngrok-free.app`) をコピーする。

GitHub App の設定ページ → **Webhook URL** を `https://xxxx.ngrok-free.app/webhook` に更新して保存する。

---

## 4. サーバーを起動する

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

## 5. 動作確認

GitHub App をインストールしたリポジトリで **PR を作成または更新** する。

数秒後、PR の **Checks タブ** に `PRism` の結果が表示される。

---

## デプロイ (Render)

1. Render で新しい Web Service を作成
2. GitHub リポジトリを接続
3. 環境変数を設定 (`.env.example` 参照)
4. ビルドコマンド: `go build -o bin/prism ./cmd/server`
5. 起動コマンド: `./bin/prism`
