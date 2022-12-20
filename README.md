Toy tool to create fake bundle load.

## How to

Initial preparations:
1. Generate mnemonic and fund wallet 0 (hdpath `m/44'/60'/0'/0/0`)
2. Deploy `MevSim.sol` contract using `deploy` command.
goerli address `0xa1c874985ec392209070Db9C613f61fF7F66d23E`

Fund searcher wallets:
1. Decide how many searchers do you want and amount of funds to send to each of them. 
2. Fund these wallets using `fund` command. (these wallets will be generated from hdpath `m/44'/60'/0'/0/i`)
Fund command can also be used to `-check` balances and addresses of these wallets.

Run tests:
1. Start sending bundles with `run` command.

To run you should describe set of searchers. Searchers compete for slots, only one tx per slot can enter the block.
Parameters:
- `-slots 0,1`     - slots to use in the test
- `-count 1,2`     - number of searchers per slot
- `-start-gp 5,5`  - effective gas price in gwei for the first bundle per block
- `-inc-gp 1,1`    - effective gas price increment. searchers will resend bundles with higher effective gas price for the same block
- `-rate 1`        - rate at which new bundles are resent

## Examples

### Goerli
```shell
export MNEMONIC="..."
export ETH_RPC_URL=https://goerli.infura.io/v3/...
export MEVSIM_ADDR=0xa1c874985ec392209070Db9C613f61fF7F66d23E
export FLASHBOTS_RPC_URL=https://relay-goerli.flashbots.net

./go-bundles-go -rpc "$ETH_RPC_URL" -mnemonic "$MNEMONIC" \
                run -fb-rpc "$FLASHBOTS_RPC_URL" \
                    -mevsim-addr "$MEVSIM_ADDR"  \
                    -rate 1 \
                    -count 2,2 \
                    -slots 0,1 \
                    -start-gp 0.01,0.01 \
                    -inc-gp 0.001,0.001
```

### Local devnet

1. Setup local devnet using https://github.com/dvush/geth-builder-local-devnet
2. Deploy with `./go-bundles-go deploy` (no args required - defaults should work)
3. Run with `./go-bundles-go run` (no args required - defaults should work)

## Usage

```
Usage of ./go-bundles-go:
  ./go-bundles-go [command] [flags]
Flags:
  -mnemonic string
    	mnemonic (default "panic keen way shuffle post attract clever country juice point pulp february")
  -rpc string
    	rpc url (default "http://localhost:8545")
Commands:
run
  -count string
    	number of agents per slot, comma separated list (default "1,1")
  -fb-rpc string
    	flashbots rpc endpoint (default "http://localhost:8545")
  -inc-gp string
    	increment effective gas price(gwei), comma separated list (default "1,2")
  -mevsim-addr string
    	mev sim address (default "0xafcb5f59eca70854780c04f4fdb04198b969b7ea")
  -rate uint
    	bids per second (default 10)
  -slots string
    	slot to bid on, comma separated list (default "0,1")
  -start-gp string
    	starting effective gas price(gwei), comma separated list (default "5,6")
fund
  -amount int
    	amount to fund (default 1000000000000000000)
  -check
    	only check balances
  -count int
    	number of accounts to fund (default 10)
deploy
```