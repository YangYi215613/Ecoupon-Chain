package main

import (
	"bytes"
	"encoding/gob"
	"github.com/sirupsen/logrus"
)

/* 账户树、存储树、交易树 */

// 账户树
type AccountTree struct {
	// 映射关系
	UserInfoMap map[string]*UserInfo // key: 公钥 value: 用户数据
}

// 存储树
type StorageTree struct {
	// 映射关系
	CouponList []*Coupon // 卡券列表
}

// 交易树
type TransTree struct {
	NorTransList      []*NorTransaction      // 正常交易列表
	GeneTransList     []*GeneTransaction     // 分发交易列表
	TransferTransList []*TransferTransaction // 转让交易列表
}

// 生成账户树，只需要在系统一开始的时候进行初始化
// pubKeyList 传入系统中用户的公钥
func NewAccountTree(pubKeyList []string) *AccountTree {
	// a 生成每个用户数据初始化value结构

	var userInfoMap = make(map[string]*UserInfo)

	accountTree := AccountTree{
		UserInfoMap: userInfoMap,
	}

	// 进行数据添加
	for _, pubKey := range pubKeyList {
		userInfo := NewUserInfo()                // 生成每个用户数据初始化value结构
		accountTree.UserInfoMap[pubKey] = userInfo // 是不是指针的问题(进行hash的话，需要传值)
	}
	return &accountTree
}

// 交易树生成
func NewTransTree() *TransTree {
	return &TransTree{
		NorTransList:      []*NorTransaction{},
		GeneTransList:     nil,
		TransferTransList: nil,
	}
}

// 账户树序列化
func (at *AccountTree) Serialize() []byte {
	// 编码的数据放到buffer中
	var buffer bytes.Buffer

	// 使用gob进行序列化(编码)得到字节流
	// 1 定义一个编码器
	// 2 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(at)
	if err != nil {
		logrus.Panic("账户树编码出错...")
	}

	return buffer.Bytes()
}

// 账户树反序列化
func AccountDeserialize(data []byte) *AccountTree {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var accountTree AccountTree

	err := decoder.Decode(&accountTree)
	if err != nil {
		logrus.Panic("账户树解码错误...")
	}
	return &accountTree
}

func (st *StorageTree) Serialize() []byte {
	// 编码的数据放到buffer中
	var buffer bytes.Buffer

	// 使用gob进行序列化(编码)得到字节流
	// 1 定义一个编码器
	// 2 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(st)
	if err != nil {
		logrus.Panic("存储树编码出错...")
	}

	return buffer.Bytes()
}

func (tt *TransTree) Serialize() []byte {
	// 编码的数据放到buffer中
	var buffer bytes.Buffer

	// 使用gob进行序列化(编码)得到字节流
	// 1 定义一个编码器
	// 2 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tt)
	if err != nil {
		logrus.Panic("交易树编码出错...")
	}

	return buffer.Bytes()
}
