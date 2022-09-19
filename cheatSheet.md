```shell
docker run -it --rm -v `pwd`:/home/lotus -w /home/lotus --entrypoint bash --user 1000 lotus-fvm
docker run -it --rm -v `pwd`:/home/lotus -v /home/lxm/.ssh:/home/lotus/.ssh -w /home/lotus --entrypoint bash --user 1000 lotus-fvm

# https://lotus.filecoin.io/lotus/developers/local-network/

export LOTUS_PATH=~/.lotus-local-net
export LOTUS_MINER_PATH=~/.lotus-miner-local-net
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"

./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false

FULL_NODE=$(./lotus auth api-info --perm admin)
export ${FULL_NODE}

./lotus-miner run --nosync

./lotus chain install-actor uptime_checker.compact.wasm

./lotus chain create-actor bafk2bzaceadu2yc3k2fozpzhieeqcryrjwpaj2lwjeoyv6u6zhvevupehy5yq ewogICAgImlkcyI6IFsKICAgICAgICAiMTJEM0tvb1dLODdVZkNTeTRRc2lTa1lua2VSRUpUQzZDcllhZEhXcTh6V2JRRW5mTGlHdiIKICAgIF0sCiAgICAiY3JlYXRvcnMiOiBbCiAgICAgICAgMTIzCiAgICBdLAogICAgImFkZHJlc3NlcyI6IFsKICAgICAgICBbCiAgICAgICAgICAgICIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iCiAgICAgICAgXQogICAgXQp9

./lotus chain create-actor bafk2bzaceadu2yc3k2fozpzhieeqcryrjwpaj2lwjeoyv6u6zhvevupehy5yq ewogICAgImlkcyI6IFtdLAogICAgImNyZWF0b3JzIjogW10sCiAgICAiYWRkcmVzc2VzIjogW10KfQ==

make uptime-checker
./uptime-checker --lotus-path /home/lotus/.lotus-local-net run --actor-address t24gt3p42wao27jqzyatalcdrbl4oqatexka6axuy

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 6

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 2 ewogICAgImlkIjogIjEyRDNLb29XSzg3VWZDU3k0UXNpU2tZbmtlUkVKVEM2Q3JZYWRIV3E4eldiUUVuZkxpR3YiLAogICAgImNyZWF0b3IiOiAxMjMsCiAgICAiYWRkcmVzc2VzIjogWyIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iXQp9

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 3 ewogICAgImlkIjogIjEyRDNLb29XSzg3VWZDU3k0UXNpU2tZbmtlUkVKVEM2Q3JZYWRIV3E4eldiUUVuZkxpR3YiLAogICAgImNyZWF0b3IiOiAxMjMsCiAgICAiYWRkcmVzc2VzIjogWyIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iXQp9

./lotus chain invoke t2w2pmfs2mu7uxr7vgnoe2fj5xpeomnkd5v2pcixa 8 eyJjaGVja2VyIjogMTAwMH0=
```