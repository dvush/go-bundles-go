package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

var MevSimBytecode = common.Hex2Bytes("0x6080806040523461001657610116908161001c8239f35b600080fdfe608080604052600480361015601357600080fd5b600091823560e01c9081637eba7ba61460c0575063b73e739914603557600080fd5b606036600319011260bc57803560243591604435430360ad5782825403609e5760018301809311608b57505580808080478181156083575b4190f11560775780f35b604051903d90823e3d90fd5b506108fc606d565b634e487b7160e01b845260119052602483fd5b6040516301b6e1e760e21b8152fd5b6040516341f833ab60e11b8152fd5b5080fd5b9190503460dc57602036600319011260dc576020925035548152f35b8280fdfea264697066735822122011f3931e3e239632427a61782e9a5c917855da6845ce582d20ce37ce417a948e64736f6c63430008110033")

func DeployBidContract(rpc string, bytecode []byte, privKey *ecdsa.PrivateKey) (common.Address, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return common.Address{}, err
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		return common.Address{}, err
	}

	// deployer address
	deployer := crypto.PubkeyToAddress(privKey.PublicKey)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Address{}, err
	}
	priorityFee, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return common.Address{}, err
	}

	// estimate gas limit
	msg := ethereum.CallMsg{
		From: deployer,
		Data: bytecode,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return common.Address{}, err
	}

	deployerBalance, err := client.BalanceAt(context.Background(), deployer, nil)
	if err != nil {
		return common.Address{}, err
	}

	// deployer balance in eth
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))

	fmt.Println("balance", WeiToUnit(deployerBalance, 1e18),
		"fee", WeiToUnit(fee, 1e18),
		"gasLimit", gasLimit,
		"gasPrice(gwei)", WeiToUnit(gasPrice, 1e9),
		"priorityFee(gwei)", WeiToUnit(priorityFee, 1e9))

	// check if deployer has enough balance
	if deployerBalance.Cmp(fee) < 0 {
		return common.Address{}, fmt.Errorf("insufficient balance")
	}

	nonce, err := client.PendingNonceAt(context.Background(), deployer)
	if err != nil {
		return common.Address{}, err
	}

	// create and sign transaction
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Nonce:     nonce,
		GasTipCap: priorityFee,
		GasFeeCap: gasPrice,
		Gas:       gasLimit,
		To:        nil,
		Data:      bytecode,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privKey)
	if err != nil {
		return common.Address{}, err
	}

	fmt.Println("tx hash", signedTx.Hash().Hex())

	// send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Address{}, err
	}

	// wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		return common.Address{}, err
	}

	// check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Address{}, fmt.Errorf("transaction failed")
	}

	// get contract address
	contractAddress := receipt.ContractAddress
	fmt.Println("contract address", contractAddress.Hex())
	return contractAddress, nil
}

func WeiToUnit(wei *big.Int, unit int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(big.NewInt(int64(unit))))
}

func PrivateKeySinger(pk *ecdsa.PrivateKey, signer types.Signer) bind.SignerFn {
	pkAddress := crypto.PubkeyToAddress(pk.PublicKey)
	return func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		if address != pkAddress {
			return nil, errors.New("incorrect signer address")
		}
		return types.SignTx(transaction, signer, pk)
	}
}
