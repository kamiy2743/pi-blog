package image

import (
	"fmt"
	"os"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/commandargs"
	"blogcmd/internal/docker"
	"blogcmd/internal/imageref"
	"blogcmd/internal/process"
	"blogcmd/internal/usage"
)

func Run(ctx *appcontext.Context, env string, args commandargs.Args) error {
	if args.IsEmpty() {
		fmt.Fprint(os.Stderr, usage.Text)
		return process.ExitCode(1)
	}
	if args.Len() > 1 {
		return fmt.Errorf("%s image は追加引数を受け取りません", env)
	}

	command := args.First()
	runtimeServices := RuntimeServices()

	switch command {
	case "build":
		if err := docker.ImageBuildCompose(ctx, env, append([]string{"build"}, runtimeServices...)...); err != nil {
			return err
		}
		if err := docker.Run(ctx, "build", "-f", ctx.Path("backend/Dockerfile.migration"), "-t", imageref.For("migration", env), ctx.Path("backend")); err != nil {
			return err
		}
		return docker.Run(ctx, "build", "-f", ctx.Path("backend/Dockerfile.seed"), "-t", imageref.For("seed", env), ctx.Path("backend"))
	case "push":
		if err := docker.ImageBuildCompose(ctx, env, append([]string{"push"}, runtimeServices...)...); err != nil {
			return err
		}
		for _, image := range imageref.JobImages(env) {
			if err := docker.Run(ctx, "push", image); err != nil {
				return err
			}
		}
		return nil
	case "build-push":
		if err := Run(ctx, env, commandargs.Args{"build"}); err != nil {
			return err
		}
		return Run(ctx, env, commandargs.Args{"push"})
	case "pull":
		if err := docker.Compose(ctx, env, append([]string{"pull"}, runtimeServices...)...); err != nil {
			return err
		}
		for _, image := range imageref.JobImages(env) {
			if err := docker.Run(ctx, "pull", image); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("未対応の %s image コマンドです: %s", env, command)
	}
}

func RuntimeServices() []string {
	return []string{"nginx", "go", "ssr"}
}
