# AGENTS.md (blog)

## 前提
- ルートの `/home/kamiy2743/workspace/AGENTS.md` の方針に従う。
- 対象: `blog` プロジェクト。`backend/`(Go + SQLite) と `frontend/`(Svelte 予定) に分割する。

## 目的
- `/blog` にブログサイトを作成する。
- `blog.panda-dev.net/article/{articleId}` で記事を表示する。
- `/admin` で記事の公開/非公開を設定できる。

## 現在の構成 (2026-03-01 時点)
### ルーティング
- `GET /health` : 監視用 (200 OK)
- `GET /` : ブログサイトのトップページ
- `GET /article` : 記事一覧
- `GET /article/` : 記事一覧
- `GET /article/{articleId}` : 記事表示
- `GET /admin` : 管理画面トップ (記事一覧・公開/非公開・新規作成)
- `GET /admin/` : 管理画面トップ
- `GET /admin/article/new` : 新規作成フォーム
- `POST /admin/article/new` : 新規作成
- `GET /admin/article/edit/{articleId}` : 編集フォーム
- `POST /admin/article/edit/{articleId}` : 更新
- `POST /admin/article/publish/{articleId}` : 公開/非公開の切替 (Basic Auth)

### データベース (SQLite)
ファイル: `/home/kamiy2743/workspace/blog/backend/data/blog.sqlite3`

テーブル案 `articles`:
- `id` TEXT PRIMARY KEY (URL の `{articleId}`)
- `title` TEXT NOT NULL
- `content_md` TEXT NOT NULL
- `content_html` TEXT NOT NULL (保存 or リクエスト時生成は要検討)
- `published` INTEGER NOT NULL DEFAULT 0
- `created_at` TEXT NOT NULL
- `updated_at` TEXT NOT NULL
- `published_at` TEXT NULL

### 表示
- Inertia + Svelte でトップページ/記事一覧/記事表示/管理画面を描画 (stub)。
- Markdown を HTML に変換 (Go のライブラリは後で決定)。
- HTML 変換は保存時に行う。
- 公開されていない記事、存在しないパスは 404

### 記事の作成・編集
- `/admin` で Markdown を貼り付けて作成・編集する。
- ローカルで AI に作成した Markdown を手動でコピペする運用。

### 管理
- `/admin` は Basic Auth で保護する。
- 管理画面で作成・編集・公開/非公開の切替を行う。

### セキュリティ方針 (強め)
- Markdown→HTML は保存時に必ずサニタイズする (許可タグ/属性を制限)。
- `/admin` は強い Basic Auth を使用し、ブルートフォース対策を入れる (レート制限)。
- CSRF 対策を導入する (CSRF トークン + SameSite Cookie)。
- `articleId` は数値のみ許可し、必ずバリデーションする。
- 非公開記事は 404 (存在を隠す)。
- DB はプレースホルダで操作し SQL インジェクションを防ぐ。
- SQLite ファイル権限は `600` を想定。

### ディレクトリ案
- `backend/` : Go API / サーバー
- `backend/cmd/blog/` : main
- `backend/internal/` : handler, model, store (予定)
- `backend/templates/` : HTML (現在は `backend/templates/root.html`)
- `backend/static/` : CSS/JS
- `frontend/server/` : Inertia SSR ビルド出力
- `backend/migrations/` : SQL
- `frontend/` : Svelte (Inertia)

### 実装メモ
- `internal/model` に `ArticleID` 型と `ParseArticleID` を実装済み。
- `internal/handler` は `ShowTop` / `ShowArticleList` / `ShowArticle` などに分割。
- `/home/kamiy2743/workspace/blog/backend/.env` に `HOST` / `PORT` / `FRONT_URL` / `INERTIA_ROOT_TEMPLATE` / `BUILD_DIR` を設定。
- systemd: `/home/kamiy2743/workspace/blog/backend/service/blog-back.service` を `blog-back.service` として配置。実行バイナリは `backend/blog-back`。
- Inertia: `backend/templates/root.html` と `/build/` 配信で利用中。
- Frontend build 出力: `frontend/` → `frontend/build/` にビルド出力。
- SSR: `frontend/src/blog-front.ts` を `frontend/server/` にビルドし、`frontend/service/blog-front.service` で起動。`frontend/.env` で `HOST` / `PORT` を設定。

## 未決事項 / 質問
- `articleId` は数値IDで運用する。
