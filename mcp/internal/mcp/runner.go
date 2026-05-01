package mcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

const timeout = 15 * time.Minute

var packagePattern = regexp.MustCompile(`^(backend|backend/[A-Za-z0-9._-]+(?:/[A-Za-z0-9._-]+)*)(/\.\.\.)?$`)

type Runner struct {
	repoRoot string
}

type CommandResult struct {
	Command  []string `json:"command"`
	ExitCode int      `json:"exitCode"`
	Stdout   string   `json:"stdout,omitempty"`
	Stderr   string   `json:"stderr,omitempty"`
}

func NewRunner(repoRoot string) *Runner {
	return &Runner{
		repoRoot: repoRoot,
	}
}

func (r *Runner) RunFmt(ctx context.Context) (CommandResult, error) {
	return r.run(ctx, []string{"back", "fmt"})
}

func (r *Runner) RunTest(ctx context.Context, packagePath string) (CommandResult, error) {
	args := []string{"back", "test"}
	if packagePath != "" {
		if err := validatePackagePath(packagePath); err != nil {
			return CommandResult{}, err
		}
		args = append(args, packagePath)
	}
	return r.run(ctx, args)
}

func (r *Runner) run(ctx context.Context, args []string) (CommandResult, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	command := append([]string{"bash", filepath.Join(r.repoRoot, "blog")}, args...)
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Dir = r.repoRoot

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := CommandResult{
		Command: command,
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
	}

	if err == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if ok := errors.As(err, &exitErr); ok {
		result.ExitCode = exitErr.ExitCode()
		return result, fmt.Errorf("コマンドが終了コード %d で失敗しました", result.ExitCode)
	}

	if ctx.Err() != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return result, fmt.Errorf("コマンドがタイムアウトしました")
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return result, fmt.Errorf("コマンドがキャンセルされました")
		}
		return result, fmt.Errorf("コマンドの実行コンテキストでエラーが発生しました: %w", ctx.Err())
	}

	return result, fmt.Errorf("コマンドの実行に失敗しました: %w", err)
}

func validatePackagePath(packagePath string) error {
	if !packagePattern.MatchString(packagePath) {
		return fmt.Errorf("package は backend 相対の package path 形式で指定してください: %q", packagePath)
	}
	return nil
}
