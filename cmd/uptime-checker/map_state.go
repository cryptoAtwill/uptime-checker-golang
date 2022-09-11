package main

import (
	"context"

	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
)

type MapState struct {
    // The list of node members in the registry
    members map[PeerID]NodeInfo
	// List of checkers registered in the system.
    checkers map[PeerID]NodeInfo
    // Data structure used to signal offline checkers.
    offlineCheckers map[PeerID]Votes
}

func Load(ctx context.Context, api v1api.FullNode, addr address.Address) (MapState, error)  {
	act, err := api.StateGetActor(ctx, addr, types.EmptyTSK)
	if err != nil {
		return MapState{}, err
	}

	var st MapState
	bs := blockstore.NewAPIBlockstore(api)
	cst := cbor.NewCborStore(bs)
	if err := cst.Get(ctx, act.Head, &st); err != nil {
		return MapState{}, err
	}
	return st, nil
}