package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/flashbotsrpc"
	"golang.org/x/time/rate"
	"math/big"
)

// simulate mev searcher activity by generating txs

type BundleAgent struct {
	// auction target
	slot *big.Int

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
