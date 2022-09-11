package main

import (
	"context"
	"sync"

	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/go-address"
)

var log = logging.Logger("uptime-checker")

type State = MapState

// UptimeChecker maintains the uptime of member nodes
type UptimeChecker struct {
	api v1api.FullNode
	actorAddress address.Address

	rwLock sync.RWMutex
	stop bool
}

func NewUptimeChecker(api v1api.FullNode, actorAddress string) (UptimeChecker, error) {
	addr, err := address.NewFromString(actorAddress)
	if err != nil {
		return UptimeChecker{}, nil
	}
	return UptimeChecker { api: api, actorAddress: addr, stop: false }, nil
}

func (u *UptimeChecker) Start(ctx context.Context) error {
	hasRegistered, err := u.HasRegistered(ctx)
	if err != nil {
		return err
	}

	if !hasRegistered {
		if err := u.Register(ctx); err != nil {
			return err
		}
	} else {
		log.Infow("already registered with the actor, skip register")
	}

	// start the polling of reported checker voting
	go func() {
		u.processReportedCheckers(ctx)
	}()

	// start the checking of member nodes
	go func() {
		u.monitorMemberNodes(ctx)
	}()

	go func() {
		u.monitorCheckerNodes(ctx)
	}()

	return nil
}

// HasRegistered checks if the current checker has already registered itself in the actor
func (u *UptimeChecker) HasRegistered(ctx context.Context) (bool, error) {
	_, err := Load(ctx, u.api, u.actorAddress)
	if err != nil {
		return false, err
	}

	return false, nil
}

// Register registers the current checker to actor
func (u *UptimeChecker) Register(ctx context.Context) error {
	log.Infow("has yet to be registered with the actor, register now")
	return nil
}

// IsStop checks if the up time checker should stop running
func (u *UptimeChecker) IsStop() bool {
	u.rwLock.RLock()
	defer u.rwLock.RUnlock()
	return u.stop
}

// Stop stops the checker
func (u *UptimeChecker) Stop() {
	u.rwLock.Lock()
	defer u.rwLock.Unlock()
	u.stop = true
}

func (u *UptimeChecker) processReportedCheckers(ctx context.Context) {
}

func (u *UptimeChecker) monitorMemberNodes(ctx context.Context) {
}

func (u *UptimeChecker) monitorCheckerNodes(ctx context.Context) {
}