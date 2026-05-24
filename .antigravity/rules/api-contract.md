# API契約 ルール（OpenAPI / 型生成）

## 真実の在りか

API の真実は **`backend/docs/swagger.yaml`**（swaggo が生成）。
- バックは swag アノテーションを書き、`swag init` でこれを生成
- フロントは `openapi-typescript` でこれを `frontend/src/api/types.ts` に変換

両方とも git で管理する。PR レビューで「API がどう変わったか」が diff で見える状態を保つ。

## バック側のフロー

1. `backend/internal/handler/*.go` に echo のハンドラを書き、swag アノテーションで仕様を記述
2. ローカルで `swag init -g cmd/server/main.go -o ./docs` を実行
3. `backend/docs/swagger.yaml` および `docs.go` をコミット
4. CI で `make swagger-check`（再生成 → 差分が出たら失敗）を実行し、コミット忘れを防ぐ

### swag アノテーションの規約

- 全エンドポイントに以下を付ける:
  - `@Summary`（日本語OK、短く）
  - `@Tags`（リソース別: `auth`, `posts`, `likes`, `users`）
  - `@Accept`, `@Produce`（基本は `json`）
  - `@Param`（path / query / body）
  - `@Success`, `@Failure`（複数）
  - `@Router`
  - `@Security BearerAuth`（認証必須エンドポイント）
- リクエスト/レスポンス型は `handler/dto/*.go` に定義。プリミティブを直接書かない（OpenAPI に component として出すため）

## フロント側のフロー

1. backend をビルドして swagger.yaml を最新化
2. `npx openapi-typescript ../backend/docs/swagger.yaml -o src/api/types.ts`（npm script `gen:api` を用意）
3. `src/api/client.ts` で `paths` 型からエンドポイント別の型を導出

### 型導出の例

```ts
import type { paths } from "./types";

export type GetPostsResponse =
  paths["/posts"]["get"]["responses"]["200"]["content"]["application/json"];

export type CreatePostRequest =
  paths["/posts"]["post"]["requestBody"]["content"]["application/json"];
```

## バージョニング

- URL に `/api/v1` を含める（破壊的変更時は `/api/v2`）
- MVP の間は v1 のみ
- swaggo の info セクションに `@version 1.0.0` を明記

## エラーレスポンス

- 全エラーは共通形式:
  ```json
  { "error": { "code": "VALIDATION_ERROR", "message": "..." } }
  ```
- handler/dto/error.go に `ErrorResponse` を定義し、全 `@Failure` で参照
- code は SPEC.md「8. API」の表に従う

## ページネーション

- カーソルベース。レスポンスに `next_cursor`（nullable string）を含める
- カーソルは `created_at|id` を Base64 エンコード（実装詳細はバックの責任、フロントは opaque 文字列として扱う）

## ヘッダ

- 認証: `Authorization: Bearer <jwt>`
- リクエスト ID: バックが `X-Request-Id` を生成・付与（クライアントから来てれば優先）
- CORS: 許可オリジンを `CORS_ALLOWED_ORIGINS` で設定

## CI チェック

- `swag init` の出力差分検出
- `openapi-typescript` 再実行で `types.ts` の差分検出（コミット忘れ防止）

## 禁止事項

- swagger.yaml を **手書きしない**（必ず swag init から生成）
- フロントで API レスポンス型を**手書きしない**（openapi-typescript 生成物のみ）
- 後方互換性を破る変更を、バージョン更新なしで入れない
