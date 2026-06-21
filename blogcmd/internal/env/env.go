package env

import (
	"fmt"
	"os"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/commandargs"
	"blogcmd/internal/deploy"
	"blogcmd/internal/docker"
	"blogcmd/internal/imageref"
	"blogcmd/internal/job"
	"blogcmd/internal/process"
	"blogcmd/internal/usage"
)

func Run(ctx *appcontext.Context, env string, args commandargs.Args) error {
	if args.IsEmpty() {
		fmt.Fprint(os.Stderr, usage.Text)
		return process.ExitCode(1)
	}

	action := args.First()
	rest := args.Rest()

	switch action {
	case "up":
		composeArgs := []string{"up", "-d"}
		if imageref.IsImageEnv(env) {
			composeArgs = append(composeArgs, "--no-build")
		} else {
			composeArgs = append(composeArgs, "--build")
		}
		return docker.Compose(ctx, env, append(composeArgs, rest...)...)
	case "down":
		return docker.Compose(ctx, env, append([]string{"down"}, rest...)...)
	case "restart":
		return docker.Compose(ctx, env, append([]string{"restart"}, rest...)...)
	case "recreate":
		if err := docker.Compose(ctx, env, append([]string{"down"}, rest...)...); err != nil {
			return err
		}
		composeArgs := []string{"up", "-d"}
		if imageref.IsImageEnv(env) {
			composeArgs = append(composeArgs, "--no-build")
		} else {
			composeArgs = append(composeArgs, "--build")
		}
		return docker.Compose(ctx, env, append(composeArgs, rest...)...)
	case "deploy":
		if !imageref.IsImageEnv(env) {
			return fmt.Errorf("deploy は stg/prd 専用です")
		}
		return deploy.Run(ctx, env, rest)
	case "mysql":
		return docker.MySQL(ctx, env, rest)
	case "migrate":
		if rest.IsEmpty() {
			fmt.Fprint(os.Stderr, usage.Text)
			return process.ExitCode(1)
		}
		return job.RunMigrate(ctx, env, rest.First())
	case "seed":
		if rest.IsEmpty() {
			fmt.Fprint(os.Stderr, usage.Text)
			return process.ExitCode(1)
		}
		return job.RunSeed(ctx, env, rest.First())
	default:
		return docker.Compose(ctx, env, append([]string{action}, rest...)...)
	}
}
