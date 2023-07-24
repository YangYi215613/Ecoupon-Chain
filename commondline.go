package main

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type CMD struct {
	bc *BlockChain
}

func NewCMD(bc *BlockChain) *CMD {
	return &CMD{
		bc: bc,
	}
}

func (cmd *CMD) getPubKey() {
	pubKey, _ := GetNodeKey(nodeName)
	fmt.Printf("<------------------------------------>\n")
	fmt.Printf("该节点的公钥>>>: %s\n",  pubKey)
	fmt.Printf("<----------------------------------->\n\n")
}

func (cmd *CMD) getPrivateKey() {
	_, privateKey := GetNodeKey(nodeName)
	fmt.Printf("<------------------------------------>\n")
	fmt.Printf("节点的私钥>>>: %s\n",  privateKey)
	fmt.Printf("<------------------------------------>\n\n")
}


func(cmd *CMD) getAllPubKey() {
	fmt.Printf("<------------------------------------>\n")
	for _, nodeInfo := range allNodeName {
		pubKey, _ := GetNodeKey(nodeInfo)
		fmt.Printf("节点 < %s > 的公钥>>>: %s\n", nodeInfo, pubKey)
	}
	fmt.Printf("<------------------------------------>\n\n")
}

func (cmd *CMD) printBlockInfo(blockID string) {
	// blockID 为需要打印的区块编号
	// 0 定义变量
	var fullBlock *FullBlock

	// 1 生成区块链迭代器
	blockChainIT := cmd.bc.NewIterator()
	// 2 循环遍历区块链
	blockIndex, _ := strconv.Atoi(blockID) // 获取所要求区块编号

	fullBlock = blockChainIT.Next()

	// 判断blockID是否大于区块最大高度
	if uint64(blockIndex) > fullBlock.Round {
		fmt.Printf("<-------------- 提示信息 -------------->\n")
		fmt.Printf("该节点的区块链中没有高度为< %d >的区块...\n", blockIndex)
		fmt.Printf("<------------------------------------->\n\n")
		return
	}

	for {
		// 判断轮数是否一致
		if fullBlock.Round == uint64(blockIndex) {
			// 打印完整区块信息
			fmt.Printf("<------------- 完整区块信息 ------------>\n")
			fmt.Printf("类型: %s\n", fullBlock.BType)
			fmt.Printf("区块哈希: %x\n", fullBlock.Hash)
			fmt.Printf("轮数: %d\n", fullBlock.Round)
			fmt.Printf("前区块哈希: %x\n", fullBlock.PrevHash)
			fmt.Printf("时间戳: %d\n", fullBlock.TimeStamp)
			fmt.Printf("账户树根哈希: %x\n", fullBlock.AccountRoot)
			fmt.Printf("存储树根哈希: %x\n", fullBlock.StorageRoot)
			fmt.Printf("交易树根哈希: %x\n", fullBlock.TransRoot)
			fmt.Printf("微块领导者列表: %x\n", fullBlock.LeaderList)
			fmt.Printf("<------------------------------------->\n\n")
		}

		if len(blockChainIT.currentHashPointer) == 0 {
			break
		}
		fullBlock = blockChainIT.Next()
	}
}

// 查看交易池
func (cmd *CMD) printTransPool() {
	fmt.Printf("<--------- 该节点的交易池 ------>\n")
	fmt.Println(transMemPool)
	fmt.Printf("<------------------------------>\n\n")
}

// 查看每轮交易数目
func (cmd *CMD) printTransNum() {
	fmt.Printf("<---------该轮已发生交易数 ------>\n")
	fmt.Println(currentTransNum)
	fmt.Printf("<------------------------------------>\n\n")
}


// 查看总交易数目
func (cmd *CMD) printTransNumAll() {
	fmt.Printf("<---------该轮已发生交易数 ------>\n")
	fmt.Println(allTransNum)
	fmt.Printf("<------------------------------------>\n\n")
}


// 查看微块池
func (cmd *CMD) printMiniPool() {
	fmt.Printf("<--------- 该节点收到的微块池 ------>\n")
	fmt.Println(nodeAcceptMiniBlock)
	fmt.Printf("<------------------------------>\n\n")
}

// 查看区块高度
func (cmd *CMD) printHeight() {

	// 生成区块链迭代器
	blockChainIT := cmd.bc.NewIterator()

	var height uint64

	for len(blockChainIT.currentHashPointer) != 0 {
		blockChainIT.Next()
		height++
	}

	fmt.Printf("<----- 节点维护的区块链高度为< %d >--->\n\n", height)
}

// 查看用户信息
func (cmd *CMD) getUserInfo(nodeID string) {
	pubKey, _ := GetNodeKey(nodeName) // 通过id获取用户公钥

	// 1 拿到最后一个区块(使用迭代器)
	blockChainIT := cmd.bc.NewIterator()
	// 2 解析区块
	fullBlock := blockChainIT.Next()
	// 3 获取账户树hash
	accountTreeRoot := fullBlock.AccountRoot
	var accountTreeGob []byte // 定义变量
	// 4 根据hash值找出账户树
	// a) 打开数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)

	if err != nil {
		logrus.Panic("打开数据库失败...")
	}

	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			logrus.Panic("bucket不应该为空，请进行检查...")
		}
		// 根据hash值找出账户树
		accountTreeGob = bucket.Get(accountTreeRoot)
		return nil
	})

	// 5 进行反解析
	accountTree := AccountDeserialize(accountTreeGob)

	var latestRecord []int
	var latestCount int

	for _, norTran := range allTransMemInfo {
		from, value, to := norTran.From, norTran.Value, norTran.To
		if pubKey == string(from) {
			latestRecord = append(latestRecord, 0-int(value))  // TODO 显示有误
			latestCount += 1
		} else if pubKey == string(to) {
			latestRecord = append(latestRecord, int(value))
			latestCount += 1
		}
	}

	// 6 找到该公钥对应的信息
	fmt.Printf("<------------ 账户信息显示 ----------->\n")
	fmt.Printf("账户余额(已确认): %d \n", accountTree.UserInfoMap[pubKey].Balance)
	fmt.Printf("账户交易次数(已确认): %d\n", accountTree.UserInfoMap[pubKey].Count)
	fmt.Printf("账户电子卡券列表(已确认): %v\n\n", len(accountTree.UserInfoMap[pubKey].CouponList))
	fmt.Printf("账户交易列表(未确认): %v \n", latestRecord)
	fmt.Printf("账户交易次数(未确认): %d\n", latestCount)
	fmt.Printf("<----------------------------------->\n\n")
}
