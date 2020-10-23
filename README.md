# Minter Validator Switch Off Service

## BUILD

- git clone github.com/daniildulin/minter-validator-switch-off

- cd ./minter-validator-switch-off

- run `go mod download`

- run `go build -o ./builds/switcher ./cmd/switch.go`

## USE

### Setup


| env | Description | Example  |
|---  |---          |---       |
| CHAIN_ID | Minter Network chain id    | 1 - Mainnet; 2 - Testnet  |
| NODES_LIST | separated space hosts list which use for a status check. !!! Important !!! Service use gRPC to connect with a node. Port 8842 by default.    | minter-node-1.testnet.minter.network:8842 minter-node-2.testnet.minter.network:8842  |
| ADDRESS  | Control address    | Mx2fbba5ac7af662043233746df101dd09fa43cefe  |
| PUB_KEY  | A node public key    | Mp972bcf14623c05eb737bdaf033b98863e586aaeb3a93985e3cb255300625441b  |
| MISSED_BLOCKS | missed block count    | 5  |


Setup environments variables in .env or in OS.

Run `./switcher -gen_tx -m="mnemonic phrase"` for generate switch off transaction.

You have to repeat this step every time when the node has been disabled.

### Run

Just run `./switcher`
