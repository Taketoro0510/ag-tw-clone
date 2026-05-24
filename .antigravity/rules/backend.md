# Backend ルール（Go / echo / gorm）

## レイヤ構成

```
backend/
  cmd/server/main.go          -- エントリポイント。依存性のワイヤリングのみ
  internal/
    domain/                    -- ドメインモデル（純粋なGo struct + メソッド）
    repository/                -- データアクセス層。interface と gorm 実装を分ける
    usecase/                   -- ビジネスロジック。トランザクション境界はここ
    handler/                   -- echo の HTTP ハンドラ。リクエスト/レスポンス整形のみ
    middleware/                -- 認証、レート制限、ロギング、recover
    auth/                      -- Firebase Admin 検証、独自JWT 発行/検証
    storage/                   -- Firebase Storage 操作（オブジェクト削除など）
    config/                    -- 環境変数読み込み（envconfig 等）
  docs/                        -- swag init 生成物（git 管理）
  migrations/                  -- golang-migrate
```

依存方向: `handler -> usecase -> repository <- domain`。逆向き禁止。

## echo ハンドラの規約

- ハンドラは `usecase` を呼ぶだけにする。SQL / 認証ロジックを直書きしない
- リクエスト/レスポンス DTO は `handler/dto/*.go` に定義する
- swag アノテーションを **必ず** 付ける（OpenAPI が正なので、コードコメントが仕様になる）
- バリデーションは `go-playground/validator` を `echo.Bind` 後に呼ぶ
- エラーは `usecase` 層から `errors.Is` で判別できる sentinel error を返し、ミドルウェアで HTTP コードに変換する

例:

```go
// CreatePost godoc
// @Summary  投稿を作成
// @Tags     posts
// @Accept   json
// @Produce  json
// @Param    body  body      dto.CreatePostRequest  true  "投稿"
// @Success  201   {object}  dto.PostResponse
// @Failure  400   {object}  dto.ErrorResponse
// @Failure  401   {object}  dto.ErrorResponse
// @Security BearerAuth
// @Router   /posts [post]
func (h *PostHandler) Create(c echo.Context) error { ... }
```

## gorm の規約

- `domain` モデルとは別に `repository/model/*.go` に gorm 用の struct を置き、変換関数で行き来する（ドメインに DB の都合を漏らさない）
- `Find` / `First` / `Where` でクエリ。生 SQL は使うときは `Raw` で明示
- ソフトデリート対象（posts）は `gorm.DeletedAt` を使う or 明示的に `deleted_at IS NULL` を WHERE に入れる
- トランザクションは `db.Transaction(func(tx *gorm.DB) error {...})` で usecase 層から張る
- N+1 を避ける（`Preload` または手動 join）

## エラーハンドリング

- ドメインエラーは sentinel: `var ErrPostNotFound = errors.New("post not found")`
- usecase はラップせずそのまま返す
- handler / middleware で `errors.Is` で判別 → エラーレスポンス DTO に変換
- レスポンス形式は SPEC.md の「8. API」セクション参照

## ロギング

- `log/slog` を使う。ロガーは `context.Context` 経由で持ち回る（`slog.SetDefault` でグローバルも設定）
- リクエストごとに traceId を発行し、すべてのログに含める（ミドルウェアで `context` に乗せる）
- パスワードや JWT 本体はログに出さない

## レート制限

- ミドルウェアで実装。キーは `userID`（未認証 GET は IP）
- ストレージは MVP では in-memory（`golang.org/x/time/rate`）で OK。水平スケール時は Redis に差し替える前提でインターフェイスを切る
- 制限値は環境変数（`RATE_LIMIT_POSTS_PER_MIN` 等）から読む

## 起動時

- `cmd/server/main.go`:
  1. config 読み込み
  2. DB 接続 + golang-migrate でマイグレーション適用
  3. Firebase Admin SDK 初期化
  4. echo 初期化、ミドルウェア、ルート登録
  5. `:8080` でリッスン

- マイグレーションは起動時 + CI デプロイ時の両方で適用できるよう CLI も用意

## 禁止事項

- `go run` で開発しない。`air` を使う（[rules/docker.md](docker.md) 参照）
- handler から直接 `gorm.DB` を触らない
- ドメインモデルに gorm タグを付けない
- 本番用 Firebase 認証情報を git にコミットしない
