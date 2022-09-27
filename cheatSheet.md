```shell
docker run -it --rm -v `pwd`:/home/lotus -w /home/lotus --network host --entrypoint bash --user 1000 lotus-fvm
docker run -it --rm -v `pwd`:/home/lotus -w /home/lotus --entrypoint bash lotus-fvm


# https://lotus.filecoin.io/lotus/developers/local-network/

export LOTUS_PATH=~/.lotus-node-0
export LOTUS_MINER_PATH=~/.lotus-miner-node-0
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"
nohup ./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false > node-0.out &

export LOTUS_PATH=~/.lotus-node-1
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"
nohup ./lotus daemon  --genesis=devgen.car --api 1235 > node-1.out &

./lotus-miner run --nosync

./lotus chain install-actor uptime_checker.compact.wasm

export MESSAGE_CID=bafk2bzaceadu2yc3k2fozpzhieeqcryrjwpaj2lwjeoyv6u6zhvevupehy5yq
./lotus chain create-actor ${MESSAGE_CID} ewogICAgImlkcyI6IFtdLAogICAgImNyZWF0b3JzIjogW10sCiAgICAiYWRkcmVzc2VzIjogW10KfQ==

export ACTOR_ADDR=t2tgcfonmyyq3tloubry4d64a54wyeuyka47g65mi && FULL_NODE=$(./lotus auth api-info --perm admin) && export ${FULL_NODE}

./uptime-checker run --checker-host 0.0.0.0 --checker-port 30001 --node-info-port 3001 --actor-address ${ACTOR_ADDR} --actor-id 1001 --wallet-index 0
./uptime-checker run --checker-host 0.0.0.0 --checker-port 30002 --node-info-port 3002 --actor-address ${ACTOR_ADDR} --actor-id 1002 --wallet-index 1
./uptime-checker run --checker-host 0.0.0.0 --checker-port 30003 --node-info-port 3003 --actor-address ${ACTOR_ADDR} --actor-id 1003 --wallet-index 2
./uptime-checker run --checker-host 0.0.0.0 --checker-port 30004 --node-info-port 3004 --actor-address ${ACTOR_ADDR} --actor-id 1004 --wallet-index 3

./lotus chain invoke ${ACTOR_ADDR} 3 ewogICAgImlkIjogIjEyRDNLb29XQlV1SnFmNFRoeTVNSER2eXZ1eUczQnAxa1BKQkM3RVc4Mnd2Q3h5ZGlDTUIiLAogICAgImFkZHJlc3NlcyI6IFsiL2lwNC8xNzIuMzEuMzUuMTk1L3RjcC80NjM2My9wMnAvMTJEM0tvb1dCVXVKcWY0VGh5NU1IRHZ5dnV5RzNCcDFrUEpCQzdFVzgyd3ZDeHlkaUNNQiJdCn0=

./lotus chain invoke ${ACTOR_ADDR} 3 ewogICAgImlkIjogIjEyRDNLb29XSkhTQ3R4bWRFWGY5QWROOTVWaFJ3N0ZiSkNWNEFOMXh0Vjd6Q1d2Z3U1b0siLAogICAgImFkZHJlc3NlcyI6IFsiL2lwNC8xNzIuMzEuMzUuMTk1L3RjcC80MDg2OS9wMnAvMTJEM0tvb1dKSFNDdHhtZEVYZjlBZE45NVZoUnc3RmJKQ1Y0QU4xeHRWN3pDV3ZndTVvSyJdCn0=

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 8 eyJjaGVja2VyIjogMTAwMH0=
./lotus chain invoke t24yl7hskmou5uhrcdz3g6prqlbhefivmocnawlda 8 eyJjaGVja2VyIjogMTAwfQ==
```

export LOTUS_PATH=~/.lotus-node-1
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"

nohup ./lotus daemon --genesis=devgen.car --api 30001 > node-1.out &