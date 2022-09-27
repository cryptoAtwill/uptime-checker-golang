# Uptime-Checker-Golang
The checker counter-part of `Uptime-Checker`.

## Setup
To setup the FVM actor, please refer to the `uptime-checker` repo: .

Once the actor is deployed, build the source code in the root folder using: `make uptime-checker`.

## Usage
`uptime-checker-golang`will look in `~/.lotus` to connect to a running daemon and resume checking of both nodes and fellow checkers.

For other usage see `./uptime-checker --help`

Before starting the checker, define the following env variable:
```
FULL_NODE=$(./lotus auth api-info --perm admin)
export ${FULL_NODE}
```
Then start the app using `./uptime-checker run ...`.