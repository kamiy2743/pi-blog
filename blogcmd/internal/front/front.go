package front

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
	case "install":
		return process.Run(ctx.Path("frontend"), "npm", append([]string{"install"}, args.Rest()...)...)
	default:
		return fmt.Errorf("未対応の front コマンドです: %s", args.First())
	}
}
