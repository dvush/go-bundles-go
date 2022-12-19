package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)

var (
	rpc      = flag.String("rpc", "http://localhost:8545", "rpc url")
	mnemonic = flag.String("mnemonic", "panic keen way shuffle post attract clever country juice point pulp february", "mnemonic")

	deployCommand = flag.NewFlagSet("deploy", flag.ExitOnError)

	fundCommand = flag.NewFlagSet("fund", flag.ExitOnError)
	fundCheck   = fundCommand.Bool("check", false, "only check balances")
	fundAmount  = fundCommand.Int64("amount", 1000000000000000000, "amount to fund")
	fundCount   = fundCommand.Int("count", 10, "number of accounts to fund")

	runCommand              = flag.NewFlagSet("run", flag.ExitOnError)
	runFlashbotsRpc         = runCommand.String("fb-rpc", "http://localhost:8545", "flashbots rpc endpoint")
	runSlots                = runCommand.String("slots", "0,1", "slot to bid on, comma separated list")
	runCount                = runCommand.String("count", "1,1", "number of agents per slot, comma separated list")
	runStartEffGasPrices    = runCommand.String("start-gp", "5,6", "starting effective gas price(gwei), comma separated list")
	runIncrementEffGasPrice = runCommand.String("inc-gp", "1,2", "increment effective gas price(gwei), comma separated list")
	runBidRate              = runCommand.Int64("rate", 10, "bids per second")
	runMevSimAddr           = runCommand.String("mevsim-addr", "0xafcb5f59eca70854780c04f4fdb04198b969b7ea", "mev sim address")
)

func ExecuteDeployCmd(args []string) error {
	err := deployCommand.Parse(args)
	if err != nil {
		deployCommand.Usage()
		return err
	}
	privateKey, _, err := DeriveWallets(*mnemonic, 1)
	if err != nil {
		return err
	}
	_, err = DeployBidContract(*rpc, MevSimBytecode, privateKey)
	return err
}

func ExecuteRunCmd(args []string) error {
	err := runCommand.Parse(args)
	if err != nil {
		runCommand.Usage()
		return err
	}

	//startEffGasPrice := new(big.Int).Mul(big.NewInt(*runStartEffGasPrice), big.NewInt(1e9))
	//incrementEffGasPrice := new(big.Int).Mul(big.NewInt(*runIncrementEffGasPrice), big.NewInt(1e9))
	//privateKey, err := crypto.HexToECDSA(*runPrivKey)
	//if err != nil {
	//	return err
	//}
	//
	//agent := BundleAgent{
	//	slot:                 big.NewInt(*runSlot),
	//	startingEffGasPrice:  startEffGasPrice,
	//	incrementEffGasPrice: incrementEffGasPrice,
	//	bidRate:              uint64(*runBidRate),
	//	pk:                   privateKey,
	//}
	//
	//mevSimAddress := common.HexToAddress(*runMevSimAddr)
	//
	//err = agent.RunBundleAgent(*runRpc, *runFlashbotsRpc, mevSimAddress)
	//if err != nil {
	//	return err
	//}
	//
	return nil
}

func ExecuteFundCmd(args []string) error {
	err := fundCommand.Parse(args)
	if err != nil {
		fundCommand.Usage()
		return err
	}

	targetBalance := big.NewInt(*fundAmount)

	client, err := ethclient.Dial(*rpc)
	if err != nil {
		return err
	}

	masterWallet, agents, err := DeriveWallets(*mnemonic, *fundCount)
	if err != nil {
		return err
	}
	if *fundCheck {
		privateKeys := append([]*ecdsa.PrivateKey{masterWallet}, agents...)
		fmt.Printf("%-42s %-20s %-20s\n", "Address", "Balance(ETH)", "Defficit(ETH)")
		for _, pk := range privateKeys {
			address := crypto.PubkeyToAddress(pk.PublicKey)
			balance, err := client.BalanceAt(context.Background(), address, nil)
			if err != nil {
				return err
			}
			deficit := new(big.Int).Sub(targetBalance, balance)
			if deficit.Cmp(big.NewInt(0)) < 0 {
				deficit = big.NewInt(0)
			}
			fmt.Printf("%-42s %-20s %-20s\n", address.Hex(), WeiToUnit(balance, 1e18).String(), WeiToUnit(deficit, 1e18).String())
		}
		return nil
	}

	fundAmounts := make([]*big.Int, len(agents))
	totalFundAmount := big.NewInt(0)
	for i := 0; i < len(agents); i++ {
		balance, err := client.BalanceAt(context.Background(), crypto.PubkeyToAddress(agents[i].PublicKey), nil)
		if err != nil {
			return err
		}
		sendValue := new(big.Int).Sub(targetBalance, balance)
		if sendValue.Cmp(big.NewInt(0)) < 0 {
			sendValue = big.NewInt(0)
		}
		fundAmounts[i] = sendValue
		totalFundAmount = new(big.Int).Add(totalFundAmount, sendValue)
	}
	fmt.Printf("Total balance needed(eth): %s\n", WeiToUnit(totalFundAmount, 1e18).String())

	balance, err := client.BalanceAt(context.Background(), crypto.PubkeyToAddress(masterWallet.PublicKey), nil)
	if err != nil {
		return err
	}
	if balance.Cmp(totalFundAmount) < 0 {
		return fmt.Errorf("master wallet balance insufficient")
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	signer := types.NewLondonSigner(chainID)
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(masterWallet.PublicKey))
	if err != nil {
		return err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}
	var lastTx *types.Transaction
	for i := 0; i < len(agents); i++ {
		if fundAmounts[i].Cmp(big.NewInt(0)) == 0 {
			continue
		}
		address := crypto.PubkeyToAddress(agents[i].PublicKey)
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasFeeCap: gasPrice,
			GasTipCap: gasTipCap,
			Gas:       21000,
			To:        &address,
			Value:     fundAmounts[i],
		})
		nonce++

		signedTx, err := types.SignTx(tx, signer, masterWallet)
		if err != nil {
			return err
		}
		fmt.Printf("Sending %s to %s, hash: %s\n", WeiToUnit(fundAmounts[i], 1e18).String(), address.Hex(), tx.Hash().Hex())
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return err
		}
		lastTx = signedTx
	}

	if lastTx != nil {
		fmt.Println("Waiting for last transaction to be mined...")
		_, err := bind.WaitMined(context.Background(), client, lastTx)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "  %s [command] [flags]\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "Commands:\n")
		_, _ = fmt.Fprintf(os.Stderr, "run\n")
		runCommand.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "fund\n")
		fundCommand.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "deploy\n")
		deployCommand.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return
	}

	command, commandArgs := args[0], args[1:]
	switch command {
	case "deploy":
		err := ExecuteDeployCmd(commandArgs)
		if err != nil {
			panic(err)
		}
	case "fund":
		err := ExecuteFundCmd(commandArgs)
		if err != nil {
			panic(err)
		}
	case "run":
		err := ExecuteRunCmd(commandArgs)
		if err != nil {
			panic(err)
		}
	default:
		flag.Usage()
	}
}
