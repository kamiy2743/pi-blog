package helper

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type InertiaRequest struct {
	Method      string
	Path        string
	QueryParams map[string]string
	Body        io.Reader
}

type InertiaResponse struct {
	StatusCode int            `json:"-"`
	Component  string         `json:"component"`
	Props      map[string]any `json:"props"`
	URL        string         `json:"url"`
}

func RequestInertia(
	t *testing.T,
	server *httptest.Server,
	inertiaRequest InertiaRequest,
) InertiaResponse {
	t.Helper()

	if inertiaRequest.Method == "" {
		t.Fatalf("Inertia リクエストの Method は必須です")
	}

	requestURL, err := url.JoinPath(server.URL, inertiaRequest.Path)
	if err != nil {
		t.Fatalf("Inertia リクエスト URL の組み立てに失敗しました: %v", err)
	}

	if len(inertiaRequest.QueryParams) > 0 {
		values := url.Values{}
		for key, value := range inertiaRequest.QueryParams {
			values.Set(key, value)
		}
		requestURL += "?" + values.Encode()
	}

	req, err := http.NewRequest(inertiaRequest.Method, requestURL, inertiaRequest.Body)
	if err != nil {
		t.Fatalf("Inertia リクエストの作成に失敗しました: %v", err)
	}
	req.Header.Set("X-Inertia", "true")

	res, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("Inertia リクエストの送信に失敗しました: %v", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Inertia レスポンスボディの読み取りに失敗しました: %v", err)
	}

	var inertiaResponse InertiaResponse
	inertiaResponse.StatusCode = res.StatusCode
	if err := json.Unmarshal(bodyBytes, &inertiaResponse); err != nil {
		t.Fatalf("Inertia レスポンス JSON のデコードに失敗しました: %v", err)
	}

	return inertiaResponse
}

func (response InertiaResponse) AssertProps(
	t *testing.T,
	expectedComponent string,
	expectedProps map[string]any,
) {
	t.Helper()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", http.StatusOK, response.StatusCode)
	}

	if response.Component != expectedComponent {
		t.Fatalf("component が不正です: expected=%q actual=%q", expectedComponent, response.Component)
	}

	actualPropsJSON, err := json.Marshal(response.Props)
	if err != nil {
		t.Fatalf("取得した props の JSON エンコードに失敗しました: %v", err)
	}

	expectedPropsJSON, err := json.Marshal(expectedProps)
	if err != nil {
		t.Fatalf("期待する props の JSON エンコードに失敗しました: %v", err)
	}

	if string(actualPropsJSON) != string(expectedPropsJSON) {
		actualPrettyJSON, _ := json.MarshalIndent(response.Props, "", "  ")
		expectedPrettyJSON, _ := json.MarshalIndent(expectedProps, "", "  ")
		t.Fatalf("props が不正です:\nexpected:\n%s\nactual:\n%s", expectedPrettyJSON, actualPrettyJSON)
	}
}
