package main

import (
	"context"
	_ "net/http/pprof"

	"fmt"
	"net/http"

	"strings"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/filecoin-project/lotus/lib/lotuslog"
	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/uptime"
	lcli "github.com/filecoin-project/lotus/cli"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

var log = logging.Logger("uptime-checker")

func main() {
	local := []*cli.Command{
		runCmd,
		versionCmd,
	}

	lotuslog.SetupLogLevels()

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
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "info",
			},
		},
		Before: func(cctx *cli.Context) error {
			return logging.SetLogLevelRegex("uptime-checker", cctx.String("log-level"))
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
			Usage:   "The address of the up time checker actor",
			Value:   "",
		},
		&cli.IntFlag{
			Name:    "actor-id",
			EnvVars: []string{"ACTOR_ID"},
			Usage:   "The actor id of the checker",
			Value:   0,
		},
		&cli.StringFlag{
			Name:    "checker-host",
			EnvVars: []string{"CHECKER_HOST"},
			Usage:   "The host of the up time checker actor",
			Value:   "0.0.0.0",
		},
		&cli.StringFlag{
			Name:    "checker-port",
			EnvVars: []string{"CHECKER_PORT"},
			Usage:   "The port of the up time checker actor",
			Value:   "30000",
		},
		&cli.StringFlag{
			Name:    "node-info-port",
			EnvVars: []string{"NODE_INFO_PORT"},
			Usage:   "The port to get info of the nodes",
			Value:   "3000",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := context.Background()

		checkerHost := cctx.String("checker-host")
		checkerPort := cctx.String("checker-port")
		nodeInfoPort := cctx.String("node-info-port")
	
		actorAddress := cctx.String("actor-address")
		self := uptime.ActorID(cctx.Int("actor-id"))

		log.Infow(
			"starting uptime checker",
			"host", checkerHost,
			"port", checkerPort,
			"nodeInfoPort", nodeInfoPort,
		)

		api, closer, err := lcli.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		// peerID, err := api.ID(ctx);
		// if err != nil {
		// 	return err
		// }

		node, ping, addrs, err := setupLibp2p(checkerHost, checkerPort)
		if err != nil {
			return err
		}

		multiAddresses := make([]uptime.MultiAddr, len(addrs))
		for i, addr := range addrs {
			multiAddresses[i] = addr.String()
		}

		checker, err := uptime.NewUptimeChecker(api, actorAddress, multiAddresses, self, node, ping)
		err = checker.Start(ctx)
		if err != nil {
			return err
		}

		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			info, _ := checker.NodeInfoJsonString()
			fmt.Fprint(writer, info)
		})
		err = http.ListenAndServe(":" + nodeInfoPort, nil)
		if err != nil {
			panic(err)
		}

		// shut the node down
		if err = node.Close(); err != nil {
			panic(err)
		}

		return nil
	},
}

func setupLibp2p(checkerHost string, checkerPort string) (host.Host, *ping.PingService, []multiaddr.Multiaddr, error) {
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/" + checkerHost + "/tcp/" + checkerPort),
		libp2p.Ping(false),
	)
	if err != nil {
		return node, nil, make([]multiaddr.Multiaddr, 0), err
	}

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)

	onlyFirst := make([]multiaddr.Multiaddr, 0)
	for _, addr := range addrs {
		if strings.HasPrefix(addr.String(), "/ip4") {
			onlyFirst = append(onlyFirst, addrs[0])
		}
	}

	log.Infow("Listen addresses:", "addrs", addrs, "first", onlyFirst)

	return node, pingService, onlyFirst, nil
}
