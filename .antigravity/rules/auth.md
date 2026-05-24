# Auth ルール（Firebase Auth + 独自JWT）

## 全体フロー

```
[Browser]
  1. Firebase Auth SDK で Google ログイン
  2. Firebase IDトークン（短命）を取得
  3. POST /api/v1/auth/sessions { idToken } をバックに送信
[Backend]
  4. firebase-admin で IDトークンを検証 → uid / email / name / picture を取得
  5. users テーブルに upsert（無ければ作成）
  6. 独自JWT（HS256, exp=24h）を発行して返す
[Browser]
  7. 受け取った独自JWTを localStorage に保存
  8. 以降のAPI呼び出しで Authorization: Bearer <jwt> を付ける
```

## 独自JWT の仕様

- アルゴリズム: HS256
- シークレット: `JWT_SECRET` 環境変数（32バイト以上のランダム）
- ペイロード:
  ```json
  {
    "sub": "<users.id (UUID v7)>",
    "iat": 1716508800,
    "exp": 1716595200
  }
  ```
- 期限切れ: 401 を返す。クライアントは Firebase IDトークンを再取得して `/auth/sessions` を再叩きで自動回復

Refresh トークンは MVP では発行しない。シンプルさを優先。

## バックエンドの認証ミドルウェア

- `Authorization: Bearer <jwt>` をパース → 署名検証 → `exp` チェック
- 成功時、`echo.Context` に `userID`（UUID）を載せる
- 失敗時、`{"error":{"code":"UNAUTHORIZED","message":"..."}}` で 401

`firebase-admin` の IDトークン検証は **`/auth/sessions` エンドポイントのみで使う**。それ以外の API は独自JWT のみで認証する。

## クライアントの実装方針

- Firebase Auth SDK は **クライアントの auth feature のみ** に閉じ込める
- 他の feature は `useAuth()` フックから「ログイン中のユーザー」「JWT」だけを取得する
- 401 を受けたら一度だけ自動リトライ:
  1. Firebase の `getIdToken(true)` で再取得
  2. `/auth/sessions` で独自JWT 再発行
  3. 元のリクエストを 1 回だけ再送
  4. それでも失敗ならログアウトしてサインイン画面へ

## エミュレータ

- ローカル開発では Firebase Auth Emulator を使う
- 環境変数 `FIREBASE_AUTH_EMULATOR_HOST=localhost:9099` を設定
- firebase-admin はこの env を見て自動でエミュレータに接続する

## セキュリティ上の注意

- **本番の Firebase サービスアカウントJSONを git にコミットしない**（GitHub Actions Secrets / Cloud Run Secret Manager 経由）
- JWT シークレットも同上
- localStorage に JWT を置く方針は XSS リスクを負う。MUI / React のデフォルトエスケープを破る（`dangerouslySetInnerHTML` 等）を絶対に使わない
- CORS は許可ドメインのみ（本番ドメイン + ローカル `http://localhost:5173`）

## ログアウト

- クライアント: localStorage から JWT 削除 + Firebase Auth SDK の `signOut()`
- サーバ: JWT 失効リストは持たない（短命 + 24h 内に再ログインで失効するため許容）
- 厳密な失効が必要になったら Redis ベースのブラックリストを別途設計

## 禁止事項

- API リクエストごとに Firebase IDトークンを送らない（独自JWT を使う）
- httpOnly cookie 認証への変更は MVP では行わない（将来検討）
- パスワード認証を追加しない（MVPは Google のみ）
