package action

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"blog/internal/config"
	"blog/internal/handler/session"
	"blog/internal/test/helper"
)

type TestActionRequest struct {
	Method       string
	Path         string
	Body         map[string]any
	UseBasicAuth bool
	Referer      string
}

type TestActionResponse struct {
	StatusCode int
	Location   string
	client     *http.Client
}

func Send(t *testing.T, server *httptest.Server, request TestActionRequest) TestActionResponse {
	t.Helper()

	if request.Method == "" {
		t.Fatalf("Inertia action リクエストの Method は必須です")
	}

	requestURL, err := url.JoinPath(server.URL, request.Path)
	if err != nil {
		t.Fatalf("Inertia action リクエスト URL の組み立てに失敗しました: %v", err)
	}

	req, err := http.NewRequest(request.Method, requestURL, toJSONBody(t, request.Body))
	if err != nil {
		t.Fatalf("Inertia action リクエストの作成に失敗しました: %v", err)
	}
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("Content-Type", "application/json")

	if request.UseBasicAuth {
		req.SetBasicAuth(config.MustGetAdminBasicAuthUser(), config.MustGetAdminBasicAuthPass())
	}

	if request.Referer != "" {
		referer, err := url.JoinPath(server.URL, request.Referer)
		if err != nil {
			t.Fatalf("Referer URL の組み立てに失敗しました: %v", err)
		}
		req.Header.Set("Referer", referer)
	}

	client := newClient(t)
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Inertia action リクエストの送信に失敗しました: %v", err)
	}
	defer res.Body.Close()

	return TestActionResponse{
		StatusCode: res.StatusCode,
		Location:   res.Header.Get("Location"),
		client:     client,
	}
}

func newClient(t *testing.T) *http.Client {
	t.Helper()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookie jar の作成に失敗しました: %v", err)
	}

	return &http.Client{
		Jar:           jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}
}

func toJSONBody(t *testing.T, value any) io.Reader {
	t.Helper()

	body, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("JSON のエンコードに失敗しました: %v", err)
	}
	return bytes.NewReader(body)
}

func (response TestActionResponse) AssertRedirectTo(t *testing.T, expectedLocation string) {
	t.Helper()

	if response.StatusCode != 303 {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", 303, response.StatusCode)
	}

	if response.Location != expectedLocation && !strings.HasSuffix(response.Location, expectedLocation) {
		t.Fatalf("リダイレクト先が不正です: expected=%s actual=%s", expectedLocation, response.Location)
	}
}

func (response TestActionResponse) AssertValidationError(
	t *testing.T,
	server *httptest.Server,
	manager *session.SessionManager,
	expectedMessages map[string]string,
) {
	t.Helper()

	payload := response.getSessionPayload(t, server, manager)
	if payload.ValidationError == nil {
		t.Fatalf("session に validation error が保存されていません")
	}
	helper.AssertEqual(t, expectedMessages, payload.ValidationError.Messages, "validation error が不正です")
}

func (response TestActionResponse) AssertFlashError(
	t *testing.T,
	server *httptest.Server,
	manager *session.SessionManager,
	expectedMessage string,
) {
	t.Helper()

	payload := response.getSessionPayload(t, server, manager)
	if payload.Flash == nil {
		t.Fatalf("session に flash が保存されていません")
	}
	helper.AssertEqual(t, expectedMessage, payload.Flash.Error, "flash error が不正です")
}

func (response TestActionResponse) getSessionPayload(
	t *testing.T,
	server *httptest.Server,
	manager *session.SessionManager,
) session.SessionPayload {
	t.Helper()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("テストサーバー URL の解析に失敗しました: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, server.URL, nil)
	for _, cookie := range response.client.Jar.Cookies(serverURL) {
		req.AddCookie(cookie)
	}

	var payload session.SessionPayload
	manager.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload = manager.GetSessionPayload(r)
	})).ServeHTTP(httptest.NewRecorder(), req)

	return payload
}
