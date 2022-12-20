package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	hdwallet "github.com/ethereum-optimism/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strconv"
	"strings"
)

var (
	MevSimBytecode       = common.Hex2Bytes("6080806040523461001657610116908161001c8239f35b600080fdfe608080604052600480361015601357600080fd5b600091823560e01c9081637eba7ba61460c0575063b73e739914603557600080fd5b606036600319011260bc57803560243591604435430360ad5782825403609e5760018301809311608b57505580808080478181156083575b4190f11560775780f35b604051903d90823e3d90fd5b506108fc606d565b634e487b7160e01b845260119052602483fd5b6040516301b6e1e760e21b8152fd5b6040516341f833ab60e11b8152fd5b5080fd5b9190503460dc57602036600319011260dc576020925035548152f35b8280fdfea264697066735822122011f3931e3e239632427a61782e9a5c917855da6845ce582d20ce37ce417a948e64736f6c63430008110033")
	MevSimDeployGasLimit = uint64(200000)
)

func DeployBidContract(rpc string, bytecode []byte, privKey *ecdsa.PrivateKey) (common.Address, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return common.Address{}, err
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		return common.Address{}, err
	}

	if len(bytecode) == 0 {
		return common.Address{}, fmt.Errorf("bytecode is empty")
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

	deployerBalance, err := client.BalanceAt(context.Background(), deployer, nil)
	if err != nil {
		return common.Address{}, err
	}

	// deployer balance in eth
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(MevSimDeployGasLimit)))

	fmt.Println("balance", WeiToUnit(deployerBalance, 1e18),
		"fee", WeiToUnit(fee, 1e18),
		"gasLimit", MevSimDeployGasLimit,
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
		Gas:       MevSimDeployGasLimit,
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

func privateKeySinger(pk *ecdsa.PrivateKey, signer types.Signer) bind.SignerFn {
	pkAddress := crypto.PubkeyToAddress(pk.PublicKey)
	return func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		if address != pkAddress {
			return nil, errors.New("incorrect signer address")
		}
		return types.SignTx(transaction, signer, pk)
	}
}

func DeriveWallets(mnemonic string, count int) (*ecdsa.PrivateKey, []*ecdsa.PrivateKey, error) {
	if count < 0 {
		return nil, nil, fmt.Errorf("count must be >= 0")
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, nil, err
	}

	var accounts []*ecdsa.PrivateKey
	for i := 0; i <= count; i++ {
		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		account, err := wallet.Derive(path, false)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to derive account %d: %w", i, err)
		}
		privKey, err := wallet.PrivateKey(account)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get private key for account %d: %w", i, err)
		}
		accounts = append(accounts, privKey)
	}
	return accounts[0], accounts[1:], nil
}

func ParseIntList(s string) ([]int, error) {
	var result []int
	for _, v := range strings.Split(s, ",") {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

func ParseFloatList(s string) ([]float64, error) {
	var result []float64
	for _, v := range strings.Split(s, ",") {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}
	return result, nil
}
