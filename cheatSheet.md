```shell
docker run -it --rm -v `pwd`:/home/lotus -w /home/lotus --network host --entrypoint bash --user 1000 lotus-fvm
docker run -it --rm -v `pwd`:/home/lotus -w /home/lotus --entrypoint bash lotus-fvm


# https://lotus.filecoin.io/lotus/developers/local-network/

export LOTUS_PATH=~/.lotus-node-0
export LOTUS_PATH=~/.lotus-local-net
export LOTUS_PATH=~/.lotus-follower
export LOTUS_MINER_PATH=~/.lotus-miner-node-0
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"

nohup ./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false > node-0.out &

FULL_NODE=$(./lotus auth api-info --perm admin)
export ${FULL_NODE}

./lotus-miner run --nosync

./lotus chain install-actor uptime_checker.compact.wasm

export MESSAGE_CID=bafk2bzaceadu2yc3k2fozpzhieeqcryrjwpaj2lwjeoyv6u6zhvevupehy5yq
./lotus chain create-actor ${MESSAGE_CID} ewogICAgImlkcyI6IFtdLAogICAgImNyZWF0b3JzIjogW10sCiAgICAiYWRkcmVzc2VzIjogW10KfQ==

export ACTOR_ADDR=

./uptime-checker --lotus-path /home/lotus/.lotus-local-net run --checker-host 0.0.0.0 --actor-address ${ACTOR_ADDR} --actor-id 100
./uptime-checker --lotus-path /home/lotus/.lotus-local-net run --checker-host 0.0.0.0 --checker-port 30001 --node-info-port 3001 --actor-address ${ACTOR_ADDR} --actor-id 1001

./lotus chain invoke ${ACTOR_ADDR} 3 ewogICAgImlkIjogIjEyRDNLb29XQ1h6VEhQZUhCZXRudWhOTGVDcmZLN21GZVpoZmVTYXZ2WkhMTEFQNVhtUnkiLAogICAgImFkZHJlc3NlcyI6IFsiL2lwNC8xNzIuMzEuMzUuMTk1L3RjcC8zNTExNy9wMnAvMTJEM0tvb1dDWHpUSFBlSEJldG51aE5MZUNyZks3bUZlWmhmZVNhdnZaSExMQVA1WG1SeSJdCn0=

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 8 eyJjaGVja2VyIjogMTAwMH0=
./lotus chain invoke t24yl7hskmou5uhrcdz3g6prqlbhefivmocnawlda 8 eyJjaGVja2VyIjogMTAwfQ==
```

export LOTUS_PATH=~/.lotus-node-1
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"

nohup ./lotus daemon --genesis=devgen.car --api 30001 > node-1.out &