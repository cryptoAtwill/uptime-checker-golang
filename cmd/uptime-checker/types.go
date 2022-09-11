package main

type PeerID = string;
type MultiAddr = string;
type ActorID = uint32;
type ChainEpoch = int64;

type NodeInfo struct {
	// PeerID of the node
	id PeerID
	// The creator of the node. Only creator can modifier other fields of this struct
	creator ActorID
	/// List of multiaddresses exposed by the node
	/// along with the supported healthcheck endpoints.
	///
	/// e.g. [ /ip4/10.1.1.1/quic/8080/p2p/<peer_id>/ping,
	///        /ip4/10.1.1.1/tcp/8081/http/get/healtcheck,
	///      ]
	/// These multiaddresses are signalling that the liveliness
	/// can be checked by using the default libp2p ping protocol
	/// in the first multiaddress, or by sending a GET HTTP
	/// query to the /healtchek endpoint at 10.1.1.1:8081.
	addresses []MultiAddr
}

type Votes struct {
    // Time of the last offline vote received by a checker.
    lastVote ChainEpoch
    // Checkers that have voted
    votes []PeerID
}