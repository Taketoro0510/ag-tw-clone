# Database ルール（PostgreSQL）

## 基本方針

- DB: PostgreSQL 16（Cloud SQL / Render Postgres / ローカルは docker-compose）
- マイグレーションは **golang-migrate**。`gorm.AutoMigrate` は使わない
- 主キーは **UUID v7**（時系列ソート可能）
- タイムスタンプは `timestamptz`（UTC 保存）
- ソフトデリート対象は `deleted_at timestamptz` を NULL 許容で持つ

## テーブル設計（MVP）

詳細は `docs/SPEC.md` の「7. データモデル」を正として参照。
ここでは運用ルールのみ記す。

### users

- `firebase_uid` に unique index
- `email` は通知 / 将来のための保持。NOT NULL
- `display_name` は Google の `name` を初期値
- `avatar_url` は Google の `picture` を初期値、NULL 可

### posts

- 必要なインデックス:
  - `(created_at DESC, id DESC) WHERE deleted_at IS NULL` — タイムライン
  - `(author_id, created_at DESC) WHERE deleted_at IS NULL` — プロフィール画面の投稿一覧
- `body` は VARCHAR/TEXT（長さ制約はアプリで140まで、DB は寛容に TEXT）
- `media_type` は `CHECK (media_type IN ('image','video'))`、または ENUM
- `media_path`, `media_url` は Storage 上のオブジェクトの位置と公開URL

### likes

- 主キー `(post_id, user_id)` の複合
- `(user_id, created_at DESC)` index（将来「いいねした投稿一覧」用）
- `like_count` は posts に持たない（MVPはCOUNTで集計、スケール時にカウンタ化）

## マイグレーション運用

- ファイル名: `migrations/{timestamp}_{snake_case_summary}.up.sql` / `.down.sql`
- 1マイグレーション = 1論理変更
- スキーマ変更は **PRに対応**。データだけのマイグレーション（バックフィル等）は別ファイルで明示
- `down` はベストエフォートで書く（本番の rollback は注意深く行う）
- 例:
  - `20260524_120000_create_users.up.sql`
  - `20260524_120000_create_users.down.sql`

## クエリの規約

- `SELECT * FROM posts WHERE deleted_at IS NULL` を素で書いていい
- カーソルページネーション:
  ```sql
  SELECT * FROM posts
  WHERE deleted_at IS NULL
    AND (created_at, id) < ($1, $2)   -- cursor 復号値
  ORDER BY created_at DESC, id DESC
  LIMIT $3;
  ```
- ベンチマークなしに新しい index を足さない。本番で必要になったら追加

## トランザクション

- 「いいね作成」「投稿削除 + Storage 削除」など複数操作は usecase 層でトランザクションを張る
- Storage への副作用はトランザクション**外**で実行（DB ロールバック後に Storage を戻せない場合があるため、まず DB をコミット → Storage の削除リクエスト → 失敗時はリトライキュー or ログのみ）

## 接続管理

- gorm の DSN は `DATABASE_URL` 環境変数から
- max open conns: 25（Cloud Run の並列性に合わせる）
- max idle conns: 10
- conn max lifetime: 5分

## 禁止事項

- `gorm.AutoMigrate` を本番で実行しない（マイグレーションは golang-migrate に一本化）
- 本番 DB に対する手動 SQL 実行は原則禁止（必要時は DBA 役と相方レビュー）
- マイグレーションファイルを **git push 後に書き換えない**（不可逆）
