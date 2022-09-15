package main

// TODO: man, how this works with local package?

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/lotus/uptime"
)

func main() {
	if err := gen.WriteMapEncodersToFile("../cbor_gen.go", "uptime",
		// uptime.NodeInfo{},
		uptime.Votes{},
		uptime.HAMTStateInner{},
	); err != nil {
		panic(err)
	}
}