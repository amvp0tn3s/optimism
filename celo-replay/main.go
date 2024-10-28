package replay

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"

	"github.com/ethereum-optimism/optimism/op-node/flags"
)

func Main(ctx *cli.Context, close context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	return Setup(context.Background(), ctx)
}

func main() {
	// Set up logger with a default INFO level in case we fail to parse flags,
	// otherwise the final critical log won't show what the parsing error was.
	// oplog.SetupDefaults()

	app := cli.NewApp()
	// app.Version = VersionWithMeta
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Name = "celo-replay"
	app.Usage = "replay"
	app.Description = "replay txa"
	app.Action = cliapp.LifecycleCmd(Main)
	// app.Commands = []*cli.Command{
	// 	{
	// 		Name:        "networks",
	// 		Subcommands: networks.Subcommands,
	// 	},
	// }

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}
