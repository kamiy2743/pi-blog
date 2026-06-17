package back

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/commandargs"
	"blogcmd/internal/docker"
	"blogcmd/internal/process"
	"blogcmd/internal/usage"
)

func Run(ctx *appcontext.Context, args commandargs.Args) error {
	if args.IsEmpty() {
		fmt.Fprint(os.Stderr, usage.Text)
		return process.ExitCode(1)
	}
	switch args.First() {
	case "fmt":
		return runFmt(ctx)
	case "mod":
		if args.Len() < 2 {
			fmt.Fprint(os.Stderr, usage.Text)
			return process.ExitCode(1)
		}
		return runMod(ctx, args.At(1), args[2:])
	case "test":
		target, err := normalizeTestTarget(ctx, args.Rest())
		if err != nil {
			return err
		}
		return runTest(ctx, target)
	default:
		return fmt.Errorf("未対応の back コマンドです: %s", args.First())
	}
}

func runFmt(ctx *appcontext.Context) error {
	goBin, err := process.FindCommand("go", "/usr/local/go/bin/go")
	if err != nil {
		return err
	}
	for _, dir := range goModDirs(ctx) {
		if err := process.Run(dir, goBin, "fmt", "./..."); err != nil {
			return err
		}
	}
	return nil
}

func runMod(ctx *appcontext.Context, subcommand string, args commandargs.Args) error {
	if subcommand != "tidy" {
		return fmt.Errorf("未対応の back mod コマンドです: %s", subcommand)
	}
	goBin, err := process.FindCommand("go", "/usr/local/go/bin/go")
	if err != nil {
		return err
	}
	for _, dir := range goModDirs(ctx) {
		if err := process.Run(dir, goBin, append([]string{"mod", "tidy"}, args...)...); err != nil {
			return err
		}
	}
	return nil
}

func runTest(ctx *appcontext.Context, target string) error {
	if err := docker.Compose(ctx, "dev", "up", "-d", "--build", "go-test"); err != nil {
		return err
	}

	script := `
set -euo pipefail

cd /app
if command -v go >/dev/null 2>&1; then
  go_bin="$(command -v go)"
elif [[ -x /usr/local/go/bin/go ]]; then
  go_bin=/usr/local/go/bin/go
else
  echo "go が見つかりません" >&2
  exit 1
fi

mapfile -t test_packages < <(
  "$go_bin" list -f '{{if or (gt (len .TestGoFiles) 0) (gt (len .XTestGoFiles) 0)}}{{.ImportPath}}{{end}}' "$1" |
    sed "/^$/d"
)

if [[ ${#test_packages[@]} -eq 0 ]]; then
  echo "テスト対象の package がありません: $1" >&2
  exit 0
fi

"$go_bin" test -mod=readonly -p 1 "${test_packages[@]}"
`
	return docker.Compose(ctx, "dev", "exec", "-T", "go-test", "bash", "-lc", script, "bash", target)
}

func normalizeTestTarget(ctx *appcontext.Context, args commandargs.Args) (string, error) {
	backendDir := ctx.Path("backend")
	target := "./..."

	if args.IsEmpty() {
		return target, nil
	}
	target = args.First()

	if stat, err := os.Stat(target); err == nil && !stat.IsDir() {
		return "", fmt.Errorf("back test は Go package 単位で実行してください。ファイル指定は未対応です: %s", target)
	}

	switch {
	case target == "backend":
		return ".", nil
	case strings.HasPrefix(target, "backend/"):
		return "./" + strings.TrimPrefix(target, "backend/"), nil
	case target == backendDir:
		return ".", nil
	case strings.HasPrefix(target, backendDir+string(os.PathSeparator)):
		return "." + strings.TrimPrefix(target, backendDir), nil
	default:
		return target, nil
	}
}

func goModDirs(ctx *appcontext.Context) []string {
	var dirs []string
	for _, dir := range []string{ctx.Path("backend"), ctx.Path("mcp"), ctx.Path("blogcmd")} {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			dirs = append(dirs, dir)
		}
	}
	return dirs
}
