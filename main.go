package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
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

func (b *BundleAgent) RunBundleAgent(rpc string, mevsimAddr common.Address) error {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	bundleAgentAddress := crypto.PubkeyToAddress(b.pk.PublicKey)

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
			Nonce:  nil,
			Signer: PrivateKeySinger(b.pk, signer),
			NoSend: true,
		},
	}

	var (
		lastBlockNumber uint64
		lastEffGasPrice *big.Int
		lastSlotValue   *big.Int
	)
	for {
		// get current block number
		blockNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			continue
		}
		if blockNumber != lastBlockNumber || lastBlockNumber == 0 {
			lastBlockNumber = blockNumber
			lastEffGasPrice = new(big.Int).Set(b.startingEffGasPrice)
			lastSlotValue, err = mevsimSession.GetSlot(b.slot)
			if err != nil {
				continue
			}
		} else {
			lastEffGasPrice = lastEffGasPrice.Add(lastEffGasPrice, b.incrementEffGasPrice)
		}

		tx, err := mevsimSession.Auction(b.slot, lastSlotValue, big.NewInt(int64(blockNumber)))
		if err != nil {
			continue
		}

		// TODO send tx
		tx.To()

	}

	return nil
}

func main() {
}
