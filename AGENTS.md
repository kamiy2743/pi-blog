# AGENTS.md (blog)

## 前提
- ルートの `/home/kamiy2743/workspace/AGENTS.md` の方針に従う。
- 対象は `blog` プロジェクト。
- 構成は `backend/` (Go + SQLite), `frontend/` (Svelte + Inertia), `nginx/`, `cloudflared/`。

## 目的
- `blog.panda-dev.net` でブログを公開する。
- `GET /article/{articleId}` で記事を表示する。
- `/admin` で記事の作成・編集・公開/非公開切替を行う。

## 現在の構成 (2026-03-08 時点)
### コンテナ構成
- `docker-compose.yml` で `nginx`, `go`, `node`, `cloudflared` を起動する。
- `cloudflared` は `cloudflared/config.yml` の ingress で `blog.panda-dev.net -> http://nginx:8000` に転送する。
- `nginx` は `location /` を `http://go:8001` に proxy する。
- `go` は Inertia SSR を使い、`FRONT_URL=http://node:8002` に SSR リクエストを送る。

### ルーティング
- `GET /health`
- `GET /`
- `GET /article`
- `GET /article/`
- `GET /article/{articleId}`
- `GET /admin`
- `GET /admin/`
- `GET /admin/article/new`
- `POST /admin/article/new`
- `GET /admin/article/edit/{articleId}`
- `POST /admin/article/edit/{articleId}`
- `POST /admin/article/publish/{articleId}`

### データベース (SQLite)
- DB ファイル: `/home/kamiy2743/workspace/blog/backend/data/blog.sqlite3`
- テーブル: `articles` (`id`, `title`, `content_md`, `content_html`, `published`, `created_at`, `updated_at`, `published_at`)

## セキュリティ方針
- `/admin` は Basic Auth で保護する。
- Markdown -> HTML は保存時にサニタイズする。
- CSRF 対策 (トークン + SameSite Cookie) を導入する。
- 非公開記事は `404` を返して存在を隠す。
- SQL はプレースホルダを使う。
