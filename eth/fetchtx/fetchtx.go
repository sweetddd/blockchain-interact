package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("https://rpc.ankr.com/eth")
	if err != nil {
		log.Fatal(err)
	}

	// 填写要查询的区块高度
	blockNumber := big.NewInt(9818181)

	// 获取指定区块
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	input := "Transfer(address,address,uint256)"

	// 计算Keccak-256哈希值
	erc20TransferHash := crypto.Keccak256([]byte(input))

	// 转换为十六进制表示
	erc20TransferTopic := "0x" + hex.EncodeToString(erc20TransferHash)

	if len(block.Transactions()) == 0 {
		fmt.Printf("block %d does not have transactions", blockNumber)
	} else {
		// 遍历区块中的交易
		for i, tx := range block.Transactions() {
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				log.Fatal(err)
			}

			// 遍历收据中的日志事件
			for _, logEntry := range receipt.Logs {

				// 检查日志topic是否匹配 ERC-20 转账topic
				if logEntry.Topics[0] == common.HexToHash(erc20TransferTopic) {
					// 获取 ERC-20 转账的相关信息
					from := common.HexToAddress(logEntry.Topics[1].Hex())
					to := common.HexToAddress(logEntry.Topics[2].Hex())
					amount := new(big.Int).SetBytes(logEntry.Data)

					fmt.Printf("tx%d: hash:%s\n", i, tx.Hash().Hex())
					fmt.Printf("Token 地址: %s\n", logEntry.Address.Hex())
					fmt.Printf("From 地址: %s\n", from.Hex())
					fmt.Printf("To 地址: %s\n", to.Hex())
					fmt.Printf("Amount: %s\n\n", amount.String())
				}
			}
		}
	}
}
