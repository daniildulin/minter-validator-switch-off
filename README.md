# Minter Validator Switch Off Service

## BUILD

- git clone github.com/daniildulin/minter-validator-switch-off

- cd ./minter-validator-switch-off

- run `go mod download`

- run `go build -o ./builds/switcher ./cmd/switch.go` if you want to generate a transaction manually

- or `go build -ldflags="-X 'github.com/daniildulin/minter-validator-switch-off/core.Vs=mnemonic phrase here'" -o ./builds/linux/switcher ./cmd/switch.go` if you want to turn it on and forget

## USE

### Setup


| env   | <div style="width:500px">Description</div> | Example   |
|---    |---  |---    |
| CHAIN_ID  | Minter Network chain id   | 1 - Mainnet; 2 - Testnet  |
| NODES_LIST    | separated space hosts list which use for a status check. !!! Important !!! Service use gRPC to connect with a node. Port 8842 by default. | minter-node-1.testnet.minter.network:8842 minter-node-2.testnet.minter.network:8842   |
| ADDRESS   | Control address   | Mx...    |
| PUB_KEY   | A node public key | Mp...    |
| MISSED_BLOCKS | missed block count    | 5 |


Setup environments variables in .env or in OS.

Run `./switcher -gen_tx -m="mnemonic phrase"` for generate switch off transaction.

You have to repeat this step every time when the node has been disabled.

### Run

Just run `./switcher`
