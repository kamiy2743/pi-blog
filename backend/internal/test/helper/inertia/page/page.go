package page

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

	"github.com/romsar/gonertia/v3"
)

type TestPageRequest struct {
	Path             string
	QueryParams      map[string][]string
	UseBasicAuth     bool
	PartialComponent string
	PartialData      []string
}

type TestPageResponse struct {
	StatusCode int            `json:"-"`
	Component  string         `json:"component"`
	Props      gonertia.Props `json:"props"`
	URL        string         `json:"url"`
}

func Send(t *testing.T, server *httptest.Server, request TestPageRequest) TestPageResponse {
	t.Helper()

	requestURL, err := url.JoinPath(server.URL, request.Path)
	if err != nil {
		t.Fatalf("Inertia page リクエスト URL の組み立てに失敗しました: %v", err)
	}

	if len(request.QueryParams) > 0 {
		values := url.Values{}
		for key, valueList := range request.QueryParams {
			for _, value := range valueList {
				values.Add(key, value)
			}
		}
		requestURL += "?" + values.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		t.Fatalf("Inertia page リクエストの作成に失敗しました: %v", err)
	}
	req.Header.Set("X-Inertia", "true")
	if request.PartialComponent != "" {
		req.Header.Set("X-Inertia-Partial-Component", request.PartialComponent)
	}
	if len(request.PartialData) > 0 {
		req.Header.Set("X-Inertia-Partial-Data", strings.Join(request.PartialData, ","))
	}

	if request.UseBasicAuth {
		req.SetBasicAuth(config.MustGetAdminBasicAuthUser(), config.MustGetAdminBasicAuthPass())
	}

	res, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("Inertia page リクエストの送信に失敗しました: %v", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Inertia page レスポンスボディの読み取りに失敗しました: %v", err)
	}

	var response TestPageResponse
	response.StatusCode = res.StatusCode
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		t.Fatalf("Inertia page レスポンス JSON のデコードに失敗しました: %v", err)
	}

	return response
}

func (response TestPageResponse) AssertFullProps(
	t *testing.T,
	expectedComponent string,
	expectedProps gonertia.Props,
) {
	t.Helper()

	response.AssertResponse(t, 200, expectedComponent, expectedProps)
}

func (response TestPageResponse) AssertError(
	t *testing.T,
	expectedStatusCode int,
	expectedMessage string,
) {
	t.Helper()

	response.AssertResponse(t, expectedStatusCode, "ErrorPage", gonertia.Props{
		"statusCode": float64(expectedStatusCode),
		"statusText": http.StatusText(expectedStatusCode),
		"message":    expectedMessage,
	})
}

func (response TestPageResponse) AssertResponse(
	t *testing.T,
	expectedStatusCode int,
	expectedComponent string,
	expectedProps gonertia.Props,
) {
	t.Helper()

	if response.StatusCode != expectedStatusCode {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", expectedStatusCode, response.StatusCode)
	}

	if response.Component != expectedComponent {
		t.Fatalf("component が不正です: expected=%q actual=%q", expectedComponent, response.Component)
	}

	expectedResponseProps := make(gonertia.Props, len(expectedProps)+3)
	for key, value := range expectedProps {
		expectedResponseProps[key] = value
	}
	expectedResponseProps["errors"] = gonertia.Props{}
	if _, ok := expectedResponseProps["validationErrors"]; !ok && response.Props["validationErrors"] != nil {
		expectedResponseProps["validationErrors"] = gonertia.Props{}
	}
	if _, ok := expectedResponseProps["flash"]; !ok && response.Props["flash"] != nil {
		expectedResponseProps["flash"] = gonertia.Props{}
	}

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

func (response TestPageResponse) AssertPartialProps(
	t *testing.T,
	expectedComponent string,
	propPath string,
	expectedProps gonertia.Props,
) {
	t.Helper()

	if response.StatusCode != 200 {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", 200, response.StatusCode)
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

func (response TestPageResponse) AssertPropsCount(
	t *testing.T,
	componentName string,
	propName string,
	expectedCount int,
) {
	t.Helper()

	if response.StatusCode != 200 {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", 200, response.StatusCode)
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

func (response TestPageResponse) AssertPropsValue(
	t *testing.T,
	componentName string,
	propPath string,
	expectedValue any,
) {
	t.Helper()

	if response.StatusCode != 200 {
		t.Fatalf("ステータスコードが不正です: expected=%d actual=%d", 200, response.StatusCode)
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

func lookupPropPath(props gonertia.Props, propPath string) (any, bool) {
	if propPath == "" {
		return nil, false
	}

	var current any = props
	for _, key := range strings.Split(propPath, ".") {
		currentMap, ok := propMap(current)
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

	expectedMap, ok := propMap(expected)
	if !ok {
		assertJSONValue(t, path, actual, expected)
		return
	}

	actualMap, ok := propMap(actual)
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

func propMap(value any) (map[string]any, bool) {
	switch typedValue := value.(type) {
	case gonertia.Props:
		return map[string]any(typedValue), true
	case map[string]any:
		return typedValue, true
	default:
		return nil, false
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
