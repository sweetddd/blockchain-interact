package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"log"
)

func main() {
	// 连接到比特币节点
	rpcConfig := &rpcclient.ConnConfig{
		Host:         "smart-wispy-glade.btc.discover.quiknode.pro/49b7059cd0b374d01b28837c7f20d0aed5da340a/",
		User:         "sweet",
		Pass:         "Qn!123456",
		HTTPPostMode: true,
		DisableTLS:   false,
	}
	client, err := rpcclient.New(rpcConfig, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// 获取指定高度的区块
	blockHash, err := client.GetBlockHash(800009)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.GetBlock(blockHash)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历区块中的每笔交易
	for txIndex, tx := range block.Transactions {
		if txIndex == 0 {
			fmt.Printf("tx%d: txHash: %s\n", txIndex, tx.TxHash())
			fmt.Printf("Input: coinbase\n")

			output := tx.TxOut[0]
			fmt.Printf("output[0]:\n")

			pkScript := output.PkScript

			// 解析输出脚本以获取地址
			_, addresses, _, err := txscript.ExtractPkScriptAddrs(pkScript, &chaincfg.MainNetParams)
			if err != nil {
				panic(err)
			}
			// 打印输出地址和金额
			for _, addr := range addresses {
				fmt.Printf("   address: %s\n", addr.EncodeAddress())
			}
			outputValue := float64(output.Value) / 1e8
			fmt.Printf("   amount: %f BTC\n\n", outputValue)

			continue
		}

		fmt.Printf("tx%d: txHash: %s\n", txIndex, tx.TxHash())

		// 每笔交易输入的总额
		var inputTotal int64

		// 遍历交易的输入
		for inputIndex, input := range tx.TxIn {

			fmt.Printf("Input[%d]:\n", inputIndex)
			previousOutPoint := input.PreviousOutPoint
			previousTxHash := previousOutPoint.Hash
			previousVoutIndex := previousOutPoint.Index

			// 获取前一笔交易的详细信息
			previousTx, err := client.GetRawTransaction(&previousTxHash)
			if err != nil {
				log.Fatal(err)
			}

			// 获取前一笔交易输出的金额
			previousTxOutput := previousTx.MsgTx().TxOut[previousVoutIndex]
			value := float64(previousTxOutput.Value) / 1e8 // 将Satoshi转换为BTC

			inputTotal += previousTxOutput.Value

			// 获取前一笔交易输出的脚本
			pkScript := previousTxOutput.PkScript

			//解析脚本以获取地址
			_, addresses, _, err := txscript.ExtractPkScriptAddrs(pkScript, &chaincfg.MainNetParams)
			if err != nil {
				panic(err)
			}

			// 打印输出地址和金额
			for _, addr := range addresses {
				fmt.Printf("   address: %s\n", addr.EncodeAddress())
			}

			fmt.Printf("   amount: %f BTC\n", value)
		}

		// 计算输出总额
		var outputTotal int64

		// 处理交易的输出
		for outputIndex, output := range tx.TxOut {
			fmt.Printf("output[%d]:\n", outputIndex)

			// 解析输出脚本
			pkScript := output.PkScript

			// 解析脚本以获取地址
			_, addresses, _, err := txscript.ExtractPkScriptAddrs(pkScript, &chaincfg.MainNetParams)
			if err != nil {
				panic(err)
			}

			// 打印输出地址
			for _, addr := range addresses {
				fmt.Printf("   address: %s\n", addr.EncodeAddress())
			}

			outputValue := float64(output.Value) / 1e8

			fmt.Printf("   amount: %f BTC\n", outputValue)

			outputTotal += output.Value
		}

		// 计算手续费
		fee := float64(inputTotal-outputTotal) / 1e8

		// 输出交易的手续费
		fmt.Printf("txfee: %f BTC\n\n", fee)
	}
}
