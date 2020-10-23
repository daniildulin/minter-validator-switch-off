# Minter Validator Switch Off Service

## BUILD

- git clone github.com/daniildulin/minter-validator-switch-off

- cd ./minter-validator-switch-off

- run `go mod download`

- run `go build -o ./builds/switcher ./cmd/switch.go`

## USE

### Setup

Setup environments variables in .env or in OS.

NODES_LIST - separated space hosts list which use for a status check

Run `./switcher -gen_tx -m="mnemonic phrase"` for generate switch off transaction.
You have to repeat this step every time when the node has been disabled.

### Run

Just run `./switcher`
