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

./lotus chain create-actor bafk2bzacear5h2ldcnjv7hyrl2lm3z7b3je74jejdzbrlfvx5pgbtv6rezi4m ewogICAgImlkcyI6IFsKICAgICAgICAiMTJEM0tvb1dLODdVZkNTeTRRc2lTa1lua2VSRUpUQzZDcllhZEhXcTh6V2JRRW5mTGlHdiIKICAgIF0sCiAgICAiY3JlYXRvcnMiOiBbCiAgICAgICAgMTIzCiAgICBdLAogICAgImFkZHJlc3NlcyI6IFsKICAgICAgICBbCiAgICAgICAgICAgICIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iCiAgICAgICAgXQogICAgXQp9

make uptime-checker
./uptime-checker --lotus-path /home/lotus/.lotus-local-net run --actor-address t27wgaqqyepncc5b6xxthgop3hugr3jvejrcqxd4y

./lotus chain invoke t2qhgvyqiozs6i5yalzyvrnwnktzne3axybgmwswq 6 ewogICAgImlkIjogIjEyRDNLb29XSzg3VWZDU3k0UXNpU2tZbmtlUkVKVEM2Q3JZYWRIV3E4eldiUUVuZkxpR3YiLAogICAgImNyZWF0b3IiOiAxMjMsCiAgICAiYWRkcmVzc2VzIjogWyIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iXQp9

./lotus chain invoke t2qhgvyqiozs6i5yalzyvrnwnktzne3axybgmwswq 2 ewogICAgImlkIjogIjEyRDNLb29XSzg3VWZDU3k0UXNpU2tZbmtlUkVKVEM2Q3JZYWRIV3E4eldiUUVuZkxpR3YiLAogICAgImNyZWF0b3IiOiAxMjMsCiAgICAiYWRkcmVzc2VzIjogWyIvaXA0LzcuNy43LjcvdGNwLzQyNDIvcDJwL1FtWXlRU28xYzFZbTdvcld4TFl2Q3JNMkVteEZUQU5mOHdYbW1FN0RXamh4NU4iXQp9
```