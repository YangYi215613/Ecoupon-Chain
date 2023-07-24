package main

import (
	"crypto/sha256"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type BlockChain struct {
	Db *bolt.DB // 数据库操作句柄
	Tail string  // 最后一个区块的哈希(从数据库中查找数据)
}

// 获取区块链对象(获取区块链操作句柄)
func NewBlockChain() *BlockChain {
	// 1 读取数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)

	if err != nil {
		logrus.Panic("打开区块链数据库失败...")
	}

	defer db.Close()

	// 定义变量，存储完整区块的hash数据
	var tail []byte

	// 2 读取db 和 tail
	db.View(func(tx *bolt.Tx) error {
		// 找到抽屉，bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			logrus.Warnf("数据抽屉 < %s >不存在，请检查程序...", blockBucket)
		}
		// 获取到完整区块的hash数据
		tail = bucket.Get([]byte(lastBlockHash))
		return nil
	})

	bc = &BlockChain{
		Db:   db,
		Tail: string(tail),
	}

	return bc
}

// 初始化系统区块链
func InitialBlockChain() {
	// 1 生成创世完整区块
	genesisFullBlock := GenesisFullBlock() // 初始化创世区块

	// 2 将该创世区块存储到该编号节点的区块链中
	// a 打开区块链数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)

	if err != nil {
		logrus.Panic("打开区块链数据库失败...")
	}

	defer db.Close()

	// b 将创世完整区块写入区块链数据库中
	db.Update(func(tx *bolt.Tx) error {
		// 找到抽屉，bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			// 创建抽屉
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				logrus.Panic("创建bucket失败...")
			}
		}

		// 写数据
		// key: hash  value: gob数据
		bucket.Put(genesisFullBlock.Hash, genesisFullBlock.Serialize())
		bucket.Put([]byte(lastBlockHash), genesisFullBlock.Hash)
		return nil
	})
}

/*
1) 交易:
a 判断节点是否该自己处理，如果是，将交易加入到自己的内存池中
b 判断该轮交易数量是否足够，如果足够，打包微块，并进行广播


2) 微块:
a 收到微块，加入到节点本地的微块池中
b 判断收到的微块数目是否足够，如果足够，打包完整区块，加入到本地区块链
*/

// 解析收到的命令行信息
func (bc *BlockChain) ParseCMDData(data []byte) {
	//fmt.Printf("节点 < %s > 收到的数据...:%s\n",bc.NodeID, data)
	//fmt.Println("哈哈哈")

	cmd := NewCMD(bc) // 生成命令行对象

	// 将数据进行拆分
	flags := strings.Split(string(data), "_") // 获取方法标识

	// 命令行功能
	switch flags[2] {
	case "getPubKey":
		cmd.getPubKey()
	case "getPrivateKey":
		cmd.getPrivateKey()
	case "getAllPubKey":
		cmd.getAllPubKey()
	case "printBlockInfo":
		// d_printBlockInfo_8000_0
		cmd.printBlockInfo(flags[4]) 
	case "printTransPool":
		// d_printTransPool_8000
		cmd.printTransPool()
	case "printTransNum":
		cmd.printTransNum()
	case "printTransNumAll":
		cmd.printTransNumAll()
	case "printMiniPool":
		cmd.printMiniPool()
	case "printHeight":
		cmd.printHeight()
	case "getUserInfo":
		cmd.getUserInfo(flags[4]) // flags[3]: 节点id
	}
}

// 2 获取接收到的数据并进行解析(交易、微块)
func (bc *BlockChain) ParseData(sourceData []byte) {
	// 进行标识符确认
	flag := string(sourceData[0])

	if string(flag) == "a" {
		// 普通交易
		// infos: ["a", "f677689a-9631-4897-8156-3053f63d573b", "from", "to", "value"]
		infos := strings.Split(string(sourceData), "_") // 获取方法标识

		from, _ := GetNodeKey(infos[2])
		to, _ := GetNodeKey(infos[3])

		value, _ := strconv.Atoi(infos[4])

		currentTransNum += 1 // 每轮交易数目 + 1
		allTransNum += 1  // 总交易数目 + 1

		// c 计算该交易所属委员会编号
		norTrans := NewNorTransaction([]byte(from), []byte(to), uint64(value))

		// 记录交易数据
		allTransMemInfo = append(allTransMemInfo, norTrans)

		// 拿到该节点委员会的信息
		// b 判断该节点是不是每一轮中的委员会领导者，如果不是，直接退出
		if !isCommitteeLeader {
			return
		}

		txCommitteeID := norTrans.GetCommitteeID() // 得到交易所属的委员会编号

		// d 判断领导者节点的编号和交易编号是否相同
		if nodeCommitteeID == int(txCommitteeID) {
			transMemPool = append(transMemPool, norTrans)
			// e 将该交易加入到领导者节点的内存池中
			logrus.Infof("委员会领导者处理交易，委员会编号 < %d >, < %s > 给 < %s > 转账 < %d > ...\n", nodeCommitteeID, from, to, value) // 打印输出
		}

		// f 判断系统中产生的交易数目是否达到最大值，如果是，则打包微块并进行转发
		if currentTransNum >= maxTransNumPerRound {
			// 1) 将内存池中的交易组织成微块
			commiteeLeaderPuk, _ := GetNodeKey(nodeName)
			miniBlock := NewMiniBlock([]byte(commiteeLeaderPuk))
			// 2) 将微块进行转发
			miniBlock.SendMiniBlock()
			// 3) 将该节点的交易池清空
			transMemPool = transMemPool[:0]
			// 4) 将每个节点收到的每轮交易次数置空
			currentTransNum = 0
		}
	} else if flag == "b" {
		// 分发交易
	} else if flag == "c" {
		// 转让交易
	} else if flag == "m" {
		defer func() {
			if r := recover(); r != nil {
				logrus.Info("异常捕获", r)
			}
		}()
		// 反序列化数据
		tempData := strings.Split(string(sourceData), "$$$")

		logrus.Infof("微块数据解析长度为 <%v> \n", len(tempData))

		if len(tempData) != 3 { logrus.Panic("获取微块数据失败...") }

		miniBlock := MiniDeserialize([]byte(tempData[2]))
		
		// 将微块放入内存中，判断微块数目是否足够触发生成完整区块
		nodeAcceptMiniBlock = append(nodeAcceptMiniBlock, miniBlock)

		MiniNumLen := len(nodeAcceptMiniBlock)

		logrus.Infof("节点当前收到 < %d > 个微块数目\n", MiniNumLen)

		// 如果不满足微块重组条件，直接退出
		if MiniNumLen < committeeNum {
			return
		}

		// 1 如果长度为M，则说明收到M个微块，需要组织完整区块
		logrus.Printf("==========> 节点触发微块重组... <=========\n")
		bc.AddFullBlock() // 触发添加完整区块功能
		logrus.Printf("==========> 节点微块重组完成... <=========\n")

		nodeAcceptMiniBlock = nodeAcceptMiniBlock[:0]
		allTransMemInfo = allTransMemInfo[:0]
	}
}

// 3 判断收到的微块个数是否足够 组装完整区块, 将完整区块追加到自己本地的区块链
func (bc *BlockChain) AddFullBlock() {
	// 1 调用数据库，调用数据库获取中的数据(直接获取的数据: round、prevHash)
	db, err := bolt.Open(blockChainDB, 0600, nil)

	if err != nil {
		logrus.Panic("打开数据库失败")
	}

	defer db.Close()

	var miniTimeStamp []uint64 // 存放收到的微块时间戳
	var maxIndex = 0           // 定义最大时间戳索引

	db.Update(func(tx *bolt.Tx) error {
		// a 查找bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			logrus.Panic("bucket不应该为空,请进行检查...")
		}

		// b 获取前一个完整区块
		fullBlockData := bucket.Get([]byte(bc.Tail))
		fullBlock := FullDeserialize(fullBlockData)
		// c 获取round、prevHash
		round := fullBlock.Round + 1
		prevHash := fullBlock.Hash
		prevAccountRoot := fullBlock.AccountRoot

		// 通过前账户树的根hash获取账户树结构
		prevAccountTreeGob := bucket.Get(prevAccountRoot)
		accountTree := AccountDeserialize(prevAccountTreeGob)

		transTree := NewTransTree() // 生成交易树

		var miniLeaderList [][]byte  // 微块领导者

		// 2 需要计算的数据: accountRoot、storageRoot
		// a 遍历该节点内存池中所有的微块，进行账户树的修改
		for _, miniBlock := range nodeAcceptMiniBlock {

			miniLeaderList = append(miniLeaderList, miniBlock.LeaderPuK)

			// 补充: 添加微块的时间戳
			miniTimeStamp = append(miniTimeStamp, miniBlock.TimeStamp)
			// b 获取微块中的普通交易列表
			norTransList := miniBlock.TransTree.NorTransList
			// c 解析每一笔交易
			for _, norTran := range norTransList {
				// --------对账户树进行修改--------
				from := norTran.From // 转账人
				to := norTran.To     // 收款人
				value := norTran.Value // 交易额度

				// d 对上一轮的账户树进行修改
				// 1) 修改转账人的信息
				accountTree.UserInfoMap[string(from)].Balance -= value 
				accountTree.UserInfoMap[string(from)].Count++
				// 2) 修改收款人的信息
				accountTree.UserInfoMap[string(to)].Balance += value

				// ----------组织交易树---------
				transTree.NorTransList = append(transTree.NorTransList, norTransList...)
				// ----------组织存储树、交易树-------------
			}
		}

		// 获取最大的时间戳
		for index, tmp := range miniTimeStamp {
			if tmp > miniTimeStamp[maxIndex] {
				maxIndex = index
			}
		}

		// 取账户树的hash
		accountTreeGob := accountTree.Serialize()
		accountTreeRootHash := sha256.Sum256(accountTreeGob[:20]) // 还需要将账户树进行存储

		// 取交易树的hash
		transTreeRootHash := sha256.Sum256(transTree.Serialize())

		
		logrus.Info("maxIndex: ", miniTimeStamp[maxIndex])

		defer func ()  {
			if r := recover(); r != nil {
				logrus.Info("异常捕获:", r, miniTimeStamp)
			}
		}()

		// 生成完整区块
		// round uint64, prevHash, accountRoot, storageRoot, transRoot, leaderList
		newFullBlock := NewFullBlock(
			round,
			prevHash,
			miniTimeStamp[maxIndex],  // TODO 时间戳出错 index out of range [0] with length 0
			accountTreeRootHash[:],
			nil,  
			transTreeRootHash[:],
			miniLeaderList,
		)

		// 完整区块序列化
		newFullBlockHash := sha256.Sum256(newFullBlock.Serialize())

		// 将数据存入区块链中
		// a 账户树哈希 账户树
		bucket.Put(accountTreeRootHash[:], accountTreeGob)
		// b 交易树哈希 交易树
		bucket.Put(transTreeRootHash[:], transTree.Serialize())
		// c 完整区块哈希 完整区块
		bucket.Put(newFullBlockHash[:], newFullBlock.Serialize())
		// 修改区块链对象的索引
		bc.Tail = string(newFullBlockHash[:])
		// 修改系统中最后指针的索引
		bucket.Put([]byte(lastBlockHash), newFullBlockHash[:])
		return nil
	})
}
