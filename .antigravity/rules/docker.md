# Docker / ローカル開発 ルール

## 目的

- ローカルで `docker-compose up` するだけで **Postgres + バックエンド（ホットリロード付）** が立ち上がる状態を維持する
- フロントエンドと Firebase エミュレータはホスト側で起動（コンテナ化しない）

## 構成

```
project_root/
  docker-compose.yml
  backend/
    Dockerfile           -- 本番用、マルチステージ
    Dockerfile.dev       -- 開発用、air が動く
    .air.toml            -- air の設定
  firebase.json          -- Firebase エミュレータ設定
```

## docker-compose.yml（要点）

```yaml
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: app
    ports:
      - "5432:5432"
    volumes:
      - db-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app"]
      interval: 5s
      timeout: 5s
      retries: 10

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgres://app:app@db:5432/app?sslmode=disable
      JWT_SECRET: dev-secret-do-not-use-in-prod-xxxxxxxxxxxxxxxx
      FIREBASE_PROJECT_ID: cloudcode-sns-local
      FIREBASE_AUTH_EMULATOR_HOST: host.docker.internal:9099
      FIREBASE_STORAGE_EMULATOR_HOST: host.docker.internal:9199
      CORS_ALLOWED_ORIGINS: http://localhost:5173
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
      - go-mod-cache:/go/pkg/mod
    depends_on:
      db:
        condition: service_healthy

volumes:
  db-data:
  go-mod-cache:
```

`host.docker.internal` で Firebase エミュレータ（ホスト起動）に接続する。Linux では `extra_hosts: ["host.docker.internal:host-gateway"]` を追加。

## backend/Dockerfile.dev（air 使用）

```dockerfile
FROM golang:1.23-alpine

RUN apk add --no-cache git curl
RUN go install github.com/air-verse/air@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
```

## backend/.air.toml（要点）

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/server ./cmd/server"
  bin = "./tmp/server"
  include_ext = ["go"]
  exclude_dir = ["tmp", "docs", "migrations"]
  delay = 200
```

`go run` を air の cmd に使わない（毎回フルコンパイルになる）。`go build` でバイナリを作り、それを air が走らせる。

## 本番イメージ（backend/Dockerfile）

- マルチステージ: builder で `go build`、distroless または alpine で起動
- 生成バイナリのみコピー、ソースは含めない
- non-root ユーザー
- `CMD ["/app/server"]`

## Firebase エミュレータ

- ホスト側で `firebase emulators:start --only auth,storage`
- `firebase.json` でポート / プロジェクトIDを固定:
  ```json
  {
    "emulators": {
      "auth": { "port": 9099 },
      "storage": { "port": 9199 },
      "ui": { "enabled": true, "port": 4000 }
    },
    "storage": { "rules": "storage.rules" }
  }
  ```

## フロントエンドの起動

- `cd frontend && npm run dev`（Vite, port 5173）
- 環境変数は `.env.local`:
  ```
  VITE_API_BASE_URL=http://localhost:8080/api/v1
  VITE_FIREBASE_API_KEY=...
  VITE_USE_FIREBASE_EMULATOR=true
  ```

## マイグレーションの実行

- バックコンテナ起動時に自動で `migrate up` を実行（main.go で）
- 手動実行: `docker compose run --rm backend migrate -path migrations -database "$DATABASE_URL" up`

## ホットリロード時の注意

- 新しいファイルを作ったらコンテナ内に確実に伝わるよう、`volumes` で `./backend:/app` をバインドマウント
- Windows + WSL では IO 遅延に注意（プロジェクトを WSL ファイルシステム上に置くと早い）

## 禁止事項

- `go run` ベースのホットリロード（毎回フルコンパイルになり遅い）
- フロントエンドを docker-compose に入れない（Vite を直接動かす方が速い）
- 本番認証情報を docker-compose.yml にハードコードしない
