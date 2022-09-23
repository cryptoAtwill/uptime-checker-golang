package uptime

import (
	"context"
	"sync"
	"fmt"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/go-address"
	chainTypes "github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"

	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/libp2p/go-libp2p-core/host"
	libp2pMultiaddr "github.com/multiformats/go-multiaddr"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

var log = logging.Logger("uptime")

const NEW_CHECKER = 2
const REPORT_CHECKER = 8

// UptimeChecker maintains the uptime of member nodes
type UptimeChecker struct {
	api v0api.FullNode

	self ActorID
	actorAddress address.Address
	
	checkerAddresses []MultiAddr
	nodeAddresses map[ActorID]map[MultiAddr]HealtcheckInfo

	// libp2p ping related
	node host.Host
	ping *ping.PingService

	rwLock sync.RWMutex
	stop bool
}

func NewUptimeChecker(
	api v0api.FullNode,
	actorAddress string,
	checkerAddresses []MultiAddr,
	self ActorID,
	node host.Host,
	ping *ping.PingService,
) (UptimeChecker, error) {
	addr, err := address.NewFromString(actorAddress)
	if err != nil {
		return UptimeChecker{}, nil
	}
	return UptimeChecker {
		api: api,

		self: self,
		actorAddress: addr,

		checkerAddresses: checkerAddresses,
		nodeAddresses: make(map[ActorID]map[MultiAddr]HealtcheckInfo),

		node: node,
		ping: ping,

		stop: false,
	}, nil
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

	go func() {
		u.processReportedCheckers(ctx)
	}()

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
	s, err := Load(ctx, u.api, u.actorAddress, u.self)
	log.Infow("Self actor id", "actorId", u.self)

	if err != nil {
		return false, err
	}
	return s.HasRegistered(u.self)
}

// Register registers the current checker to actor
func (u *UptimeChecker) Register(ctx context.Context) error {
	log.Infow("has yet to be registered with the actor, register now")

	peerID := u.node.ID();
	log.Infow("register new checker with peer id", "peerID", peerID.String())

	params, err := encodeJson(NodeInfo {
		Id: peerID.String(),
		Addresses: u.checkerAddresses,
	})
	if err != nil {
		return err
	}

	fromAddr, err := u.api.WalletDefaultAddress(ctx)
	if err != nil {
		return err
	}

	return u.execSync(ctx, NEW_CHECKER, fromAddr, params)
}

// Reports to the actor that the checker is down
func (u *UptimeChecker) ReportChecker(ctx context.Context, actor ActorID) error {
	log.Infow("report checker as down", "checker", actor)

	params, err := encodeJson(PeerReportPayload {
		Checker: actor,
	})
	if err != nil {
		return err
	}

	fromAddr, err := u.api.WalletDefaultAddress(ctx)
	if err != nil {
		return err
	}

	return u.execSync(ctx, REPORT_CHECKER, fromAddr, params)
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

func (u *UptimeChecker) CheckChecker(ctx context.Context, actorID ActorID, addrs *[]MultiAddr) error {
	infos := u.multiAddrsUp(addrs)
	
	if !allUp(&infos) {
		log.Warnw("actor down, report now", "actorID", actorID)

		state, err := Load(ctx, u.api, u.actorAddress, u.self)
		if err != nil {
			log.Errorw("cannot load state", "err", err)
			return err
		}

		hasVoted, err := state.HasVotedReportedPeer(actorID)
		if !hasVoted {
			return u.ReportChecker(ctx, actorID)
		}
		log.Debugw("has already reported actor", "actor", actorID)
		
	}
	return nil
}

func (u *UptimeChecker) CheckMember(actorID ActorID, addrs *[]MultiAddr) error {
	infos := u.multiAddrsUp(addrs)
	return u.recordMemberHealthInfo(actorID, &infos, addrs)
}

// /// =================== Private Functions ====================

// Records and aggregate on the health info of membership nodes
func (u *UptimeChecker) recordMemberHealthInfo(actorID ActorID, upInfos *[]UpInfo, addrs *[]MultiAddr) error {
	healthInfos, ok := u.nodeAddresses[actorID]
	if !ok {
		healthInfos = make(map[MultiAddr]HealtcheckInfo, len(*upInfos))
	}

	for i, addr := range(*addrs) {
		val, ok := healthInfos[addr]
		if !ok {
			val = HealtcheckInfo{
				HealtcheckAddr: addr,
				
				AvgLatency: (*upInfos)[i].latency,
				LatencyCounts: 1,
				
				IsOnline: (*upInfos)[i].isOnline,
				Latency: (*upInfos)[i].latency,
				LastChecked: (*upInfos)[i].checkedTime,
			}
		} else {
			val.IsOnline = (*upInfos)[i].isOnline
			val.Latency = (*upInfos)[i].latency
			val.LastChecked = (*upInfos)[i].checkedTime

			// moving average calculation
			newCount := val.AvgLatency + 1
			val.AvgLatency = val.AvgLatency * val.AvgLatency / newCount + val.Latency / newCount
			val.AvgLatency = newCount
		}
		healthInfos[addr] = val
	}

	u.nodeAddresses[actorID] = healthInfos

	return nil
}

func (u *UptimeChecker) multiAddrsUp(addrs *[]MultiAddr) []UpInfo {
	upInfos := make([]UpInfo, 0)
	for _, addr := range(*addrs) {

		isCheck := true
		for _, selfAddr := range u.checkerAddresses {
			if selfAddr == addr {
				isCheck = false
				break
			}
		}

		if !isCheck {
			continue
		}

		upInfos = append(upInfos, u.isUp(addr))
	}
	return upInfos
}

func (u *UptimeChecker) processReportedCheckers(ctx context.Context) error {
	for {
		if u.IsStop() {
			break
		}

		state, err := Load(ctx, u.api, u.actorAddress, u.self)
		if err != nil {
			log.Errorw("cannot load state", "err", err)
			continue
		}

		listToCheck, err := state.ListReportedCheckerNotVoted()
		if err != nil {
			log.Errorw("cannot list repored checkers not voted", "err", err)
			continue
		}

		for toCheckPeerID, addrs := range(listToCheck) {
			u.CheckChecker(ctx, toCheckPeerID, addrs)
		}

		time.Sleep(time.Duration(5) * time.Second)
	}

	return nil
}

func (u *UptimeChecker) monitorMemberNodes(ctx context.Context) error {
	for {
		if u.IsStop() {
			break
		}

		state, err := Load(ctx, u.api, u.actorAddress, u.self)
		if err != nil {
			log.Errorw("cannot load state", "err", err)
			continue
		}

		listToCheck, err := state.ListMembers()
		if err != nil {
			log.Errorw("cannot list members", "err", err)
			continue
		}

		for _, toCheckActorID := range listToCheck {
			addrs, err := state.ListMemberMultiAddrs(toCheckActorID)
			if err != nil {
				log.Errorw("cannot list member multi addrs", "actor", toCheckActorID, "err", err)
				continue
			}

			log.Debugw("member info", "actor", toCheckActorID, "addrs", addrs)

			u.CheckMember(toCheckActorID, addrs)
		}

		time.Sleep(time.Duration(5) * time.Second)
	}

	return nil
}

func (u *UptimeChecker) monitorCheckerNodes(ctx context.Context) error {
	for {
		if u.IsStop() {
			break
		}

		state, err := Load(ctx, u.api, u.actorAddress, u.self)
		if err != nil {
			log.Errorw("cannot load state", "err", err)
			continue
		}

		listToCheck, err := state.ListCheckers()
		if err != nil {
			log.Errorw("cannot list members", "err", err)
			continue
		}

		for _, toCheckPeerID := range listToCheck {
			addrs, err := state.ListCheckerMultiAddrs(toCheckPeerID)
			if err != nil {
				log.Errorw("cannot list member multi addrs", "peer", toCheckPeerID, "err", err)
				continue
			}

			log.Debugw("toCheckPeerID addrs", "toCheckPeerID", toCheckPeerID, "addrs", addrs)

			err = u.CheckChecker(ctx, toCheckPeerID, addrs)
			if err != nil {
				log.Errorw("cannot check checker", "peer", toCheckPeerID, "err", err)
			}
		}

		time.Sleep(time.Duration(5) * time.Second)
	}

	return nil
}

func (u *UptimeChecker) execSync(ctx context.Context, method uint32, from address.Address, params []byte) error {
	smsg, err := u.execSyncNoWait(ctx, method, from, params)
	if err != nil {
		return err
	}

	log.Infow("waiting for message to execute...")
	return u.wait(ctx, smsg)
}

func (u *UptimeChecker) execSyncNoWait(ctx context.Context, method uint32, from address.Address, params []byte) (*chainTypes.SignedMessage, error) {
	msg := &chainTypes.Message{
		To:     u.actorAddress,
		From:   from,
		Value:  big.Zero(),
		Method: abi.MethodNum(method),
		Params: params,
	}
	return u.api.MpoolPushMessage(ctx, msg, nil)
}

func (u *UptimeChecker) wait(ctx context.Context, smsg *chainTypes.SignedMessage) (error) {
	wait, err := u.api.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return err
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}

	return nil
}

// Checks is up and also record the latency
func (u *UptimeChecker) isUp(addrStr MultiAddr) UpInfo {
	upInfo := UpInfo{
		isOnline: false,
		latency: uint64(0),
		checkedTime: uint64(time.Now().Unix()),
	}

	addr, err := libp2pMultiaddr.NewMultiaddr(addrStr)
	if err != nil {
		log.Errorw("cannot parse multi addr", "addr", addr)
		return upInfo
	}

	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Errorw("cannot add multi addr", "addr", addr)
		return upInfo
	}

	log.Debugw("addr for peer", "peer", peer)
	if err := u.node.Connect(context.Background(), *peer); err != nil {
		log.Errorw("cannot connect to multi addr", "peer", peer.ID, "err", err, "addr", addr)
		return upInfo
	}

	ch := u.ping.Ping(context.Background(), peer.ID)
	res := <-ch
	log.Debugw("got ping response!", "RTT:", res.RTT, "res", res)

	upInfo.isOnline = true
	upInfo.checkedTime = uint64(time.Now().Unix())

	return upInfo
}

func (u *UptimeChecker) NodeInfo() map[ActorID]map[MultiAddr]HealtcheckInfo {
	return u.nodeAddresses
}

func (u *UptimeChecker) NodeInfoJsonString() (string, error) {
	log.Debugw("node map", "nodes",  u.nodeAddresses)

	data := make(map[ActorID]map[MultiAddr]HealtcheckInfo, 0)

	for k, v := range u.nodeAddresses {
		data[k] = v
	}

	bytes, err := encodeJson(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}