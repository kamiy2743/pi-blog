package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type testArgs struct {
	Package string `json:"package,omitempty" jsonschema:"backend 相対の Go package パス。例: backend/internal/handler/..."`
}

func RegisterTools(server *mcp.Server, runner *Runner) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "back_fmt",
		Description: "ホスト上の blog リポジトリで ./blog back fmt を実行する",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, CommandResult, error) {
		result, err := runner.RunFmt(ctx)
		return formatToolResult("back_fmt", result, err), result, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "back_test",
		Description: "ホスト上の blog リポジトリで ./blog back test [package] を実行する",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args testArgs) (*mcp.CallToolResult, CommandResult, error) {
		result, err := runner.RunTest(ctx, args.Package)
		return formatToolResult("back_test", result, err), result, nil
	})
}

func formatToolResult(toolName string, result CommandResult, err error) *mcp.CallToolResult {
	lines := []string{fmt.Sprintf("%s を実行しました。", toolName)}

	if len(result.Command) > 0 {
		lines = append(lines, fmt.Sprintf("実行コマンド:\n%s", strings.Join(result.Command, " ")))
	}

	if err != nil {
		lines = append(lines, fmt.Sprintf("エラー:\n%v", err))
	}
	if result.ExitCode != 0 {
		lines = append(lines, fmt.Sprintf("終了コード:\n%d", result.ExitCode))
	}
	if result.Stdout != "" {
		lines = append(lines, "標準出力:\n"+result.Stdout)
	}
	if result.Stderr != "" {
		lines = append(lines, "標準エラー出力:\n"+result.Stderr)
	}
	if result.Stdout == "" && result.Stderr == "" {
		lines = append(lines, "標準出力と標準エラー出力は空です。")
	}

	return &mcp.CallToolResult{
		IsError: err != nil,
		Content: []mcp.Content{
			&mcp.TextContent{Text: strings.Join(lines, "\n\n")},
		},
	}
}
