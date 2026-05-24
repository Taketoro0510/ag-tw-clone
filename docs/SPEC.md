# SNSアプリ 仕様書（v1 / MVP）

Twitterライクな投稿・いいねを行えるSNSアプリのMVP仕様書。
本ドキュメントは最終の正となる要件定義書であり、`.antigravity/ANTIGRAVITY.md` および `.antigravity/rules/*.md` はここから派生する実装ガイドラインとして位置づける。

---

## 1. プロダクト概要

- **コンセプト**: ユーザーがテキスト・画像・動画を投稿し、互いに「いいね」を押せるシンプルなSNS
- **対象ユーザー**: 個人ユーザー（PC / スマートフォン両対応）
- **MVPの方向性**: 「投稿」「いいね」「タイムライン閲覧」を最短経路で成立させる。フォロー、リプライ、リポスト、通知、検索などは v2 以降

## 2. MVPスコープ

### 含めるもの

| 機能 | 概要 |
| --- | --- |
| Google ログイン | Firebase Auth 経由でGoogleアカウントでログイン |
| 投稿（作成） | テキスト最大140文字 + 画像1枚 または 動画1本（最大30秒） |
| 投稿（削除） | 自分の投稿のみ削除可能（ソフトデリート + Storage同期削除） |
| 投稿（編集） | 不可（Twitter準拠） |
| いいね / いいね解除 | 投稿に対していいね、解除可能 |
| グローバルタイムライン | 全ユーザーの投稿を新しい順に表示。カーソルベース無限スクロール |
| プロフィール画面 | `/users/{id}` で表示名 + そのユーザーの投稿一覧 |
| テーマ切替 | Light / Dark の2種（localStorage 保存、OS設定追従もオプション） |
| レスポンシブUI | PC（≧1024px）・タブレット（768〜）・モバイル（〜767px） |
| レート制限 | サーバ側ミドルウェアで abuse 対策 |

### 含めないもの（v2以降）

- フォロー / フォロワー、フォロー中タイムライン
- リプライ、リポスト（リツイート）、引用
- ハッシュタグ、メンション、検索
- 通知（in-app, push）
- DM
- ユーザーのブロック / ミュート、投稿の通報
- 多言語対応（i18n）
- PWA対応（アイコンマニフェスト、Service Worker）
- 動画のサムネイル自動生成
- アバター画像の差し替え（Googleアカウント画像をそのまま使用）
- 自己紹介文（bio）、ユーザー名（@username）

## 3. 用語

| 用語 | 定義 |
| --- | --- |
| ユーザー | 本アプリにログインしているアカウント。Googleアカウントと1対1 |
| 投稿 (Post) | 1人のユーザーが作成する1つのコンテンツ。テキスト + 任意のメディア1点 |
| メディア (Media) | 投稿に紐づく画像または動画。Firebase Storage 上に保存 |
| いいね (Like) | 投稿とユーザーの組。同じ投稿に対して同じユーザーは1つまで |
| タイムライン | 投稿の一覧。MVPではグローバル（全ユーザー）のみ |

## 4. ユーザーストーリー（MVP）

1. ユーザーは Google アカウントで本アプリにサインインできる
2. サインイン済のユーザーはタイムラインで他ユーザーの投稿を新しい順に閲覧できる
3. サインイン済のユーザーはテキスト140文字以内 + 画像1枚 または 動画1本（30秒以内）を投稿できる
4. サインイン済のユーザーは任意の投稿に「いいね」「いいね解除」を行える
5. サインイン済のユーザーは自分の投稿を削除できる
6. ユーザーは `/users/{id}` で他人を含むユーザーのプロフィール（表示名 + 投稿一覧）を閲覧できる
7. ユーザーは Light / Dark テーマを切り替えられる

未ログイン時の挙動は v2 で検討（MVPでは「サインインを要求」する画面を出す）。

## 5. 機能要件

### 5.1 認証

- **方式**: Firebase Authentication（Googleプロバイダ）でフロント側がIDトークンを取得 → バックエンドが Firebase Admin SDK で検証 → 独自JWTを発行
- **独自JWT**: HS256、ペイロードに `sub`（ユーザーID, UUID v7）、`iat`, `exp`（24時間）。Refresh トークンは MVP では使わず、期限切れ時は再度Firebaseで取得した IDトークンで再発行
- **保存**: フロントでは `localStorage` に保存（XSS リスクは認識した上でMVPの簡便さを優先。将来的に httpOnly cookie 化を検討）
- **送信**: API リクエストの `Authorization: Bearer <jwt>` ヘッダで送信

ユーザー初回ログイン時は Firebase の `uid`, `email`, `name`, `picture` を使って `users` レコードを作成（upsert）。

### 5.2 投稿

#### 投稿作成

- フロー:
  1. フロントが Firebase Storage SDK でメディアを直接アップロード（バックは経由しない）
  2. アップロード完了で Storage 上のオブジェクトパスを取得
  3. バックエンドに `POST /api/v1/posts` を投げる（本文 + メディア情報）
- バリデーション:
  - 本文: 0〜140文字（メディアありなら本文0文字も可、メディアなしなら本文1文字以上）
  - メディア:
    - 画像: `image/jpeg`, `image/png`, `image/webp`, `image/gif` のいずれか、5MB以下
    - 動画: `video/mp4`, `video/quicktime` のいずれか、50MB以下、長さ30秒以内（クライアント側で `<video>` の `duration` で確認、サーバ側はファイルメタは信頼せず recordingの長さ情報をオプションパラメータとして受ける）
- 同一ユーザーは1投稿につきメディア最大1点

#### 投稿削除

- 自分の投稿のみ削除可能
- DB は `deleted_at` を更新（ソフトデリート）
- Storage 側のオブジェクトは同期的に削除リクエスト（失敗時はエラーログを残しジョブ再試行可能にする）

### 5.3 いいね

- 1投稿 × 1ユーザーで最大1いいね（複合ユニーク制約）
- フロントには各投稿につき `like_count`（int）と `liked_by_me`（bool）を返す
- いいね一覧画面、いいねしたユーザー名の表示はしない（MVP）

### 5.4 タイムライン

- 全投稿を `created_at DESC, id DESC` の順で取得
- カーソルベースのページネーション:
  - クエリパラメータ: `limit`（default 20, max 50）, `cursor`（直前ページ末尾の `created_at|id` をBase64エンコード）
  - レスポンス: `items[]`, `next_cursor`（次が無いと null）
- ソフトデリート済の投稿は返さない

### 5.5 プロフィール

- `GET /api/v1/users/{id}` でユーザーの基本情報を取得
- `GET /api/v1/users/{id}/posts` でそのユーザーの投稿一覧を取得（タイムラインと同じカーソル仕様）

### 5.6 テーマ

- 切替は画面右上のスイッチで `light` / `dark` を切り替える
- 選択は `localStorage` に保存し再訪時に復元（OS設定追従モードはオプションとして用意）
- MUI の `createTheme` でブランドカラーをカスタマイズ。配色はデザイントークンとして `frontend/src/theme/` 配下に集約

### 5.7 レート制限（初期値）

| エンドポイント区分 | レート | キー |
| --- | --- | --- |
| 投稿作成 `POST /api/v1/posts` | 5件/分 | `userId` 単位 |
| いいね操作 `POST/DELETE /api/v1/posts/{id}/likes` | 30件/分 | `userId` 単位 |
| GET系全般 | 60件/分 | `userId` 単位（未認証時は IP 単位） |

数値はミドルウェアの設定で容易に変更できるよう、コードと環境変数で分離する。

## 6. 非機能要件

| 項目 | 内容 |
| --- | --- |
| パフォーマンス | タイムライン初回ロード LCP 2.5秒以内（モバイル4Gで） |
| 可用性 | 99% / 月（MVPのため緩め） |
| セキュリティ | Firebase Admin SDK でのIDトークン検証必須。CORS は許可ドメインのみ。SQL/XSS/CSRF 対策（gorm のプレースホルダ使用、MUI の標準エスケープ、APIはCookie認証ではないためCSRFリスクは低） |
| ロギング | 構造化ログ（Go: `log/slog`）。HTTP リクエストごとに traceId を発行 |
| エラー監視 | クラウドの標準ログビューア（Cloud Run の Logs Explorer / Render Logs）。Sentry は v2 以降 |
| アクセシビリティ | キーボード操作可、`alt` 属性必須、コントラスト比4.5:1以上 |
| ブラウザ対応 | 直近2バージョンの Chrome / Edge / Safari / Firefox |
| 国際化 | 日本語のみ（i18n ライブラリは入れない） |

## 7. データモデル

PostgreSQL。主キーは UUID v7（`gen_random_uuid()` ではなく v7 ライブラリで生成）。タイムスタンプは `timestamptz`、ソフトデリートは `deleted_at`。

```
users
  id              uuid (PK, v7)
  firebase_uid    text (unique, not null)
  email           text (not null)
  display_name    text (not null)            -- Googleの name を初期値、ユーザー編集可（MVPでは表示のみ）
  avatar_url      text (nullable)            -- Googleの picture を保存
  created_at      timestamptz (not null)
  updated_at      timestamptz (not null)

posts
  id              uuid (PK, v7)
  author_id       uuid (FK -> users.id, not null, indexed)
  body            text (not null, length <= 140 at app level)
  media_type      text (nullable, enum: image|video)
  media_path      text (nullable, Firebase Storage object path)
  media_url       text (nullable, 公開URL or 署名URL)
  created_at      timestamptz (not null, indexed)
  deleted_at      timestamptz (nullable)
  -- index: (created_at DESC, id DESC) WHERE deleted_at IS NULL
  -- index: (author_id, created_at DESC) WHERE deleted_at IS NULL

likes
  post_id         uuid (FK -> posts.id, not null)
  user_id         uuid (FK -> users.id, not null)
  created_at      timestamptz (not null)
  PRIMARY KEY (post_id, user_id)
  -- index: (user_id, created_at DESC)  -- 将来「自分のいいね一覧」用
```

`posts.like_count` は集計用キャッシュとして持たない（MVPでは都度 `COUNT` で良い。スケール時にカウンタカラム + トランザクション加算へ移行）。

## 8. API（抜粋）

完全な定義は backend の swaggo アノテーションから生成される `backend/docs/swagger.yaml` を正とする。フロントは `openapi-typescript` で型生成。

ベースURL: `/api/v1`

| メソッド | パス | 概要 | 認証 |
| --- | --- | --- | --- |
| POST | `/auth/sessions` | Firebase IDトークンを受け取り独自JWTを返す | 不要 |
| GET | `/me` | 自分のユーザー情報を返す | 必要 |
| GET | `/posts` | グローバルタイムライン（カーソル） | 必要 |
| POST | `/posts` | 投稿作成 | 必要 |
| GET | `/posts/{id}` | 投稿詳細 | 必要 |
| DELETE | `/posts/{id}` | 投稿削除（本人のみ） | 必要 |
| POST | `/posts/{id}/likes` | いいね | 必要 |
| DELETE | `/posts/{id}/likes` | いいね解除 | 必要 |
| GET | `/users/{id}` | ユーザー情報 | 必要 |
| GET | `/users/{id}/posts` | そのユーザーの投稿一覧（カーソル） | 必要 |
| GET | `/healthz` | 死活監視 | 不要 |

### エラーレスポンス共通形式

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "post not found"
  }
}
```

主要エラーコード:
- `UNAUTHORIZED` (401)
- `FORBIDDEN` (403)
- `RESOURCE_NOT_FOUND` (404)
- `VALIDATION_ERROR` (400)
- `RATE_LIMITED` (429)
- `INTERNAL_ERROR` (500)

## 9. 画面構成

| 画面 | パス | 内容 |
| --- | --- | --- |
| サインイン | `/login` | Google ログインボタン |
| タイムライン | `/` | グローバルタイムライン、無限スクロール、投稿ボタン（FAB） |
| 投稿作成モーダル | `/` 内のオーバーレイ | テキスト + メディア選択 + 投稿ボタン |
| 投稿詳細 | `/posts/{id}` | 1投稿の表示 + 削除ボタン（本人時のみ） |
| プロフィール | `/users/{id}` | 表示名 + アバター + 投稿一覧 |
| 設定 | `/settings` | テーマ切替、ログアウト |

ヘッダ・サイドナビは MUI の `AppBar` + `Drawer`（モバイルではボトムナビに切替）。

## 10. アーキテクチャ

### 10.1 全体構成図（論理）

```
[Browser]
  | (HTTPS)
  v
[Firebase Hosting] -- 静的SPA配信
  | (XHR/fetch)
  v
[Cloud Run or Render] -- Go (echo) API
  |                       |
  |                       +--> [Cloud SQL or 管理Postgres]
  |
  +--> [Firebase Auth]   (IDトークン検証)
  |
  +--> [Firebase Storage] (メディアはブラウザから直接アップロード/ダウンロード、APIは経由しない)
```

### 10.2 バックエンドのレイヤ構成

```
backend/
  cmd/server/main.go             -- エントリポイント
  internal/
    domain/                       -- ドメインモデル（純粋なGo型）
    repository/                   -- gorm を使う実装。インターフェイスとSeparate
    usecase/                      -- ユースケース。トランザクション境界もここ
    handler/                      -- echo のハンドラ。リクエストパース・レスポンス整形
    middleware/                   -- 認証、レート制限、ロギング、リカバリ
    auth/                         -- Firebase Admin の薄いラッパ、JWT発行/検証
    storage/                      -- Firebase Storage の薄いラッパ（オブジェクト削除用）
    config/                       -- 環境変数読み込み
  docs/                           -- swag init で生成（git管理）
  migrations/                     -- golang-migrate の SQL ファイル
```

### 10.3 フロントエンドの構成

```
frontend/
  src/
    api/                          -- openapi-typescript で生成された型、fetch wrapper
    components/                   -- 汎用 UI コンポーネント
    features/                     -- 機能単位（posts, likes, auth, profile, theme）
    pages/                        -- ルーティング対応のページ
    hooks/                        -- 共通フック
    theme/                        -- MUI テーマ定義（light / dark）
    i18n/                         -- （v2 用空ディレクトリ。MVPでは未使用）
    main.tsx
  index.html
  vite.config.ts
  package.json
```

ルーティングは `react-router`、状態管理は `@tanstack/react-query` を中心に。グローバル状態（テーマ、認証ユーザー）は React Context。

## 11. 開発環境

### 11.1 ローカル

- `docker-compose up` で Postgres と Backend（air でホットリロード）を起動
- Frontend は `npm run dev`（Vite dev server, port 5173）
- Firebase Local Emulator Suite で Auth / Storage をエミュレート（`firebase emulators:start`）
- 各プロセスは別ターミナルで運用

### 11.2 docker-compose（要点）

- `db`: postgres:16-alpine、ヘルスチェック、ボリュームマウント
- `backend`: マルチステージで dev は air、本番イメージは別Dockerfile

### 11.3 環境変数

| 名前 | 用途 | 例 |
| --- | --- | --- |
| `DATABASE_URL` | Postgres 接続 | `postgres://app:app@db:5432/app?sslmode=disable` |
| `JWT_SECRET` | 独自JWT署名キー | （ランダム32バイト以上） |
| `FIREBASE_PROJECT_ID` | Firebase Admin 検証用 | `cloudcode-sns-dev` |
| `FIREBASE_CREDENTIALS_JSON` | サービスアカウントJSON（本番） | （秘匿） |
| `FIREBASE_AUTH_EMULATOR_HOST` | エミュレータ接続（ローカルのみ） | `localhost:9099` |
| `FIREBASE_STORAGE_EMULATOR_HOST` | エミュレータ接続（ローカルのみ） | `localhost:9199` |
| `CORS_ALLOWED_ORIGINS` | CORS 許可元 | `http://localhost:5173` |
| `RATE_LIMIT_POSTS_PER_MIN` | 投稿レート制限 | `5` |
| `RATE_LIMIT_LIKES_PER_MIN` | いいねレート制限 | `30` |
| `RATE_LIMIT_GET_PER_MIN` | GET レート制限 | `60` |

## 12. テスト戦略

- **バックエンド**:
  - usecase はモック repository で unit テスト
  - repository / handler は `testcontainers-go` で Postgres を立てて integration テスト
  - 認証ミドルウェアは Firebase Admin のモックで検証
- **フロントエンド**:
  - 主要コンポーネントを Vitest + Testing Library で unit テスト
  - 「ログイン → 投稿 → いいね → 削除」の golden path を Playwright で1〜2本
- **E2E のCI実行**: PRごとは省略可、main マージ前にナイトリーで実行（コスト都合で判断）

## 13. デプロイ

### 13.1 環境

MVPは本番（prod）環境のみ。

### 13.2 パイプライン

- GitHub Actions
- main への push トリガで:
  - backend: Docker イメージビルド → Cloud Run（または Render）にデプロイ → デプロイ後に golang-migrate でマイグレーション
  - frontend: `npm ci && npm run build` → Firebase Hosting にデプロイ
- PR では: lint / typecheck / test を実行（デプロイなし）

### 13.3 シークレット管理

- GitHub Actions Secrets（CI用）
- Cloud Run の Secret Manager 連携 / Render の Environment（実行時用）

## 14. オープン課題（仕様確定後に詰める）

- ログアウト時のJWT失効方針（短命JWT + ブラックリストなしで運用するか、Refreshトークン導入か）
- Firebase Storage のオブジェクト命名規則（`users/{userId}/posts/{postId}/{filename}` 案）
- バックエンドホスティング先の最終決定（Cloud Run か Render か）
- 動画長 30秒のサーバ側検証（ffprobe を入れるか、クライアント申告を信用するか）

## 15. 改訂履歴

| 日付 | 変更内容 |
| --- | --- |
| 2026-05-24 | 初版（MVPスコープ確定） |
