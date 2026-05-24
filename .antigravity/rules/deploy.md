# Deploy ルール（GitHub Actions / 本番）

## 構成

- バックエンド: **Cloud Run** または **Render**（MVP段階で最終決定。両者は無料枠で軽量に運用可能）
- フロントエンド: **Firebase Hosting**
- DB: 本番は **Cloud SQL for PostgreSQL** または Render PostgreSQL
- CI/CD: **GitHub Actions**
- 環境: MVP は **prod 1つのみ**

## ブランチ戦略

- `main` を保護ブランチに（PR 必須、CI 通過必須）
- 機能ブランチ → PR → main へマージ → 自動デプロイ

## ワークフロー

### PR ワークフロー（`.github/workflows/ci.yml`）

トリガ: PR 作成 / 更新

ジョブ:

1. **backend-lint-test**
   - `go vet`, `golangci-lint run`
   - `swag init` 差分チェック（再生成して `git diff` がないか）
   - `go test ./...`（unit）
   - `go test -tags=integration ./...`（integration、testcontainers-go）

2. **frontend-lint-test**
   - `npm ci`
   - `npm run typecheck`（tsc --noEmit）
   - `npm run lint`（ESLint）
   - `openapi-typescript` 再実行 + 差分チェック
   - `npm test`（Vitest）

### main マージ時のデプロイ（`.github/workflows/deploy.yml`）

トリガ: `push` to `main`

ジョブ（順序）:

1. **backend-deploy**
   - Docker イメージビルド（マルチアーキ不要、amd64 で OK）
   - Cloud Run / Render にデプロイ
   - デプロイ後に `migrate up` を実行（Cloud Run のジョブ or デプロイコンテナ内）
2. **frontend-deploy**
   - `npm ci && npm run build`
   - Firebase Hosting にデプロイ（`firebase deploy --only hosting`）
   - VITE 環境変数はビルド時に注入

ジョブ間に依存関係を持たせる場合、`needs:` で順序を明示。

## シークレット管理

GitHub Actions Secrets:

| 名前 | 用途 |
| --- | --- |
| `GCP_SA_KEY` | Cloud Run / Cloud SQL / Firebase Hosting 用サービスアカウント |
| `FIREBASE_TOKEN` | firebase CLI 用（Hosting デプロイ） |
| `JWT_SECRET_PROD` | 本番 JWT シークレット |
| `DATABASE_URL_PROD` | 本番 DB 接続文字列 |
| `FIREBASE_ADMIN_CREDENTIALS_JSON` | バックエンドの Firebase Admin 認証 |

本番ランタイムは Secret Manager 連携が望ましいが、MVP は GitHub Secrets → 環境変数注入で十分。

## マイグレーションの安全策

- ロールバック可能であることが望ましい（破壊的変更は2段階デプロイ）
- 大量データに対する変更（インデックス追加など）は本番で時間を計測してから実施
- マイグレーション中に古いコードが動いてもクラッシュしないよう、「カラム追加 → デプロイ → カラム使用するデプロイ」の順を守る

## ロールバック

- バック: 直前の Cloud Run リビジョンに切り戻し（`gcloud run services update-traffic` 等）
- フロント: Firebase Hosting の rollback コマンドで前リリースに戻す
- DB: マイグレーションの down は best-effort。データ整合性が壊れる場合は手作業で復旧

## モニタリング

- Cloud Run / Render の標準ログ・メトリクス
- バックは `/healthz` を返すエンドポイントを必ず維持
- アラートはまず設定せず、Issue 起票駆動で運用

## 禁止事項

- 本番デプロイを **手動 push で行わない**（必ず main 経由）
- main に直接 commit / push しない（PR 必須）
- シークレットを `.env` のままリポジトリにコミットしない
- 本番 DB に対する任意の手動マイグレーションを GitHub Actions 外から実行しない
