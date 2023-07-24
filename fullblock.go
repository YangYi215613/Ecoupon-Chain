package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

// 完整区块
type FullBlock struct {
	BType       []byte // 区块类型
	Round       uint64 // 系统轮次
	PrevHash    []byte // 前完整区块hash
	TimeStamp   uint64 // 时间戳
	AccountRoot []byte // 账户树根哈希(账户MPT树)
	StorageRoot []byte // 存储树根哈希(Merkle树)
	TransRoot   []byte // 交易树根哈希(Merkle树)
	LeaderList  [][]byte // 微块领导者信息列表(先置为空)

	// 附带数据(方便操作)
	// a 当前区块hash
	Hash []byte
}

// 生成完整区块
// 执行前提: 收到M个微块才进行重组
func NewFullBlock(round uint64, prevHash []byte, miniMaxTimeStamp uint64, accountRoot, storageRoot, transRoot []byte, leaderList [][]byte) *FullBlock {
	// 组装完整区块
	fullBlock := FullBlock{
		BType:    []byte("fullBlock"),
		Round:    round,    // 查数据库得到
		PrevHash: prevHash, // 查数据库得到
		TimeStamp:   miniMaxTimeStamp, // 时间戳
		AccountRoot: accountRoot,      // 计算
		StorageRoot: storageRoot,      // 计算
		TransRoot:   transRoot,        // 计算
		LeaderList:  leaderList,       
	}

	fullBlock.SetHash() // 设置区块hash字段

	return &fullBlock
}

// 生成完整区块
func GenesisFullBlock() *FullBlock {
	// 1 生成账户树
	// a 获取所有节点的公钥列表
	nodeAllPub := GetAllNodeBubKey()

	// logrus.Infof("nodeAllPub数据为 < %v >\n", nodeAllPub)

	// b 生成账户树结构
	accountTree := NewAccountTree(nodeAllPub)

	// logrus.Infof("账户树数据为 < %v >\n", accountTree)

	// c 进行gob序列化
	accountGob := accountTree.Serialize()

	// logrus.Infof("账户树gob之后的数据为 < %v >\n", accountGob)

	// b 进行hash运算
	accountRoot := sha256.Sum256(accountGob[:20])

	logrus.Infof("账户树哈希之后的数据为 < %x >\n", accountRoot)


	// 2 生成创世完整区块
	fullBlock := &FullBlock{
		BType:       []byte("genesisBlock"),
		Round:       1,
		PrevHash:    []byte{},
		TimeStamp:   genesisTimeStamp, // 创世区块时间戳为程序启动时间戳
		AccountRoot: accountRoot[:],
		StorageRoot: []byte{},  // 创世区块无存储树
		TransRoot:   []byte{},  // 创世区块无交易树
		LeaderList:  [][]byte{},  // 创世区块委员会领导者集合
	}

	fullBlock.SetHash() // 设置区块hash字段

	// logrus.Infof("创世区块 < %v >\n", fullBlock)

	// 附录: 将完整区块中的AccountRoot对应的账户树进行存储
	// a 读取数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)

	if err != nil {
		logrus.Panic("打开区块链数据库失败...")
	}

	defer db.Close()

	// b 存储数据
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
		// key: 账户树hash, value: 账户树 gob数据
		bucket.Put(accountRoot[:], accountGob)
		return nil
	})

	return fullBlock
}

// 生成完整区块hash
func (fb *FullBlock) SetHash() {
	// 组装数据
	tmp := [][]byte{
		fb.BType,
		Uint64ToByte(fb.Round),
		fb.PrevHash,
		Uint64ToByte(fb.TimeStamp),
		fb.AccountRoot,
		fb.StorageRoot,
		fb.TransRoot,
	}
	
	// 数据拼接
	data := bytes.Join(tmp, []byte{})
	leaderData := bytes.Join(fb.LeaderList, []byte{})
	result := append(data, leaderData...)

	hash := sha256.Sum256(result)
	fb.Hash = hash[:]
}

// 完整区块序列化
func (fb *FullBlock) Serialize() []byte {
	// 编码的数据放到buffer中
	var buffer bytes.Buffer

	// 使用gob进行序列化(编码)得到字节流
	// 1 定义一个编码器
	// 2 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(fb)
	if err != nil {
		logrus.Panic("完整区块编码出错...")
	}
	return buffer.Bytes()
}

// fullBlock反序列化: 将字节流转为完整区块信息
func FullDeserialize(data []byte) *FullBlock {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var fullBlock FullBlock

	err := decoder.Decode(&fullBlock)
	if err != nil {
		logrus.Panic("完整区块解码错误...")
	}
	return &fullBlock
}
