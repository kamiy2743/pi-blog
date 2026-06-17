package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/commandargs"
	"blogcmd/internal/process"
)

func Run(ctx *appcontext.Context, args ...string) error {
	return process.Run(ctx.RepoRoot, "docker", args...)
}

func Compose(ctx *appcontext.Context, env string, args ...string) error {
	composeFile := ctx.Path(fmt.Sprintf("docker-compose.%s.yml", env))
	composeArgs := []string{"compose", "-p", "blog-" + env, "-f", composeFile}
	composeArgs = append(composeArgs, args...)
	return Run(ctx, composeArgs...)
}

func ImageBuildCompose(ctx *appcontext.Context, env string, args ...string) error {
	composeArgs := []string{
		"compose", "-p", "blog-" + env,
		"-f", ctx.Path(fmt.Sprintf("docker-compose.%s.yml", env)),
		"-f", ctx.Path(fmt.Sprintf("docker-compose.%s-build.yml", env)),
	}
	composeArgs = append(composeArgs, args...)
	return Run(ctx, composeArgs...)
}

func MySQL(ctx *appcontext.Context, env string, args commandargs.Args) error {
	user, err := readSecret(ctx.Path(filepath.Join("secrets", env, "mysql_user")))
	if err != nil {
		return err
	}
	password, err := readSecret(ctx.Path(filepath.Join("secrets", env, "mysql_user_password")))
	if err != nil {
		return err
	}

	mysqlArgs := []string{"exec", "mysql", "mysql", "-u" + user, "-p" + password}
	mysqlArgs = append(mysqlArgs, args...)
	return Compose(ctx, env, mysqlArgs...)
}

func readSecret(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\r\n"), nil
}
