package main

import (
	"crypto/ecdsa"
	"math/big"
)

// simulate mev searcher activity by generating txs

type BundleAgent struct {
	// auction target
	slot uint64

	// efficiency
	additionalGas uint64

	// bid parameters
	startingEffGasPrice  *big.Int
	incrementEffGasPrice *big.Int
	bidRate              uint64 // bids per second

	pk *ecdsa.PrivateKey
}

// deploy bid contract

func main() {
}
