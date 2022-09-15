package uptime

import (
	"bytes"
	"time"
	"strconv"
	"encoding/json"
)

func encodeJson(payload interface{}) ([]byte, error) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(payload)
	return reqBodyBytes.Bytes(), err
}

// Checks is up and also record the latency
func isUp(addr MultiAddr) UpInfo {
	// TODO: with libp2p
	return UpInfo{
		isOnline: true,
		latency: uint64(0),
		checkedTime: uint64(time.Now().Unix()),
	}
}

func allUp(infos *[]UpInfo) bool {
	for _, info := range *infos {
		if !info.isOnline {
			return false
		}
	}
	return true
}

func keysOfMap(target *map[PeerID]NodeInfo) []PeerID {
	keys := make([]PeerID, len(*target))

	i := 0
	for k := range *target {
		keys[i] = k
		i++
	}
	return keys
}

func parseActorIDFromString(s string) (ActorID, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return ActorID(0), err
	}
	return ActorID(v), nil
}
