package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/flashbotsrpc"
	"golang.org/x/time/rate"
	"math/big"
	"os"
)

// simulate mev searcher activity by generating txs

type BundleAgent struct {
	// auction target
	slot *big.Int

	// efficiency
	additionalGas uint64

	// bid parameters
	startingEffGasPrice  *big.Int
	incrementEffGasPrice *big.Int
	bidRate              uint64 // bids per second

	pk *ecdsa.PrivateKey
}

func (b *BundleAgent) RunBundleAgent(rpc string, flashbotsRpc string, mevsimAddr common.Address) error {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	bundleAgentAddress := crypto.PubkeyToAddress(b.pk.PublicKey)

	flashbotsClient := flashbotsrpc.New(flashbotsRpc)

	chainid, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	signer := types.NewLondonSigner(chainid)

	mevsim, err := NewMevSim(mevsimAddr, client)
	if err != nil {
		return err
	}
	mevsimSession := MevSimSession{
		Contract: mevsim,
		CallOpts: bind.CallOpts{
			Pending: false,
			From:    bundleAgentAddress,
		},
		TransactOpts: bind.TransactOpts{
			From:   bundleAgentAddress,
			Signer: privateKeySinger(b.pk, signer),
			Value:  big.NewInt(0),
			NoSend: true,
		},
	}

	var (
		lastBlockNumber uint64
		lastEffGasPrice *big.Int
		lastSlotValue   *big.Int
		lastNonce       uint64
	)

	limiter := rate.NewLimiter(rate.Limit(b.bidRate), 1)
	for {
		err = limiter.Wait(context.Background())
		if err != nil {
			return err
		}

		// get current block number
		blockNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			continue
		}
		if blockNumber != lastBlockNumber || lastBlockNumber == 0 {
			fmt.Println("switching to new block", blockNumber)
			lastSlotValue, err = mevsimSession.GetSlot(b.slot)
			if err != nil {
				fmt.Println("error getting slot value", err)
				continue
			}
			lastNonce, err = client.PendingNonceAt(context.Background(), bundleAgentAddress)
			if err != nil {
				fmt.Println("error getting nonce", err)
				continue
			}
			lastEffGasPrice = new(big.Int).Set(b.startingEffGasPrice)
			lastBlockNumber = blockNumber
		} else {
			lastEffGasPrice = lastEffGasPrice.Add(lastEffGasPrice, b.incrementEffGasPrice)
		}

		mevsimSession.TransactOpts.Nonce = big.NewInt(int64(lastNonce))
		mevsimSession.TransactOpts.GasFeeCap = lastEffGasPrice
		mevsimSession.TransactOpts.GasTipCap = lastEffGasPrice
		mevsimSession.TransactOpts.GasLimit = 100000
		tx, err := mevsimSession.Auction(b.slot, lastSlotValue, big.NewInt(int64(blockNumber+1)))
		if err != nil {
			fmt.Println("error sending tx", err)
			continue
		}

		fmt.Println("created tx", tx.Hash().Hex())
		txBytes, err := tx.MarshalBinary()
		if err != nil {
			fmt.Println("error marshalling tx", err)
			continue
		}

		//send tx as a bundle
		callBundleArgs := flashbotsrpc.FlashbotsSendBundleRequest{
			Txs:         []string{fmt.Sprintf("0x%s", common.Bytes2Hex(txBytes))},
			BlockNumber: fmt.Sprintf("0x%x", blockNumber+1),
		}

		_, err = flashbotsClient.FlashbotsSendBundle(b.pk, callBundleArgs)
		if err != nil {
			fmt.Println("error sending bundle", err)
			continue
		}
	}

	return nil
}

var (
	deployCommand = flag.NewFlagSet("deploy", flag.ExitOnError)
	deployRpc     = deployCommand.String("rpc", "http://localhost:8545", "rpc endpoint")
	deployPrivKey = deployCommand.String("privkey", "ea14b05f4419f1dd22081c213b67b11164296e607373445ce972419eb895ce17", "private key")

	runCommand              = flag.NewFlagSet("run", flag.ExitOnError)
	runRpc                  = runCommand.String("rpc", "http://localhost:8545", "rpc endpoint")
	runFlashbotsRpc         = runCommand.String("fb-rpc", "http://localhost:8545", "flashbots rpc endpoint")
	runSlot                 = runCommand.Int64("slot", 0, "slot to bid on")
	runStartEffGasPrice     = runCommand.Int64("startEffGasPrice", 5, "starting effective gas price")
	runIncrementEffGasPrice = runCommand.Int64("incrementEffGasPrice", 1, "increment effective gas price")
	runBidRate              = runCommand.Int64("bidRate", 10, "bids per second")
	runPrivKey              = runCommand.String("privKey", "ea14b05f4419f1dd22081c213b67b11164296e607373445ce972419eb895ce17", "private key")
	runMevSimAddr           = runCommand.String("mevSimAddr", "0xf28AB6dBC5917732e9B1B6bc204ef5680E5805ec", "mev sim address")
)

func ExecuteDeployCmd(args []string) error {
	err := deployCommand.Parse(args[1:])
	if err != nil {
		deployCommand.Usage()
		return err
	}
	privateKey, err := crypto.HexToECDSA(*deployPrivKey)
	if err != nil {
		return err
	}
	_, err = DeployBidContract(*deployRpc, MevSimBytecode, privateKey)
	return err
}

func ExecuteRunCmd(args []string) error {
	err := runCommand.Parse(args[1:])
	if err != nil {
		runCommand.Usage()
		return err
	}

	startEffGasPrice := new(big.Int).Mul(big.NewInt(*runStartEffGasPrice), big.NewInt(1e9))
	incrementEffGasPrice := new(big.Int).Mul(big.NewInt(*runIncrementEffGasPrice), big.NewInt(1e9))
	privateKey, err := crypto.HexToECDSA(*runPrivKey)
	if err != nil {
		return err
	}

	agent := BundleAgent{
		slot:                 big.NewInt(*runSlot),
		additionalGas:        0,
		startingEffGasPrice:  startEffGasPrice,
		incrementEffGasPrice: incrementEffGasPrice,
		bidRate:              uint64(*runBidRate),
		pk:                   privateKey,
	}

	mevSimAddress := common.HexToAddress(*runMevSimAddr)

	err = agent.RunBundleAgent(*runRpc, *runFlashbotsRpc, mevSimAddress)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "  %s [command] [flags]\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Commands:\n")
		_, _ = fmt.Fprintf(os.Stderr, "deploy\n")
		deployCommand.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return
	}

	command := args[0]
	switch command {
	case "deploy":
		err := ExecuteDeployCmd(args)
		if err != nil {
			panic(err)
		}
	case "run":
		err := ExecuteRunCmd(args)
		if err != nil {
			panic(err)
		}
	default:
		flag.Usage()
	}
}
