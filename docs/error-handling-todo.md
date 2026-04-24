# Error Handling Migration TODO

## 方針

- 段階的に移行する。
- まず共通基盤を入れ、その後 endpoint ごとに route / handler / test を移行する。
- 1 endpoint ごとに「実装」「integration test 更新」「不要コード削除」まで完了させる。
- 進捗はこのファイルのチェックボックスで管理する。

## 進め方

1. 共通基盤を追加する。
2. 単純な GET endpoint から `InertiaPage` へ移行する。
3. partial reload を使う GET endpoint を移行する。
4. POST action を `InertiaAction` へ移行する。
5. 旧実装の helper / route / test を整理する。

## 共通基盤

- [x] `handlerresult.HandlerResult` の concrete type を追加する。
  - `PageResult`
  - `RedirectResult`
  - `RedirectBackResult`
- [x] `Page` / `Redirect` / `RedirectBack` helper を追加する。
- [x] `InertiaPage` adapter を追加する。
- [x] `InertiaAction` adapter を追加する。
- [x] adapter 内で `HandlerResult` を type switch して response へ変換する。
- [x] `DisplayableError` の GET page 変換を adapter に寄せる。
- [x] `RedirectBack` の `Referer` 解決を実装する。
  - `Referer` を優先
  - query string は保持
  - ない場合は `/`

## Session / Flash

- [x] `scs` の session manager を導入する。
- [x] session middleware を追加する。
- [x] validation error の一時保存を実装する。
- [x] global error / success flash の一時保存を実装する。
- [x] GET page adapter で session から `validationErrors` / `flash` を props に展開する。
- [x] 展開後に session から消費する。

## Frontend 共通部品

- [x] `flash` 表示用の共通 component を追加する。
  - 例: `frontend/src/components/FlashMessage.svelte`
- [x] `flash` 型を共通化する。
- [ ] page 上部または共通 layout に配置する。
- [ ] `validationErrors` と `flash` の props 命名を揃える。

## Test Helper

- [ ] Inertia test helper を `validationErrors` 前提へ更新する。
- [ ] flash props を検証できる helper を追加する。
- [ ] redirect follow を扱える helper を追加または更新する。
- [ ] session cookie を跨ぐ test を書けるようにする。

## Phase 1: 単純な GET を移行

### GET /

- [x] `ShowTopHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [x] route を `handler.InertiaPage(...)` へ差し替える。
- [x] `ShowTop` の integration test を新契約に合わせて更新する。
- [x] 旧 `inertia.RenderError` 依存を削除する。

### GET /article/{articleId}

- [x] `ShowArticleHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [x] 不正 ID / not found を `DisplayableError{StatusCode: 404}` に統一する。
- [x] route を `handler.InertiaPage(...)` へ差し替える。
- [x] `ShowArticle` の integration test を更新する。
- [x] 旧 `RenderNotFound` 依存を削除する。

### GET /admin/category

- [x] `EditCategoryHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [x] route を `handler.InertiaPage(...)` へ差し替える。
- [x] `EditCategory` の integration test を更新する。

## Phase 2: partial reload あり GET を移行

### GET /article

- [x] `SearchArticleHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [x] GET query validation error を `PageResult.ValidationErrors` に載せる。
- [x] partial reload failure 時に表示用 props と `flash.error` を返す。
- [x] frontend の `only` に `flash` を含める。
- [x] `FlashMessage` component を画面へ適用する。
- [x] integration test を追加 / 更新する。
  - [x] 通常 GET の displayable error
  - [x] partial reload の validation error
  - [x] partial reload の displayable error
  - [x] error 時の表示用 props

### GET /admin

- [x] `ShowAdminHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [x] GET query validation error を `PageResult.ValidationErrors` に載せる。
- [x] partial reload failure 時に表示用 props と `flash.error` を返す。
- [x] frontend の `only` に `flash` を含める。
- [x] `FlashMessage` component を画面へ適用する。
- [x] integration test を追加 / 更新する。
  - 通常 GET の displayable error
  - partial reload の validation error
  - partial reload の displayable error
  - error 時の表示用 props

## Phase 3: POST action を移行

### POST /admin/category

- [ ] `StoreCategoryHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [ ] validation error を `RedirectBackResult.ValidationErrors` に載せる。
- [ ] displayable / unexpected error を `flash.error` に載せる。
- [ ] success message を必要なら `flash.success` に載せる。
- [ ] route を `handler.InertiaAction(...)` へ差し替える。
- [ ] integration test を更新する。

### POST /admin/category/{categoryId}

- [ ] `UpdateCategoryHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [ ] 複数 form 用の scoped `validationErrors` を実装する。
  - 例: `update.{categoryId}.name`
- [ ] route を `handler.InertiaAction(...)` へ差し替える。
- [ ] integration test を更新する。

### POST /admin/category/{categoryId}/delete

- [ ] `DestroyCategoryHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [ ] route を `handler.InertiaAction(...)` へ差し替える。
- [ ] integration test を更新する。

### POST /admin/article/new

- [ ] `StoreArticleHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [ ] validationErrors / flash / redirect ルールを適用する。
- [ ] route を `handler.InertiaAction(...)` へ差し替える。
- [ ] integration test を更新する。

### POST /admin/article/{articleId}

- [ ] `UpdateArticleHandler` を `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
- [ ] validationErrors / flash / redirect ルールを適用する。
- [ ] route を `handler.InertiaAction(...)` へ差し替える。
- [ ] integration test を更新する。

## Phase 4: Route 整理

- [ ] `route.go` から `inertiaApp.Middleware(...)` の直接指定を減らす。
- [ ] endpoint ごとの adapter 利用に統一する。
- [ ] Basic Auth と adapter の適用順を確認する。
- [ ] health check など plain HTTP route が巻き込まれていないか確認する。

## Phase 5: 旧実装の掃除

- [ ] 旧 `inertia.RenderError` 依存を削除する。
- [ ] 旧 `inertia.RenderNotFound` 依存を削除または縮小する。
- [ ] `PrepareInput` など旧 helper の責務を見直す。
- [ ] 不要になった route / helper / test code を削除する。

## 完了条件

- [ ] 全 Inertia page route が `InertiaPage` 経由になっている。
- [ ] 全 POST action route が `InertiaAction` 経由になっている。
- [ ] validation error の最終返し先が常に `props.validationErrors` になっている。
- [ ] global error / success が共通 `FlashMessage` component で表示される。
- [ ] GET / POST / partial reload / not found の integration test が通る。
