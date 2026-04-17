package helper

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"blog/internal/config"
)

type TestInertiaRequest struct {
	Method      string
	Path        string
	QueryParams map[string][]string
	Body        io.Reader
	UseBasicAuth bool
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
		for key, valueList := range inertiaRequest.QueryParams {
			for _, value := range valueList {
				values.Add(key, value)
			}
		}
		requestURL += "?" + values.Encode()
	}

	req, err := http.NewRequest(inertiaRequest.Method, requestURL, inertiaRequest.Body)
	if err != nil {
		t.Fatalf("Inertia リクエストの作成に失敗しました: %v", err)
	}
	req.Header.Set("X-Inertia", "true")
	if inertiaRequest.UseBasicAuth {
		req.SetBasicAuth(config.MustGetAdminBasicAuthUser(), config.MustGetAdminBasicAuthPass())
	}

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

func (response TestInertiaResponse) AssertFullProps(
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

func (response TestInertiaResponse) AssertNotFound(t *testing.T) {
	t.Helper()

	response.AssertResponse(t, http.StatusNotFound, "NotFound", map[string]any{}, map[string]any{})
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

func (response TestInertiaResponse) AssertPartialProps(
	t *testing.T,
	expectedComponent string,
	propPath string,
	expectedProps map[string]any,
) {
	t.Helper()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", http.StatusOK, response.StatusCode)
	}

	if response.Component != expectedComponent {
		t.Fatalf("component が不正です: expected=%q actual=%q", expectedComponent, response.Component)
	}

	actualProps, ok := lookupPropPath(response.Props, propPath)
	if !ok {
		t.Fatalf("props に %q が存在しません", propPath)
	}

	assertPartialJSON(t, "props."+propPath, actualProps, expectedProps)
}

func (response TestInertiaResponse) AssertPropsCount(
	t *testing.T,
	componentName string,
	propName string,
	expectedCount int,
) {
	t.Helper()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", http.StatusOK, response.StatusCode)
	}

	if response.Component != componentName {
		t.Fatalf("component が不正です: expected=%q actual=%q", componentName, response.Component)
	}

	prop, ok := lookupPropPath(response.Props, propName)
	if !ok {
		t.Fatalf("props に %q が存在しません", propName)
	}

	value := reflect.ValueOf(prop)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		t.Fatalf("props[%q] は件数を検証できる型ではありません: actual=%T", propName, prop)
	}

	if value.Len() != expectedCount {
		t.Fatalf("props[%q] の件数が不正です: expected=%d actual=%d", propName, expectedCount, value.Len())
	}
}

func (response TestInertiaResponse) AssertPropsValue(
	t *testing.T,
	componentName string,
	propPath string,
	expectedValue any,
) {
	t.Helper()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", http.StatusOK, response.StatusCode)
	}

	if response.Component != componentName {
		t.Fatalf("component が不正です: expected=%q actual=%q", componentName, response.Component)
	}

	actualValue, ok := lookupPropPath(response.Props, propPath)
	if !ok {
		t.Fatalf("props に %q が存在しません", propPath)
	}

	assertJSONValue(t, "props."+propPath, actualValue, expectedValue)
}

func lookupPropPath(props map[string]any, propPath string) (any, bool) {
	if propPath == "" {
		return nil, false
	}

	var current any = props
	for _, key := range strings.Split(propPath, ".") {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = currentMap[key]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

func assertPartialJSON(t *testing.T, path string, actual any, expected any) {
	t.Helper()

	expectedMap, ok := expected.(map[string]any)
	if !ok {
		assertJSONValue(t, path, actual, expected)
		return
	}

	actualMap, ok := actual.(map[string]any)
	if !ok {
		t.Fatalf("%s は object ではありません: actual=%T", path, actual)
	}

	for key, expectedValue := range expectedMap {
		actualValue, ok := actualMap[key]
		if !ok {
			t.Fatalf("%s.%s が存在しません", path, key)
		}

		assertPartialJSON(t, path+"."+key, actualValue, expectedValue)
	}
}

func assertJSONValue(t *testing.T, path string, actual any, expected any) {
	t.Helper()

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("取得した %s の JSON エンコードに失敗しました: %v", path, err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("期待する %s の JSON エンコードに失敗しました: %v", path, err)
	}

	if string(actualJSON) != string(expectedJSON) {
		t.Fatalf("%s の値が不正です: expected=%s actual=%s", path, expectedJSON, actualJSON)
	}
}
