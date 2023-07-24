package main

import (
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type BlockChainIterator struct {
	db *bolt.DB
	currentHashPointer []byte  	// 游标
}

// 创建迭代器
func(bc *BlockChain) NewIterator() *BlockChainIterator{
	return &BlockChainIterator{
		db: bc.Db,
		// 最初指向区块链的最后一个区块，随着Next的调用，不断变化
		currentHashPointer: []byte(bc.Tail),
	}
}


// Next 
// 1 返回当前的区块
// 2 指针前移
func(it *BlockChainIterator) Next() *FullBlock {
	var fullBlock *FullBlock

	// 打开区块链数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)

	if err != nil {
		logrus.Panic("打开区块链数据库失败...")
	}

	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			logrus.Panic("迭代器遍历时bucket不应该为空，请检查!")
		}

		blockTmp := bucket.Get(it.currentHashPointer)

		fullBlock = FullDeserialize(blockTmp)
		// 游标hash左移
		it.currentHashPointer = fullBlock.PrevHash
		return nil
	})
	return fullBlock
}

