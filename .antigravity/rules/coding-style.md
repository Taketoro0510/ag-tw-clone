# コーディングスタイル ルール（共通）

## 一般原則

- **読みやすさ > 賢さ**。3行のシンプルなコードは、抽象化された1行に勝る
- **早期 return** で nest を浅く保つ
- **意味のある名前**。`data`, `info`, `obj`, `tmp` は禁止。`post`, `userId` のような名詞を使う
- **コメントは「なぜ」を書く**。「何を」はコードが語る。書かないのが基本
- **不要な抽象化を作らない**。同じパターンが3箇所目に現れたら抽出を検討
- **未来の機能のための余白を残さない**。YAGNI

## Go（バックエンド）

- フォーマット: `gofmt` / `goimports`（CI で強制）
- リンタ: `golangci-lint`（errcheck, staticcheck, gosec, govet, ineffassign, unused）
- パッケージ名: 短く、複数形にしない（`user` ではなく `users` は避ける）
- インターフェイスは **使う側のパッケージ** に定義（小さく保つ）
- エラー: `fmt.Errorf("doing X: %w", err)` で wrap
- nil チェック: `if err != nil { return err }` の典型形を保つ
- struct タグ: `db:"col_name" json:"camelCase"`、固定の順序で
- goroutine: `context.Context` を渡してキャンセル可能に

## TypeScript / React（フロントエンド）

- フォーマット: Prettier（CI で強制）
- リンタ: ESLint（`@typescript-eslint/recommended`, `react-hooks/recommended`）
- import 順: 外部ライブラリ → 内部モジュール → 同階層
- 関数コンポーネント + フック。クラスコンポーネントは禁止
- props は型を明示（`type Props = { ... }`）
- 1コンポーネント1ファイル。ファイル名は PascalCase（`PostCard.tsx`）
- ロジックを含むカスタムフックは `useXxx` 命名で `hooks/` または feature 内
- 副作用は `useEffect` で。依存配列を正確に書く（ESLint の `exhaustive-deps` を守る）
- 早期 return で「ローディング中 / エラー」を先に処理してから本体を描画

## ファイル / ディレクトリ命名

- Go: `snake_case.go`（複数語のみ、単語1つなら lowercase）
- TS: `PascalCase.tsx`（コンポーネント）、`camelCase.ts`（その他）
- ディレクトリ: `kebab-case` か `lowercase`（混在しない、Repo 内で統一）

## コミット

- 1コミット1論理変更
- メッセージ:
  - 1行目: `<area>: <短い変更内容>`（`backend:`, `frontend:`, `infra:`, `docs:`）
  - 必要なら本文で「なぜ」
- 例: `backend: add rate limit middleware`

## PR

- タイトル: コミット先頭行と同じ規約
- 本文:
  - **Summary**: 何を変えたか（箇条書き2〜4個）
  - **Why**: なぜそれが必要か
  - **Test plan**: 動作確認の手順 / テスト追加した範囲
- 1 PR 1 トピック。混ぜない

## レビュー観点

- 仕様書（SPEC.md）と矛盾していないか
- 認証 / 認可が必要な API に付いているか
- ソフトデリート対象クエリで `deleted_at IS NULL` を忘れていないか
- N+1 になっていないか
- エラーがユーザーに見える文言として適切か（個人情報や内部詳細を漏らしていないか）
- フロントが API レスポンス型を手書きしていないか（openapi-typescript 生成物を使うこと）

## 禁止事項

- 仕様書を更新せずに大きな仕様変更を実装する
- 「将来のため」のフラグ / hook / 抽象を入れる
- コードを消さずに `// 旧実装` のコメントとして残す
- secrets を含むファイルを git に追加する（`.gitignore` を守る、`git secrets` 推奨）
- `any` を使う（TS） / `interface{}` をパブリック API で使う（Go）
