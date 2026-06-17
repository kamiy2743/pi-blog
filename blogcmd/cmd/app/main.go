package main

import (
	"fmt"
	"os"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/back"
	"blogcmd/internal/commandargs"
	"blogcmd/internal/ent"
	"blogcmd/internal/env"
	"blogcmd/internal/front"
	"blogcmd/internal/process"
	"blogcmd/internal/repo"
	"blogcmd/internal/usage"
)

func main() {
	repoRoot, err := repo.FindRoot()
	if err != nil {
		process.ExitError(err)
	}

	ctx := appcontext.New(repoRoot)
	if err := run(ctx, commandargs.Args(os.Args[1:])); err != nil {
		process.ExitError(err)
	}
}

func run(ctx *appcontext.Context, args commandargs.Args) error {
	if args.IsEmpty() {
		fmt.Fprint(os.Stderr, usage.Text)
		return process.ExitCode(1)
	}

	command := args.First()
	rest := args.Rest()

	switch command {
	case "dev", "stg", "prd":
		return env.Run(ctx, command, rest)
	case "back":
		return back.Run(ctx, rest)
	case "ent":
		return ent.Run(ctx, rest)
	case "front":
		return front.Run(ctx, rest)
	default:
		fmt.Fprint(os.Stderr, usage.Text)
		return process.ExitCode(1)
	}
}
