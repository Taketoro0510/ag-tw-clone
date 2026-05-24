# ANTIGRAVITY.md

このプロジェクトは Twitterライクな投稿・いいねを行うSNSアプリ（MVP）です。
詳細な仕様は **[../docs/SPEC.md](../docs/SPEC.md) を一次資料として参照** してください。
本ファイルと `.antigravity/rules/*.md` は SPEC.md から派生した、実装時のガイドラインです。

## プロジェクト概要

- **目的**: ユーザーがテキスト・画像・動画を投稿し、互いに「いいね」を押せる SNS の MVP を作る
- **MVPスコープ**: Googleログイン / 投稿（テキスト140文字 + 画像1枚 or 動画1本30秒） / いいね / グローバルタイムライン / プロフィール画面 / Light・Dark テーマ
- **MVPで含めないもの**: フォロー、リプライ、リポスト、通知、検索、i18n、PWA、サムネ自動生成

## 技術スタック

| 領域 | 採用技術 |
| --- | --- |
| フロントエンド | React + TypeScript + Vite + MUI + react-router + @tanstack/react-query + openapi-typescript |
| バックエンド | Go + echo + gorm + swaggo/echo-swagger + log/slog |
| データベース | PostgreSQL 16 |
| 認証 | Firebase Authentication（Google） + 独自JWT（HS256） |
| ファイルストレージ | Firebase Storage（クライアント直アップロード） |
| ローカル開発 | docker-compose（DB + backend）+ Firebase Local Emulator Suite + frontend は Vite dev server |
| バックエンドのホットリロード | air（go run でのフルコンパイルは禁止） |
| マイグレーション | golang-migrate（SQLファイル、git管理） |
| テスト | バック: Go unit + testcontainers-go、フロント: Vitest + Playwright（E2Eは主要パスのみ） |
| CI/CD | GitHub Actions、main push で本番デプロイ |
| デプロイ先 | バック: Cloud Run または Render、フロント: Firebase Hosting |

## ディレクトリ構成（想定）

```
.
├── .antigravity/
│   ├── ANTIGRAVITY.md       ← 本ファイル
│   └── rules/               ← 領域別ルール（後述）
├── docs/
│   └── SPEC.md              ← 仕様書（一次資料）
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/ repository/ usecase/ handler/ middleware/
│   │   ├── auth/ storage/ config/
│   ├── docs/                ← swag init で生成（git管理）
│   ├── migrations/          ← golang-migrate
│   ├── .air.toml
│   ├── Dockerfile.dev
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/ components/ features/ pages/ hooks/ theme/
│   └── vite.config.ts
├── docker-compose.yml
└── firebase.json            ← エミュレータ設定
```

## 開発の基本フロー

1. **API を変える** → `backend/internal/handler/*.go` の swag アノテーションを更新 → `swag init` で `backend/docs/swagger.yaml` 再生成 → コミット
2. **フロントの型を更新** → `openapi-typescript backend/docs/swagger.yaml -o frontend/src/api/types.ts`
3. **DBスキーマを変える** → `backend/migrations/` に新しいファイル（up/down）を追加 → 起動時または CI で `migrate up`
4. **テスト** → バックは `go test ./...`、フロントは `npm test`、必要に応じて `npx playwright test`
5. **PR** → CI の lint/typecheck/test が通ったらマージ、main で自動デプロイ

## 守ってほしい原則

- **仕様の変更は SPEC.md を更新してから実装する**（コードが先行しない）
- **ホストOSでgoコマンドを直接実行しないでください。すべてのgoコマンドは `docker compose exec api` の中で実行してください。**
- **バックエンドのホットリロードには必ず air を使う**。`go run` でのフルコンパイルループは禁止
- **API の真実は swagger.yaml**。フロントの fetch は openapi-typescript の型を経由して書く
- **シークレットや本番のサービスアカウントJSONをコミットしない**
- **ソフトデリート対象（posts）に対するクエリは必ず `deleted_at IS NULL` を条件に入れる**
- **メディアはクライアントから Firebase Storage に直接アップロードする**。バックエンド経由のmultipart送受信は実装しない
- **新しい依存ライブラリを入れる前に、本当に必要か考える**。MVPは最小構成を維持
- **swaggerドキュメントを生成・更新したら必ずopenapi.yamlへの変換も実行すること**

## 領域別ルール

実装時は対応する rule ファイルも参照すること。

- [rules/backend.md](rules/backend.md) — Go / echo / gorm のコーディング規約とディレクトリ運用
- [rules/frontend.md](rules/frontend.md) — React / TS / MUI のコーディング規約と状態管理
- [rules/database.md](rules/database.md) — PostgreSQL スキーマ・マイグレーションのルール
- [rules/auth.md](rules/auth.md) — Firebase Auth + 独自JWT の認証フロー
- [rules/storage.md](rules/storage.md) — Firebase Storage の使い方とライフサイクル
- [rules/api-contract.md](rules/api-contract.md) — OpenAPI 連携と型生成の運用
- [rules/testing.md](rules/testing.md) — テストの粒度と運用
- [rules/docker.md](rules/docker.md) — ローカル開発の Docker 構成、air の使い方
- [rules/deploy.md](rules/deploy.md) — GitHub Actions と本番デプロイ
- [rules/coding-style.md](rules/coding-style.md) — 共通のコーディングスタイルとレビュー観点

## 参照ドキュメント

- [仕様書 docs/SPEC.md](../docs/SPEC.md) — 機能要件・非機能要件・データモデルの一次資料
