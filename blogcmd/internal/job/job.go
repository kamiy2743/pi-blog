package job

import (
	"fmt"
	"path/filepath"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/docker"
	"blogcmd/internal/imageref"
)

func RunMigrate(ctx *appcontext.Context, env string, command string) error {
	switch command {
	case "up", "down", "refresh":
	default:
		return fmt.Errorf("未対応の migrate コマンドです: %s", command)
	}

	image := fmt.Sprintf("blog-migration:%s", env)
	if imageref.IsImageEnv(env) {
		image = imageref.For("migration", env)
		if err := docker.Run(ctx, "pull", image); err != nil {
			return err
		}
	} else if err := docker.Run(ctx, "build", "-f", ctx.Path("backend/Dockerfile.migration"), "-t", image, ctx.Path("backend")); err != nil {
		return err
	}

	return docker.Run(
		ctx, "run", "--rm",
		"--network", fmt.Sprintf("blog-%s_private", env),
		"--env-file", ctx.Path(fmt.Sprintf("backend/.env.%s", env)),
		"--mount", fmt.Sprintf("type=bind,src=%s,dst=/run/secrets,readonly", ctx.Path(filepath.Join("secrets", env))),
		image, command,
	)
}

func RunSeed(ctx *appcontext.Context, env string, seed string) error {
	image := fmt.Sprintf("blog-seed:%s", env)
	if imageref.IsImageEnv(env) {
		image = imageref.For("seed", env)
		if err := docker.Run(ctx, "pull", image); err != nil {
			return err
		}
	} else if err := docker.Run(ctx, "build", "-f", ctx.Path("backend/Dockerfile.seed"), "-t", image, ctx.Path("backend")); err != nil {
		return err
	}

	return docker.Run(
		ctx, "run", "--rm",
		"--network", fmt.Sprintf("blog-%s_private", env),
		"--env-file", ctx.Path(fmt.Sprintf("backend/.env.%s", env)),
		"--mount", fmt.Sprintf("type=bind,src=%s,dst=/run/secrets,readonly", ctx.Path(filepath.Join("secrets", env))),
		image, seed,
	)
}
