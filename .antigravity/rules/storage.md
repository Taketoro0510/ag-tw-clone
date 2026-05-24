# Storage ルール（Firebase Storage）

## 方針

- メディア（画像 / 動画）はブラウザから Firebase Storage に **直接アップロード / ダウンロード** する
- バックエンドは Storage 上のオブジェクトの **パス・URL を DB に記録** するだけ。ファイル本体はバック経由で送らない
- 削除時のみ、バックエンドが Firebase Admin SDK でオブジェクトを削除する

## オブジェクトパス命名規則

```
users/{userId}/posts/{postId}/{uuid}.{ext}
```

- `userId`, `postId` は UUID v7
- `uuid` はクライアントが生成（衝突回避用、1投稿に1ファイルだが命名衝突を避ける）
- 拡張子: `jpg / jpeg / png / webp / gif / mp4 / mov`

## アップロード手順（フロント）

1. ユーザーがファイルを選択
2. クライアント側でバリデーション:
   - 画像: MIME `image/*`、サイズ 5MB 以下
   - 動画: MIME `video/mp4` or `video/quicktime`、サイズ 50MB 以下、長さ 30秒以内（`<video>.duration`）
3. postId（UUID v7）をクライアントで生成
4. Firebase Storage SDK で上記パスにアップロード
5. アップロード完了 → `downloadURL` を取得
6. `POST /api/v1/posts` に `{ body, media_type, media_path, media_url, post_id, media_duration_seconds? }` を送信

## Firebase Storage のセキュリティルール

`storage.rules` で:

- `users/{userId}/...` への書き込みは、認証済かつ `request.auth.uid == userId` のみ
- 読み込みは全員許可（MVP は公開SNSなので）
- ファイルサイズ・MIME のチェックも storage.rules で二重に行う

```
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    match /users/{userId}/posts/{postId}/{fileName} {
      allow read: if true;
      allow write: if request.auth != null
                   && request.auth.uid == userId
                   && request.resource.size < 50 * 1024 * 1024
                   && (request.resource.contentType.matches('image/.*')
                       || request.resource.contentType.matches('video/.*'));
    }
    match /users/{userId}/{allPaths=**} {
      allow read: if true;
      allow write: if false;
    }
  }
}
```

## 削除フロー

- `DELETE /api/v1/posts/{id}` を受けたバックは:
  1. DB トランザクションで `posts.deleted_at` を更新
  2. コミット後、Firebase Admin の `bucket.file(media_path).delete()` を呼ぶ
  3. 削除失敗時は warn ログ + 後続のクリーンアップジョブで再試行（ジョブはMVPでは未実装、ログを見つけたら手動）

## エミュレータ

- ローカル開発では Firebase Storage Emulator を使う（ポート 9199 デフォルト）
- 環境変数 `FIREBASE_STORAGE_EMULATOR_HOST=localhost:9199` をフロント/バック両方に設定
- エミュレータは Firebase CLI で起動: `firebase emulators:start --only auth,storage`

## バケット運用

- 1プロジェクトに1バケット（デフォルトバケット）
- 本番 / ローカルでバケットは分けない（エミュレータでローカルは完結する）
- 命名は Firebase プロジェクト名のデフォルトに従う

## 禁止事項

- バックエンド経由の multipart アップロードを実装しない（Storage 直アップロードに統一）
- 認証なしでの書き込みを許可しない（storage.rules で厳格化）
- Firebase の本番認証情報を git にコミットしない
