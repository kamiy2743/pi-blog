package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/commandargs"
	"blogcmd/internal/image"
	"blogcmd/internal/imageref"
	"blogcmd/internal/process"
)

const (
	defaultRemoteRepoRoot    = "/opt/blog"
	defaultRemoteSecretsRoot = "/etc/blog/secrets"
)

func Run(ctx *appcontext.Context, env string, upArgs commandargs.Args) error {
	if err := image.Run(ctx, env, commandargs.Args{"build-push"}); err != nil {
		return err
	}

	if isLocalDeployEnv(env) {
		return runPullUp(ctx, env, upArgs)
	}

	host := envValue("BLOG_DEPLOY_HOST", defaultRemoteHost(env))
	repoRoot := envValue("BLOG_DEPLOY_REPO_ROOT", defaultRemoteRepoRoot)
	secretsRoot := envValue("BLOG_DEPLOY_SECRETS_ROOT", defaultRemoteSecretsRoot)
	if err := syncRemoteFiles(ctx, env, host, repoRoot, secretsRoot); err != nil {
		return err
	}
	return process.Run(ctx.RepoRoot, "ssh", host, pullUpCommand(env, repoRoot, upArgs))
}

func syncRemoteFiles(ctx *appcontext.Context, env string, host string, repoRoot string, secretsRoot string) error {
	mkdirCommand := "mkdir -p " + shellQuote(repoRoot) + " " + shellQuote(secretsRoot)
	if err := process.Run(ctx.RepoRoot, "ssh", host, mkdirCommand); err != nil {
		return err
	}

	paths := []string{
		"blog",
		"blogcmd/bin/blogcmd",
		fmt.Sprintf("docker-compose.%s.yml", env),
		fmt.Sprintf("backend/.env.%s", env),
		"frontend/.env",
		fmt.Sprintf("cloudflared/config.%s.yml", env),
	}
	args := append([]string{"-az", "--relative"}, paths...)
	args = append(args, host+":"+repoRoot+"/")
	if err := process.Run(ctx.RepoRoot, "rsync", args...); err != nil {
		return err
	}

	secretsArgs := []string{"-az", filepath.Join("secrets", env) + "/", host + ":" + secretsRoot + "/"}
	if err := process.Run(ctx.RepoRoot, "rsync", secretsArgs...); err != nil {
		return err
	}

	chmodCommand := "chmod 711 " + shellQuote(secretsRoot) + " && chmod 644 " + shellQuote(secretsRoot) + "/*"
	return process.Run(ctx.RepoRoot, "ssh", host, chmodCommand)
}

func runPullUp(ctx *appcontext.Context, env string, upArgs commandargs.Args) error {
	compose := []string{"compose", "-p", "blog-" + env, "-f", ctx.Path(fmt.Sprintf("docker-compose.%s.yml", env))}
	if err := process.Run(ctx.RepoRoot, "docker", append(append([]string{}, compose...), append([]string{"pull"}, image.RuntimeServices()...)...)...); err != nil {
		return err
	}
	for _, image := range imageref.JobImages(env) {
		if err := process.Run(ctx.RepoRoot, "docker", "pull", image); err != nil {
			return err
		}
	}
	return process.Run(ctx.RepoRoot, "docker", append(append([]string{}, compose...), append([]string{"up", "-d", "--no-build"}, upArgs...)...)...)
}

func pullUpCommand(env string, repoRoot string, upArgs commandargs.Args) string {
	compose := []string{
		"docker", "compose",
		"-p", "blog-" + env,
		"-f", filepath.Join(repoRoot, fmt.Sprintf("docker-compose.%s.yml", env)),
	}

	commands := [][]string{
		append(append([]string{}, compose...), append([]string{"pull"}, image.RuntimeServices()...)...),
	}
	for _, image := range imageref.JobImages(env) {
		commands = append(commands, []string{"docker", "pull", image})
	}
	commands = append(commands, append(append([]string{}, compose...), append([]string{"up", "-d", "--no-build"}, upArgs...)...))

	parts := []string{"cd " + shellQuote(repoRoot)}
	for _, command := range commands {
		parts = append(parts, shellJoin(command))
	}
	return strings.Join(parts, " && ")
}

func isLocalDeployEnv(env string) bool {
	return env == "stg"
}

func defaultRemoteHost(env string) string {
	if env == "prd" {
		return "eq14-prd"
	}
	return "eq14-dev"
}

func envValue(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func shellJoin(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, shellQuote(arg))
	}
	return strings.Join(quoted, " ")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}
