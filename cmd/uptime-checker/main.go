package main

import (
	"context"
	_ "net/http/pprof"
	"os"
	"net/http"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/uptime"
	lcli "github.com/filecoin-project/lotus/cli"
)

var log = logging.Logger("uptime-checker")

func main() {
	local := []*cli.Command{
		runCmd,
		versionCmd,
	}

	app := &cli.App{
		Name:    "uptime-checker",
		Usage:   "Checks the uptime of UptimeCheckerActor member nodes",
		Version: build.UserVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lotus-path",
				EnvVars: []string{"LOTUS_PATH"},
				Value:   "~/.lotus", // TODO: Consider XDG_DATA_HOME
			},
			&cli.StringFlag{
				Name:    "log-level",
				EnvVars: []string{"LOTUS_STATS_LOG_LEVEL"},
				Value:   "debug",
			},
		},
		Commands: local,
	}

	if err := app.Run(os.Args); err != nil {
		log.Errorw("exit in error", "err", err)
		os.Exit(1)
		return
	}
}

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Print version",
	Action: func(cctx *cli.Context) error {
		cli.VersionPrinter(cctx)
		return nil
	},
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "actor-address",
			EnvVars: []string{"ACTOR_ADDRESS"},
			Usage:   "The addres of the up time checker actor",
			Value:   "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := context.Background()

		actorAddress := cctx.String("actor-address")
		self := uptime.ActorID(100) // TODO: what to put as self address?

		api, closer, err := lcli.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		multiAddresses := make([]uptime.MultiAddr, 0)
		multiAddresses = append(multiAddresses, uptime.MultiAddr("/ip4/1.1.1.1/tcp/4242/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N"))
		multiAddresses = append(multiAddresses, uptime.MultiAddr("/ip4/2.2.2.2/tcp/4242/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N"))
		checker, err := uptime.NewUptimeChecker(api, actorAddress, multiAddresses, self)
		err = checker.Start(ctx)
		if err != nil {
			return err
		}

		return http.ListenAndServe(":3333", nil)

		// return nil
	},
}
