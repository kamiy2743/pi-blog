package ent

import (
	"fmt"
	"os"

	"blogcmd/internal/appcontext"
	"blogcmd/internal/commandargs"
	"blogcmd/internal/process"
	"blogcmd/internal/usage"
)

func Run(ctx *appcontext.Context, args commandargs.Args) error {
	if args.IsEmpty() {
		fmt.Fprint(os.Stderr, usage.Text)
		return process.ExitCode(1)
	}
	switch args.First() {
	case "generate":
		goBin, err := process.FindCommand("go", "/usr/local/go/bin/go")
		if err != nil {
			return err
		}
		return process.Run(ctx.Path("backend"), goBin, append([]string{"generate", "./internal/db/ent"}, args.Rest()...)...)
	default:
		return fmt.Errorf("未対応の ent コマンドです: %s", args.First())
	}
}
