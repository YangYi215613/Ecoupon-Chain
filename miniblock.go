package main

import (
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/google/uuid"
)

// 微块
type MiniBlock struct {
	BType           []byte     // 区块类型
	Round           uint64     // 系统版本
	PrevHash        []byte     // 前完整区块hash
	TimeStamp       uint64     // 时间戳
	TransTree       *TransTree // 微块中交易的根哈希(Merkle) 微块中不需要这样
	PrevAccountRoot []byte     // 微块本地的前一个账户树根哈希
	PrevStorageRoot []byte     // 微块本地的前一个存储树根哈希
	PrevTransRoot   []byte     // 微块本地的前一个交易树根哈希

	// 附带
	LeaderPuK []byte  // 打包该微块的领导者
}

// 1 生成微块
// 需要传入bucket名，查找自己的数据库
func NewMiniBlock(leaderPuk []byte) *MiniBlock {
	// 生成微块中的交易树
	transTree := &TransTree{
		NorTransList:      transMemPool,
		GeneTransList:     nil,  
		TransferTransList: nil,
	}

	miniBlock := MiniBlock{
		BType: []byte("miniBlock"),  // 区块标识
		Round:           0,              // 查数据库
		PrevHash:        []byte{},       // 查数据库
		TimeStamp:       uint64(time.Now().Unix()),
		TransTree:       transTree, // 计算
		PrevAccountRoot: []byte{},  // 查数据库
		PrevStorageRoot: []byte{},  // 查数据库
		PrevTransRoot:   []byte{},  // 查数据库
		LeaderPuK: leaderPuk,
	}

	// 1 查数据库获取微块数据
	// a 打开数据库(没有则创建)
	db, err := bolt.Open(blockChainDB, 0600, nil)

	if err != nil {
		logrus.Panic("打开数据库失败")
	}

	defer db.Close()

	// b 操作数据库(查看)
	db.View(func(tx *bolt.Tx) error {
		// 拼接bucket名称
		bucketName := blockBucket
		// 找到抽屉 bucket
		bucket := tx.Bucket([]byte(bucketName))
		// 找出数据库中的数据的hash索引
		lastBlockDataIndex := bucket.Get([]byte(lastBlockHash))
		lastBlockData := bucket.Get(lastBlockDataIndex)
		// 数据进行反序列化
		fullBlock := FullDeserialize(lastBlockData)
		// 拿到前一个区块的数据，并进行赋值
		miniBlock.Round = fullBlock.Round + 1 // round+1
		miniBlock.PrevHash = fullBlock.Hash
		miniBlock.PrevAccountRoot = fullBlock.AccountRoot // 进行数据比对使用
		miniBlock.PrevStorageRoot = fullBlock.StorageRoot
		miniBlock.PrevTransRoot = fullBlock.TransRoot

		return nil
	})

	return &miniBlock
}

// 2 发送微块
// 传入当前节点的编号
func (mb *MiniBlock) SendMiniBlock() {
	// 1 添加标识
	data := mb.Serialize()

	uuid := uuid.New()
	prefix := []byte("m" + "$$$" + uuid.String() + "$$$")
	data = append(prefix, data...)

	// 2 发送微块
	for _, item := range FindNodes {
		go sendMessage(item, data)
	}
}

// miniBlock序列化
func (mb *MiniBlock) Serialize() []byte {
	// 编码的数据放到buffer中
	var buffer bytes.Buffer

	// 使用gob进行序列化(编码)得到字节流
	// 1 定义一个编码器
	// 2 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(mb)
	if err != nil {
		logrus.Panic("微块编码出错...")
	}

	return buffer.Bytes()
}

// miniBlock反序列化: 将字节流转为微块信息
func MiniDeserialize(data []byte) *MiniBlock {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var miniBlock MiniBlock

	err := decoder.Decode(&miniBlock)
	if err != nil {
		logrus.Panic("微块解码错误...")
	}
	return &miniBlock
}