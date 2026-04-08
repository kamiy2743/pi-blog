package helper

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestInertiaRequest struct {
	Method      string
	Path        string
	QueryParams map[string]string
	Body        io.Reader
}

type TestInertiaResponse struct {
	StatusCode int            `json:"-"`
	Component  string         `json:"component"`
	Props      map[string]any `json:"props"`
	URL        string         `json:"url"`
}

func RequestInertia(
	t *testing.T,
	server *httptest.Server,
	inertiaRequest TestInertiaRequest,
) TestInertiaResponse {
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

	var inertiaResponse TestInertiaResponse
	inertiaResponse.StatusCode = res.StatusCode
	if err := json.Unmarshal(bodyBytes, &inertiaResponse); err != nil {
		t.Fatalf("Inertia レスポンス JSON のデコードに失敗しました: %v", err)
	}

	return inertiaResponse
}

func (response TestInertiaResponse) AssertProps(
	t *testing.T,
	expectedComponent string,
	expectedProps map[string]any,
) {
	t.Helper()

	response.AssertResponse(t, http.StatusOK, expectedComponent, expectedProps, map[string]any{})
}

func (response TestInertiaResponse) AssertError(
	t *testing.T,
	expectedStatusCode int,
	expectedMessage string,
	expectedDescription string,
) {
	t.Helper()

	response.AssertResponse(t, expectedStatusCode, "ErrorPage", map[string]any{
		"statusCode":  float64(expectedStatusCode),
		"statusText":  http.StatusText(expectedStatusCode),
		"message":     expectedMessage,
		"description": expectedDescription,
	}, map[string]any{})
}

func (response TestInertiaResponse) AssertResponse(
	t *testing.T,
	expectedStatusCode int,
	expectedComponent string,
	expectedProps map[string]any,
	expectedErrors map[string]any,
) {
	t.Helper()

	if response.StatusCode != expectedStatusCode {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", expectedStatusCode, response.StatusCode)
	}

	if response.Component != expectedComponent {
		t.Fatalf("component が不正です: expected=%q actual=%q", expectedComponent, response.Component)
	}

	expectedResponseProps := make(map[string]any, len(expectedProps)+1)
	for key, value := range expectedProps {
		expectedResponseProps[key] = value
	}
	expectedResponseProps["errors"] = expectedErrors

	actualPropsJSON, err := json.Marshal(response.Props)
	if err != nil {
		t.Fatalf("取得した props の JSON エンコードに失敗しました: %v", err)
	}

	expectedPropsJSON, err := json.Marshal(expectedResponseProps)
	if err != nil {
		t.Fatalf("期待する props の JSON エンコードに失敗しました: %v", err)
	}

	if string(actualPropsJSON) != string(expectedPropsJSON) {
		actualPrettyJSON, _ := json.MarshalIndent(response.Props, "", "  ")
		expectedPrettyJSON, _ := json.MarshalIndent(expectedResponseProps, "", "  ")
		t.Fatalf("props が不正です:\nexpected:\n%s\nactual:\n%s", expectedPrettyJSON, actualPrettyJSON)
	}
}
