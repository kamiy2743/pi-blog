# AGENTS.md (blog)

## 前提
- ルートの `/home/kamiy2743/workspace/AGENTS.md` の方針に従う。
- 対象: `blog` プロジェクト。Go + SQLite でブログを構築する。

## 目的
- `/blog` にブログサイトを作成する。
- `blog.panda-dev.net/{articleId}` で記事を表示する。
- `/admin` で記事の公開/非公開を設定できる。

## 想定構成 (たたき台)
### ルーティング
- `GET /health` : 監視用 (200 OK)
- `GET /` : トップページ (最新記事一覧)
- `GET /{articleId}` : 記事表示
- `GET /admin` : 管理画面トップ (記事一覧・公開/非公開・新規作成)
- `GET /admin/new` : 新規作成フォーム
- `POST /admin/new` : 新規作成
- `GET /admin/edit/{articleId}` : 編集フォーム
- `POST /admin/edit/{articleId}` : 更新
- `POST /admin/publish/{articleId}` : 公開/非公開の切替 (Basic Auth)

### データベース (SQLite)
ファイル: `/home/kamiy2743/workspace/blog/data/blog.sqlite3`

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
- HTML テンプレートでトップページと記事ページを描画。
- Markdown を HTML に変換 (Go のライブラリは後で決定)。
- HTML 変換は保存時に行う。
- 公開されていない記事、存在しないパスはトップページ

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
- `cmd/blog/` : main
- `internal/` : handler, store, model
- `templates/` : HTML
- `static/` : CSS/JS
- `migrations/` : SQL

## 未決事項 / 質問
- `articleId` は数値IDで運用する。
