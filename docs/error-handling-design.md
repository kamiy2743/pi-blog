# Error Handling Design

## 目的

Go handler / usecase / Inertia response のエラーハンドリングを統一する。

実装は endpoint ごとに段階移行する。進捗管理は [error-handling-todo.md](/workspace/blog/docs/error-handling-todo.md) で行う。

現在は各 feature handler が `inertia.RenderError`、`inertia.RenderNotFound`、`http.Redirect` を直接呼ぶ箇所があり、GET / POST / Inertia request / URL 直アクセスごとの判断が分散しやすい。今後、記事作成・更新・カテゴリ管理の validation や保存失敗を実装すると、個別 handler ごとに response 方針がばらつく可能性がある。

この設計では、feature handler と usecase は成功結果または error を返し、HTTP response への変換は route adapter 側に集約する。

## 基本方針

- usecase は HTTP response を知らない。
- feature handler は request parse、validation、usecase 呼び出し、成功時 result の組み立てを担当する。
- feature handler は `http.ResponseWriter` に直接書き込まない。
- route adapter は result / error を HTTP response に変換する。
- GET のページ表示エラーは `ErrorPage` を表示する。404 も `ErrorPage` の一種として扱う。
- POST / PUT / PATCH / DELETE の入力エラーや処理失敗は、session / flash に一時保存して `303 See Other` で戻す。
- Basic Auth の `401` は `WWW-Authenticate` header が必要なので、middleware の plain HTTP response のまま扱う。

## レスポンス方針

| リクエスト | エラー | レスポンス |
| --- | --- | --- |
| GET 通常ページ表示 | query validation error | `validationErrors` props として同一ページに返す |
| GET 通常ページ表示 | path / resource validation error | `ErrorPage` with `404` または指定 status |
| GET 通常ページ表示 | displayable error | `ErrorPage` with 指定 status |
| GET partial reload | validation error | `validationErrors` props として同一ページに返す |
| GET partial reload | displayable error | 共通エラーメッセージ領域に表示 |
| POST / PUT / PATCH / DELETE | validation error | session に validation errors を一時保存して `303` redirect back。次の GET で `validationErrors` props に展開 |
| POST / PUT / PATCH / DELETE | displayable error | flash global error + `303` redirect back 後、共通エラーメッセージ領域に表示 |
| POST / PUT / PATCH / DELETE | success | `303` redirect to 成功後の画面 |
| Basic Auth 失敗 | unauthorized | plain `401` |

404 は原則 redirect back しない。リンク切れや存在しない記事詳細は URL と状態が一致する必要があるため、`ErrorPage` with `404` として表示する。

検索画面の query validation error は、リロードやブックマークからの再訪問でも同じ検索画面に表示する。`/article?title=<長すぎる値>` のような URL をリロードしただけで `ErrorPage` に遷移すると、URL と検索 UI の永続性が悪くなるため。

## package 責務

### usecase

usecase は domain / repository を呼び出し、成功値または error を返す。

usecase は以下を知らない。

- Inertia component 名
- HTTP status code
- redirect 先
- flash message
- `http.ResponseWriter`

ただし、現在の実装では `handlererror.DisplayableError` を usecase から返している。移行初期はこれを許容し、段階的に domain/application error へ寄せる。

### feature handler

feature handler は request を受け取り、成功時は `handlerresult.HandlerResult`、失敗時は `*handlererror.DisplayableError` を返す。

```go
func (h *Handler) Handle(r *http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)
```

feature handler は以下を担当する。

- path value / query / form の parse
- validation
- usecase 呼び出し
- Inertia component 名と props の組み立て
- 成功時 redirect result の組み立て

feature handler は以下をしない。

- `inertia.Render`
- `inertia.RenderError`
- `inertia.RenderNotFound`
- `http.Redirect`
- `http.Error`

### route adapter

route adapter は feature handler を `http.Handler` に変換する。

```go
handleAdmin("GET /admin", handler.InertiaPage(inertiaApp, container.ShowAdminHandler.Handle))
handleAdmin("POST /admin/category", handler.InertiaAction(inertiaApp, container.StoreCategoryHandler.Handle))
```

adapter は以下を担当する。

- `inertiaApp.Middleware` の適用
- result の描画
- error の種類ごとの分岐
- validation errors / flash message の保存
- redirect back
- `ErrorPage` の描画
- 想定外 error の logging

## HandlerResult

response は `handlerresult.HandlerResult` で表す。

```go
type HandlerResult interface {
	isResult()
}
```

adapter は `handlerresult.HandlerResult` の concrete type を type switch して処理を分ける。

```go
type PageResult struct {
	Component        string
	Props            gonertia.Props
	StatusCode       int
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}

type RedirectResult struct {
	To    string
	Flash *session.Flash
}

type RedirectBackResult struct {
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}
```

```go
type PageOptions struct {
	StatusCode       int
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}

type RedirectOptions struct {
	Flash *session.Flash
}

type RedirectBackOptions struct {
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}
```

```go
func Page(component string, props gonertia.Props, options ...PageOptions) HandlerResult
func Redirect(to string, options ...RedirectOptions) HandlerResult
func RedirectBack(options ...RedirectBackOptions) HandlerResult
```

options は可変長引数として省略可能にする。status code や validation errors / flash を付けたい場合だけ `PageOptions` / `RedirectOptions` / `RedirectBackOptions` を渡す。

`Page` は Inertia page を返す。status code が必要な場合は `PageOptions.StatusCode` を使う。`Redirect` / `RedirectBack` はそれぞれ `RedirectResult` / `RedirectBackResult` を返し、HTTP response としては `303 See Other` を返す。

`RedirectBack` は `Referer` header を優先して戻り先に使う。`Referer` に query string が含まれる場合はそのまま保持する。`Referer` がない場合の fallback はトップページ `/` とする。

validation errors / flash message も `handlerresult.HandlerResult` に載せる。feature handler は validation error や action failure 用の response 情報を別チャネルで返さず、常に `handlerresult.HandlerResult` に寄せる。

POST の成功時は feature handler が `http.Redirect` を呼ばず、`handlerresult.Redirect("/admin")` のように返す。

## Error Types

### ValidationError

validation error は既存の `handlererror.ValidationError` を維持する。

```go
type ValidationError struct {
	Field   string
	Message string
}
```

validation error の扱いは feature handler と adapter の組み合わせで決める。

- action request なら session に validation errors を一時保存して redirect back。次の GET で `validationErrors` props に展開する
- GET query validation error なら、通常表示でも partial reload でも feature handler が `validationErrors` 付きの `Page` result を返す
- GET path / resource validation error なら `DisplayableError` として `ErrorPage` を表示する

GET query validation error は error として early return しない。adapter だけでは component 名や props を再構築できないため、feature handler が入力値、画面を壊さない props、validation errors を含む `Page` result を返す。

`validationErrors` の key は form scope を含める前提にする。単一 form の画面では `title` や `name` のような field 名だけを使い、複数 form が同一画面にある場合は scope を必ず付ける。

例:

```txt
name
create.name
update.42.name
```

複数 form 画面では feature handler が scope 付き key を返し、画面側も対応する scope の error だけを表示する。

### DisplayableError

ユーザーに表示してよい message / description / status code を持つ error。

```go
type DisplayableError struct {
	StatusCode  int
	Message     string
	Description string
	Err         error
}
```

GET page adapter は `DisplayableError` を `ErrorPage` component に変換する。

404 も専用 component に分けず、`DisplayableError{StatusCode: 404}` として扱う。

```go
return &DisplayableError{
	StatusCode:  404,
	Message:     "ページが見つかりません。",
	Description: "URL が変わったか、公開が終了した可能性があります。",
	Err:         err,
}
```

通常 GET では `DisplayableError` を `ErrorPage` に変換する。

GET partial reload ではページ全体を `ErrorPage` に差し替えず、共通エラーメッセージ領域に表示する。検索条件、スクロール位置、現在表示中の一覧を維持し、ユーザーが同じ画面上で再試行できるようにするため。

action request では `DisplayableError` も unexpected error も flash global error + redirect back に変換する。

redirect back 後の GET page では flash global error を `flash.error` props に載せ、共通エラーメッセージ領域に表示する。

### Unexpected Error

型付けされていない error はログに出し、ユーザーには汎用 `500` を返す。

通常 GET では `ErrorPage` with `500`。

GET partial reload では、可能な限り共通エラーメッセージ領域に汎用エラーを表示する。ただし render / SSR など response 生成そのものの失敗は、共通エラー props を返せないため静的エラーへ逃がす。

action request では flash global error + redirect back。ただし render / SSR など response 生成そのものの失敗は redirect back せず、静的エラーに逃がす。

## Adapter

GET page 用と action 用は adapter を分ける。

```go
func InertiaPage(
	inertiaApp *gonertia.Inertia,
	handle func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError),
) http.Handler

func InertiaAction(
	inertiaApp *gonertia.Inertia,
	handle func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError),
) http.Handler
```

### InertiaPage

GET のページ表示向け。

- 成功時は result を描画する。
- query validation error は feature handler が `validationErrors` props を含む `Page` result として返す。adapter は validation error だけから page を再構築しない。
- path / resource validation error は `DisplayableError` にする。
- `DisplayableError` は `ErrorPage` にする。
- unexpected error は log + `ErrorPage` with `500` にする。
- adapter は原則 `err != nil` なら error response、`err == nil` なら result を描画する。ただし `PageResult` と `DisplayableError` が同時に返り、かつ partial reload の場合だけ、同じ page response に `flash.error` を載せる。

通常 GET ページ遷移では redirect back しない。URL の意味を維持するため。

GET query validation error は、通常表示と partial reload のどちらでも redirect back しない。検索 URL、入力欄、現在の一覧表示を維持することを優先し、feature handler が `validationErrors` props と画面を壊さない props を持つ同じ page response を返す。

GET partial reload の `DisplayableError` は `flash.error` props として同じ page response に返す。

### InertiaAction

POST / PUT / PATCH / DELETE 向け。

- 成功時は result を返す。原則 `Redirect`。
- validation error は session に一時保存して redirect back し、次の GET で `validationErrors` props に展開する。
- `DisplayableError` は flash global error + redirect back にする。
- unexpected error は log + flash global error + redirect back にする。

`inertiaApp.Middleware` の適用は route 定義側ではなく adapter 側に寄せる。移行後は `InertiaPage` / `InertiaAction` が必要な Inertia middleware を内部で適用する。Inertia form submit と通常 form submit の両方を扱えるよう、flash の読み出しは次回 GET の page adapter 側で行う。

## Flash

本設計では session flash を使う。

global message flash に載せる値は以下。

```go
// handler/session
type Flash struct {
	Error   string
	Success string
}
```

Inertia page 描画時に flash を読み出し、props に共有する。

validation error の最終的な返し先は、GET / POST を問わず常に `props.validationErrors` に統一する。

- GET query validation error は feature handler がその response の `validationErrors` props を返す
- POST validation error は session に一時保存し、redirect back 後の GET page adapter が `validationErrors` props に展開する

session は POST-redirect-GET を跨ぐための一時保管場所として使い、画面が直接参照する値の置き場所にはしない。

```ts
type SharedProps = {
  validationErrors: Record<string, string>
  flash: {
    error?: string
    success?: string
  }
}
```

GET partial reload の displayable error は session flash ではなく、その response の props として `flash.error` 相当の値を返す。redirect を挟まないため、検索 URL と画面状態を維持できる。

画面側には共通の flash message 表示 component を用意し、各 page の上部または共通 layout に配置する。各 page が `flash.error` / `flash.success` を個別に描画するのではなく、この component に集約する。

```ts
flash: {
  error?: string
  success?: string
}
```

共通 component は `flash.error` があればエラー alert、`flash.success` があれば成功 alert を表示する。validation は従来通り `validationErrors` props を各 input の近くに表示する。

### session 実装

backend には `scs` の session manager を導入し、`session.SessionManager.Middleware()` を共通 middleware chain に組み込む。

利点:

- POST redirect back と相性がよい。
- validation errors / global error / success message を同じ仕組みにできる。
- cookie にエラー本文を直接載せずに済む。

注意点:

- test helper に session cookie の扱いが必要。

現在の実装では `session.NewSessionManager(appEnv)` で以下を設定する。

- cookie name: `blog_session`
- lifetime: `5 * time.Minute`
- idle timeout: `3 * time.Minute`
- `HttpOnly: true`
- `Secure: true` は本番のみ。開発環境では localhost / SSH tunnel の都合で `false` にする。
- `SameSite: http.SameSiteLaxMode`

session cookie 設定は現時点では `handler/session` package 内の定数と `appEnv` 判定で持つ。今後 config に寄せる場合は、cookie name、lifetime、idle timeout、Secure の環境別切り替えを `config` から読めるようにする。

session には validation errors / flash message のような短い一時値だけを入れる。大きな payload や機密値は載せない。

## Partial Reload

`/article` と `/admin` の検索は Inertia の partial reload を使う。

```ts
router.visit(buildUrl(1), {
  method: 'get',
  preserveState: true,
  preserveScroll: true,
  replace: true,
  only: ['partialSearch', 'validationErrors', 'flash'],
})
```

GET の検索 query validation error や partial reload の displayable error を redirect back にすると、URL、入力欄、検索結果、errors の同期が崩れやすい。

そのため、GET query validation error は通常表示と partial reload のどちらでも redirect back ではなく、feature handler が同じ page response の `validationErrors` props として返す。

たとえば `/article?title=<長すぎる値>` を開いた場合、またはその URL をリロードした場合は、`ErrorPage` ではなく記事一覧画面を表示し、タイトル入力欄の近くに validation error を出す。

GET partial reload の displayable error も `ErrorPage` にはせず、同じ page response の `flash.error` として返す。たとえば検索中に repository error が起きた場合、検索画面自体は維持し、画面上部の共通エラーメッセージ領域に「記事の読み込みに失敗しました。」を表示する。

partial reload で共通メッセージ領域を更新したい route は、client 側の `only` に `flash` を明示的に含める。shared props であっても `only` に依存した配信漏れを避けるため、省略しない。

この場合、失敗した partial prop は現在の表示を壊さない値を返す。adapter は props の中身を推測しないため、partial reload で表示したい props は handler 側で `PageResult.Props` として組み立てる。

- 一覧検索なら直前の一覧を維持できるなら維持する。
- server response だけでは直前状態を再構築できない場合は、空一覧など layout が壊れない fallback を返す。
- validation error は `validationErrors`、画面全体の処理エラーは `flash.error` に分ける。

## Partial Props

検索画面は render 時の lazy function ではなく、feature handler 内で必要な props だけを eager に組み立てる。

adapter は `X-Inertia-Partial-Component` と `X-Inertia-Partial-Data` を見て partial reload かを判定する。feature handler は `inertia.ShouldIncludeProp` を使い、全体ロードなら全 props、partial reload なら要求された props だけを組み立てる。

```go
props := gonertia.Props{}

if inertia.ShouldIncludeProp(r, component, "initial") {
	result, err := h.usecase.runInitial(r.Context())
	if err != nil {
		return nil, err
	}
	props["initial"] = formatInitial(result)
}
```

usecase error は lazy props 内ではなく、feature handler の `Handle` から返す。これにより render 中に application error を捕捉する必要をなくす。

partial reload 中に displayable error が起きた場合だけ同じ画面を維持するため、feature handler は表示したい props を載せた `PageResult` と error を同時に返す。

```go
props["partialSearch"] = formatPartialSearch(emptyPartialSearchResult(input))
return handlerresult.Page(component, props, handlerresult.PageOptions{
	ValidationErrors: validationErrors,
}), err
```

adapter は `PageResult` と `DisplayableError` の同時返却を partial reload に限って許容する。`PageResult.Props` をそのまま response props に使い、`flash.error` に displayable error message を載せる。

通常 GET の error は `ErrorPage` にする。

## Route Example

現在:

```go
handleAdmin("GET /admin", inertiaApp.Middleware(http.HandlerFunc(container.ShowAdminHandler.Handle)))
handleAdmin("POST /admin/category", http.HandlerFunc(container.StoreCategoryHandler.Handle))
```

移行後:

```go
handleAdmin("GET /admin", handler.InertiaPage(inertiaApp, container.ShowAdminHandler.Handle))
handleAdmin("POST /admin/category", handler.InertiaAction(inertiaApp, container.StoreCategoryHandler.Handle))
```

`BasicAuth` は今と同じく `handleAdmin` 側で適用する。

## Handler Example

GET page:

```go
func (h *Handler) Handle(r *http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError) {
	input, validationErrors, err := toInput(r)
	if err != nil {
		return nil, err
	}

	return handlerresult.Page("admin/ShowAdmin", gonertia.Props{
		"initial": func(ctx context.Context) (any, error) {
			result, err := h.usecase.runInitial(ctx)
			if err != nil {
				return nil, err
			}
			return formatInitial(result), nil
		},
		"partialSearch": func(ctx context.Context) (any, error) {
			if len(validationErrors) > 0 {
				return formatPartialSearch(partialSearchResult{
					Title:       input.Title,
					CategoryIDs: input.CategoryIDs,
					Page:        1,
					TotalCount:  0,
					TotalPages:  1,
				}), nil
			}

			result, err := h.usecase.runPartialSearch(ctx, input)
			if err != nil {
				return nil, err
			}
			return formatPartialSearch(result), nil
		},
	}, handlerresult.PageOptions{
		ValidationErrors: validationErrors,
	}), nil
}
```

GET query validation error では validation error を `error` として early return しない。adapter は component 名や props を知らないため、同じ page response を再構築できない。feature handler が `Page` result を返し、その中に validation errors と画面を壊さない props を含める。

POST action:

```go
func (h *Handler) Handle(r *http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError) {
	input, validationErrors, err := toInput(r)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return handlerresult.RedirectBack(handlerresult.RedirectBackOptions{
			ValidationErrors: validationErrors,
		}), nil
	}

	if err := h.usecase.run(r.Context(), input); err != nil {
		return nil, err
	}
	return handlerresult.Redirect("/admin/category"), nil
}
```

## 移行手順

1. `handlerresult.HandlerResult` と `InertiaPage` / `InertiaAction` adapter を追加する。
2. 既存の `handlererror.ValidationError` と `handlererror.DisplayableError` を継続利用する。
3. `scs` の session manager を導入する。
4. session cookie の name、`HttpOnly`、`Secure`、`SameSite`、lifetime、idle timeout を config に追加する。
5. session に保存した validation errors / flash global error / flash success を GET page adapter で props に展開する。
6. partial reload で error 時にも表示したい props を `PageResult.Props` に載せる。
7. 新規または未完成の POST handler から `func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError)` に移行する。
8. `GET /admin`、`GET /article` など partial reload を持つ画面を adapter 対応に移行する。
9. `GET /article/{articleId}` の not found を `DisplayableError{StatusCode: 404}` に移行する。
10. 各 route から `http.HandlerFunc(...)` と `inertiaApp.Middleware(...)` の直接指定を減らす。
11. integration test で GET / POST / validation / not found / displayable error の期待レスポンスを固定する。

## テスト方針

最低限、以下を integration test で確認する。

- GET page の repository error が通常アクセスでは `ErrorPage` を返す。
- GET partial reload の repository error が画面を壊さない props と `flash.error` を返す。
- GET article detail の不正 ID が `ErrorPage` with `404` を返す。
- GET query validation error の URL 直アクセス / リロードが同じ page の `validationErrors` props を返す。
- GET partial reload の validation error が `validationErrors` props を返す。
- POST validation error が redirect back し、次の GET で `validationErrors` props が出る。
- POST repository error が redirect back し、次の GET で `flash.error` が出る。
- POST success が成功後 URL に redirect し、次の GET で `flash.success` が出る。
- Basic Auth 失敗は `401` と `WWW-Authenticate` header を返す。
- 本番 session cookie が `HttpOnly`、`Secure`、`SameSite=Lax` を満たす。

## 注意点

- feature handler に `http.ResponseWriter` を渡さない。渡すと response を書けてしまい、adapter で統一できなくなる。
- `X-Inertia` だけを redirect back 判定に使わない。Inertia の通常 GET ページ遷移まで redirect back されてしまう。
- 404 は原則 redirect back しない。URL と状態の整合性を優先する。
- GET partial reload の validation error は redirect back しない。URL と検索状態の同期を優先する。
- render / SSR 失敗は flash redirect で隠さない。ログに出し、エラーページまたは静的エラーへ逃がす。
