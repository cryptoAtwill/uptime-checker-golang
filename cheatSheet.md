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

./lotus chain create-actor bafk2bzaceadu2yc3k2fozpzhieeqcryrjwpaj2lwjeoyv6u6zhvevupehy5yq ewogICAgImlkcyI6IFsKICAgICAgICAiMTJEM0tvb1dLODdVZkNTeTRRc2lTa1lua2VSRUpUQzZDcllhZEhXcTh6V2JRRW5mTGlHdiIKICAgIF0sCiAgICAiY3JlYXRvcnMiOiBbCiAgICAgICAgMTIzCiAgICBdLAogICAgImFkZHJlc3NlcyI6IFsKICAgICAgICBbCiAgICAgICAgICAgICIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iCiAgICAgICAgXQogICAgXQp9

./lotus chain create-actor bafk2bzaceadu2yc3k2fozpzhieeqcryrjwpaj2lwjeoyv6u6zhvevupehy5yq ewogICAgImlkcyI6IFtdLAogICAgImNyZWF0b3JzIjogW10sCiAgICAiYWRkcmVzc2VzIjogW10KfQ==

make uptime-checker
./uptime-checker --lotus-path /home/lotus/.lotus-local-net run --checker-host 0.0.0.0 --actor-address t2xvrdducml35whderlotzgrn6j42tpgxz34rdklq --actor-id 1001

./uptime-checker --lotus-path /home/lotus/.lotus-local-net run --checker-port 3001 --actor-address t2id4hwrn5osc2es6muzoblsgw2la4idd6gklvlhi --actor-id 101

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 6

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 2 ewogICAgImlkIjogIjEyRDNLb29XSzg3VWZDU3k0UXNpU2tZbmtlUkVKVEM2Q3JZYWRIV3E4eldiUUVuZkxpR3YiLAogICAgImNyZWF0b3IiOiAxMjMsCiAgICAiYWRkcmVzc2VzIjogWyIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iXQp9

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 3 ewogICAgImlkIjogIjEyRDNLb29XQ1h6VEhQZUhCZXRudWhOTGVDcmZLN21GZVpoZmVTYXZ2WkhMTEFQNVhtUnkiLAogICAgImFkZHJlc3NlcyI6IFsiL2lwNC8xNzIuMzEuMzUuMTk1L3RjcC8zNTExNy9wMnAvMTJEM0tvb1dDWHpUSFBlSEJldG51aE5MZUNyZks3bUZlWmhmZVNhdnZaSExMQVA1WG1SeSJdCn0=

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 8 eyJjaGVja2VyIjogMTAwMH0=
./lotus chain invoke t24yl7hskmou5uhrcdz3g6prqlbhefivmocnawlda 8 eyJjaGVja2VyIjogMTAwfQ==
```

export LOTUS_PATH=~/.lotus-node-1
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"

nohup ./lotus daemon --genesis=devgen.car --api 30001 > node-1.out &