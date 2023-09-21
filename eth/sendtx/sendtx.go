package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	eth_goerli_url := "https://rpc.ankr.com/eth_goerli"
	privateKeyHex := "d04dd360acc65336bc904b2561d8bf936b5a746e1d1b315507d58a911b0693fc"
	senderAddressHex := "0x8945c26BE9Ea13e11fEa6A927353fd30507E894B"
	recipientAddressHex := "0xB9e352320E61c5ca8fadB15ab6A5AA1ca917E268"
	contractAddressHex := "0x1E1469ACE0D0313bC1B36C8a920c5f3116b38bB7"
	amount := big.NewInt(1)

	client, err := ethclient.Dial(eth_goerli_url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum node: %v", err)
	}

	defer client.Close()

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	sender := common.HexToAddress(senderAddressHex)
	recipient := common.HexToAddress(recipientAddressHex)
	contract := common.HexToAddress(contractAddressHex)

	contractABI, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// 获取当前nonce值
	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Fatalf("Failed to retrieve nonce: %v", err)
	}

	// 获取gas价格建议值
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to retrieve gas price: %v", err)
	}

	// 创建ERC20代币转账交易
	transferData, err := contractABI.Pack("transfer", recipient, amount)
	if err != nil {
		log.Fatalf("Failed to pack transfer data: %v", err)
	}

	gasLimit := uint64(50000) // 估算的gas限制

	//tx := types.NewTransaction(nonce, contract, big.NewInt(0), gasLimit, gasPrice, transferData)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &contract,
		Value:    big.NewInt(0),
		Data:     transferData,
	})

	// 签名交易
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to retrieve chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

	//receipt, err := client.TransactionReceipt(ctx, txHash)
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatalf("tx mining error:%v\n", err)
	}

	if receipt.Status == 1 {
		fmt.Printf("Transaction succeed\n")
	}
}
