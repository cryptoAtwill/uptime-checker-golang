package main

import (
	"context"
	_ "net/http/pprof"
	"os"

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
		Before: func(cctx *cli.Context) error {
			return logging.SetLogLevelRegex("stats/*", cctx.String("log-level"))
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

		checker, err := uptime.NewUptimeChecker(api, actorAddress, make([]uptime.PeerID, 0), self)
		err = checker.Start(ctx)
		if err != nil {
			return err
		}

		// TODO: start http server

		return nil
	},
}
