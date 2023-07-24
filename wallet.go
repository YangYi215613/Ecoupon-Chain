package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)


type Wallet struct {
	Private *ecdsa.PrivateKey
	// 约定: 这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在接收端重新拆分(参考r,s传递)
	PubKey []byte
}

// 创建钱包(直接进行持久化存储)
func NewWallet(nodeID string) *Wallet {
	// 创建曲线
	curve := elliptic.P256()
	// 生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		logrus.Panic("私钥生成失败...")
	}
	// 生成公钥对象
	pubKeyOrig := privateKey.PublicKey
	// 对公钥对象做一个拼接X和Y，生成公钥字节
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)

	// 生成钱包
	wallet := &Wallet{
		privateKey,
		pubKey,
	}

	// 存储钱包
	wallet.SaveToDB(nodeID)
	return wallet
}

// 新钱包持久化存储
func (wt *Wallet) SaveToDB(nodeID string) {
	// 1 钱包数据序列化
	var buffer bytes.Buffer

	// gob注册interface(不然报错)
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wt)
	// 一定要注意校验
	if err != nil {
		logrus.Panic(err)
	}

	// 2 钱包数据持久化存储
	db, err := bolt.Open("wallets.dat", 0600, nil)

	if err != nil {
		logrus.Panic("打开 wallets 钱包失败！")
	}

	defer db.Close()

	// 操作数据库(改写)(事务，参数为匿名函数)
	db.Update(func(tx *bolt.Tx) error {
		// 2 找到抽屉 bucket(如果没有，就创建)
		bucket := tx.Bucket([]byte("wallet"))
		if bucket == nil {
			// 没有抽屉，创建抽屉
			bucket, err = tx.CreateBucket([]byte("wallet"))
			if err != nil {
				logrus.Panic("创建wallet数据库失败！")
			}
		}
		// 3 写数据
		if len(bucket.Get([]byte(nodeID))) == 0 {
			bucket.Put([]byte(nodeID), buffer.Bytes())
		} else {
			logrus.Warnf("该节点 < %s > 已经存在公私钥对, 无需重新创建...\n", nodeID)
		}
		return nil
	})
}

// 根据节点编号获取钱包公私钥
func GetNodeKey(nodeID string) (string, string) {
	// 1 读取钱包数据库
	db, err := bolt.Open("wallets.dat", 0600, nil)

	if err != nil {
		logrus.Panic("打开 wallets 钱包失败！")
	}

	defer db.Close()

	// 2 根据指定编号获取序列化数据
	var data []byte // 存放读出来的数据

	db.View(func(tx *bolt.Tx) error {
		// 2 找到抽屉 bucket(如果没有，就创建)
		bucket := tx.Bucket([]byte("wallet"))
		if bucket == nil {
			fmt.Println("没有wallet数据，请先创建钱包...")
		}
		// 3 读数据
		data = bucket.Get([]byte(nodeID))
		return nil
	})
	// 3 反序列化
	// gob注册interface
	gob.Register(elliptic.P256())
	// 进行gob反序列化
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var wallet Wallet

	err = decoder.Decode(&wallet)
	if err != nil {
		logrus.Panic("wallet钱包解码错误...", err)
	}
	// 4 返回公私钥对字符串
	// 公钥数据
	publicKey := fmt.Sprintf("%x", wallet.PubKey)
	// 私钥数据
	privateKey := fmt.Sprintf("%x", wallet.Private.D.Bytes())
	return publicKey, privateKey
}

// 获取所有节点的公钥信息
func GetAllNodeBubKey() []string {
	// 根据节点编号循环获取公钥信息
	var nodeKeyList []string

	for _, item := range allNodeName {
		pub, _ :=GetNodeKey(item)
		nodeKeyList = append(nodeKeyList, pub)
	}

	return nodeKeyList
}