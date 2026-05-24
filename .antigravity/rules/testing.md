# Testing ルール

## 全体方針

- MVPは「バックエンドに重点」を置きつつ、フロントの主要パスは Playwright で押さえる
- テストの目的: 「リファクタしても壊れないという信頼」を作る
- 100% カバレッジは目指さない。価値の高い箇所（ユースケース / API / golden path）を厚く

## バックエンド（Go）

### Unit テスト

- 対象: `usecase/`、`auth/`、`middleware/` の純粋ロジック
- リポジトリはモック化（インターフェイス → モック実装。`mockery` または手書き）
- 標準の `testing` パッケージ + `github.com/stretchr/testify` のアサーション

### Integration テスト

- 対象: `repository/`、`handler/` の主要 API
- **`testcontainers-go`** で Postgres コンテナを毎テスト or 毎パッケージで立ち上げる
- マイグレーションは `golang-migrate` をテスト setup で適用
- Firebase Admin は **モック**（IDトークン検証部分をインターフェイス化）

### ファイル配置

- `*_test.go` を対象パッケージと同じ場所に
- integration テストは `//go:build integration` ビルドタグで分離（`go test -tags=integration ./...` で実行）

### 命名

- `TestPostUsecase_Create_ReturnsValidationError_WhenBodyExceeds140Chars` のように
  `Test{Subject}_{Method}_{ExpectedOutcome}_{Condition}` パターン

### CI

- PR で `go test ./...`（unit）+ `go test -tags=integration ./...` の両方を実行

## フロントエンド

### Unit / コンポーネントテスト

- **Vitest** + **@testing-library/react**
- 対象:
  - feature 配下の主要コンポーネント（PostCard、LikeButton、PostComposer）
  - カスタムフック（`useAuth`、`useTimeline`）
- API モック: **MSW**（msw）でネットワーク層をモック
- 1コンポーネント1ファイル: `PostCard.test.tsx`

### E2E（Playwright）

- 対象: 「ログイン → 投稿 → いいね → 自分のプロフィールで確認 → 投稿削除」の golden path
- 環境: ローカル / CI とも Firebase エミュレータ + Postgres コンテナを使用
- Firebase エミュレータの Auth でテストユーザーを作成
- テストは `frontend/e2e/` に配置

### CI

- PR で Vitest を実行
- Playwright は **PR では選択的、main マージ前の nightly で全件**（コスト考慮、必要に応じて変更）

## モック方針

- バック: ドメインモデル外の I/O（DB / Firebase）はインターフェイス越しにし、テストでモック
- フロント: ネットワーク層を MSW でモック。Firebase SDK 自体のモックは最小限（feature/auth にラッパを置き、それをモック）

## テストデータ

- ユーザー / 投稿 / いいねのファクトリ関数を `testutil/` に用意（バック）、`tests/factories/` に用意（フロント）
- ファクトリは「最小フィールドだけ受け取り、残りはデフォルト埋め」のパターン

## 禁止事項

- テストでデータベースの **本番接続情報を使わない**
- スリープでタイミング合わせをしない（`waitFor` / `expect.poll` を使う）
- 1テスト関数で複数のシナリオを混ぜない（1 Arrange → 1 Act → 1 Assert）
