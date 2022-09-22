package main

import (
	"context"
	_ "net/http/pprof"

	"fmt"
	"io"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"strconv"

	logging "github.com/ipfs/go-log/v2"
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
		&cli.IntFlag{
			Name:    "checker-port",
			EnvVars: []string{"CHECKER_PORT"},
			Usage:   "The port of the up time checker actor",
			Value:   3000,
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := context.Background()

		checkerHost := cctx.String("checker-host")
		checkerPort := cctx.Int("checker-port")

		actorAddress := cctx.String("actor-address")
		self := uptime.ActorID(cctx.Int("actor-id"))

		api, closer, err := lcli.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		peerID, err := api.ID(ctx);
		if err != nil {
			return err
		}

		node, ping, addrs, err := setupLibp2p(peerID, checkerHost, checkerPort)
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

		// http.HandleFunc("/", getPing)

		// err := http.ListenAndServe(":3333", nil)

		// wait for a SIGINT or SIGTERM signal
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info("Received signal, shutting down...")

		// shut the node down
		if err := node.Close(); err != nil {
			panic(err)
		}

		return nil
	},
}

func setupLibp2p(peerID peerstore.ID, hostStr string, port int) (host.Host, *ping.PingService, []multiaddr.Multiaddr, error) {
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/" + hostStr + "/tcp/" + strconv.Itoa(port)),
		libp2p.Ping(false),
	)
	if err != nil {
		return node, nil, make([]multiaddr.Multiaddr, 0), err
	}

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	peerInfo := peerstore.AddrInfo{
		ID:    peerID,
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)

	log.Infow("Listen addresses:", "addrs", addrs)
	return node, pingService, addrs, nil
}

func getPong(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Pong!\n")
}