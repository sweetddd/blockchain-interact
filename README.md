# blockchain-interact

区块链笔试题
下面的task都是用go语言去实现，然后最后的结果放到github上，到时候把github repo链接发给我们即可。
提示:quicknode 可以创建免费的节点账号，然后可以作为rpc节点(BTC和ETH都有)
ETH主网的rpc节点可以用 https://rpc.ankr.com/eth
ETH测试网的rpc节点可用 https://rpc.ankr.com/eth_goerli

Task1
1.给定一个比特币区块，取出该区块中每笔交易的txhash，inputs(address,amount,index)和outputs(address,amount,index)、txfee，并依次打印上述信息

Task2
2.发起一笔ERC20 token的转账，并广播到链上

Task3
3.给定一个区块高度，完成对应区块的，解析里面所有的ERC20 Trasnfer的交易，并打印出来对应的token地址，以及from、to 、amount
