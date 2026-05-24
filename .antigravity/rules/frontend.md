# Frontend ルール（React / TypeScript / MUI）

## ディレクトリ構成

```
frontend/src/
  api/            -- openapi-typescript の生成型 + fetch wrapper（共通エラー処理、認証ヘッダ付与）
  components/     -- 汎用 UI（Button のラッパ等、機能横断）
  features/       -- 機能単位（posts/, likes/, auth/, profile/, theme/）
                  --   各 feature 内に components / hooks / api を置く
  pages/          -- react-router のルーティング対象
  hooks/          -- アプリ横断のカスタムフック
  theme/          -- MUI テーマ定義（lightTheme.ts, darkTheme.ts, ThemeProvider）
  i18n/           -- v2 用の空ディレクトリ（MVPでは使わない）
  main.tsx
  App.tsx
```

`features/` は機能凝集。1機能で完結する小宇宙にして、ページからは feature の公開コンポーネント / フックを呼ぶだけにする。

## TypeScript

- `tsconfig.json` で `strict: true`, `noUncheckedIndexedAccess: true`
- API のレスポンス型は `openapi-typescript` 生成の `types.ts` から import する（手書きしない）
- `any` 禁止。やむを得ない場合は `unknown` → 型ガードで絞り込む

## 状態管理

- サーバ状態は **`@tanstack/react-query` で統一**（投稿一覧、ユーザー情報、いいね状態など）
  - キャッシュキーは `["posts", { cursor }]` のように構造化
  - mutation 成功時の楽観更新でいいねボタンの即時反応を実現
- グローバルクライアント状態（認証ユーザー、テーマ）は React Context
- ローカル UI 状態は `useState` で十分

Redux / Zustand 等は MVP では入れない。

## ルーティング

- `react-router` v6+
- ルート定義は `App.tsx` または `src/routes.tsx` に集約
- 認証が必要なルートは `<RequireAuth>` ラッパで囲む

## MUI のテーマ

- `src/theme/` に `lightTheme.ts` / `darkTheme.ts` を置く
- パレット、タイポグラフィ、コンポーネントオーバーライドはここに集約
- ブランドカラーは MUI のデフォルトから少しずらしたカスタムカラーを定義（SPEC 5.6 参照）
- 切替は `<ThemeProvider>` で実現、選択は `localStorage` に保存
- システム設定追従モードは `prefers-color-scheme` メディアクエリで実装

## レスポンシブ

- ブレークポイントは MUI デフォルト（xs/sm/md/lg/xl）を使用
- 1024px 以上で PC レイアウト（サイドナビ）、767px 以下でモバイル（ボトムナビ）、間はタブレット
- 画像 / 動画は親要素にフィット（`object-fit: cover`、CSS は MUI の `sx` で）

## fetch / API クライアント

- `src/api/client.ts` に共通 fetch wrapper を作る
- JWT を `Authorization` ヘッダに自動付与
- 401 を受けたら一度だけ Firebase IDトークン再取得 → JWT 再発行を試み、失敗したらログアウト
- レスポンスは openapi-typescript の `paths` 型からエンドポイント別の型を導出（`client.get('/api/v1/posts')` のような薄い API でも、テンプレートリテラル型で型安全に）

## アップロード

- 画像/動画は Firebase Storage SDK で **直接** アップロード
- パス命名: `users/{userId}/posts/{postId}/{filename}`（ただし postId はアップロード時点では未確定なので、クライアントで UUID v7 を生成して使う）
- アップロード成功後にバックの `POST /posts` に `media_path` を渡す

## アクセシビリティ

- ボタンには `aria-label`、画像には `alt` を必ず付与
- フォーカスリングを消さない（MUI デフォルトでOK）
- カラーコントラスト 4.5:1 以上を確保（Light/Dark とも）

## エラー表示

- API エラーは `react-query` の `onError` でトーストに変換
- 422 / 400 のバリデーションエラーはフォームの該当フィールドに表示

## 禁止事項

- `npm install` で安易に依存追加しない（PR で必ず妥当性を説明）
- `any` を使わない
- インライン CSS は最小限。スタイルは MUI の `sx` か theme で
- API レスポンス型を手書きしない（openapi-typescript で生成）
