package article

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

func ConvertMarkdownToHTML(markdown string) (string, error) {
	var html bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &html); err != nil {
		return "", err
	}
	return bluemonday.UGCPolicy().Sanitize(html.String()), nil
}
