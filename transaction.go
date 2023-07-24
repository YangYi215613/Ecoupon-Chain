package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/sirupsen/logrus"
	"math/big"
	"strconv"
)

// 普通交易(发送数据的时候，使用私钥进行签名)
type NorTransaction struct {
	TxType []byte // 交易类型
	TxID   []byte // 交易ID  // 交易的hash值
	From   []byte // 发起人(公钥)
	To     []byte // 收款人(公钥)
	Value uint64 // 转账金额
	//TimeStamp uint64 // 时间戳
}

// 转让交易
// a 请求信息
type Request struct {
	RType     []byte // 交易类型
	From      []byte // 发起人(公钥)
	To        []byte // 收款人(公钥)
	Msg       []byte // 请求信息
	TimeStamp uint64 // 时间戳
}

// b 响应信息
type Response struct {
	RType     []byte   // 交易类型
	To        []byte   // 收款人(公钥)
	From      []byte   // 发起人(公钥)
	Req       *Request // 请求信息
	Con       *Coupon  // 转移的卡券证书
	Msg       []byte   // 收款人附带信息
	TimeStamp uint64   // 时间戳
	Value     uint64   // 交易额度
}

// c 真正交易
type TransferTransaction struct {
	TxType    []byte    // 交易类型
	TxID      []byte    // 交易ID
	Count     uint64    // 交易次数
	Req       *Request  // 请求信息
	Res       *Response // 响应信息
	TimeStamp uint64    // 时间戳
}

// 分发交易
type GeneTransaction struct {
	TxType    []byte // 交易类型
	TxID      []byte // 交易ID
	Count     uint64 // 交易次数
	Fee       uint64 // 交易费用
	TimeStamp uint64 // 时间戳
}

func NewNorTransaction(from, to []byte, value uint64) *NorTransaction {
	norTrans := NorTransaction{
		TxType: []byte("norTrans"),
		TxID:   []byte{},
		From:   from,
		To:     to,
		Value:  value,
	}

	tmp := [][]byte{
		norTrans.TxType,
		norTrans.From,
		norTrans.To,
		Uint64ToByte(norTrans.Value),
	}

	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	norTrans.TxID = hash[:]

	return &norTrans
}

/*// 发送正常交易
func (nt *NorTransaction) SendNorTrans() {
	data := nt.Serialize()
	data = append([]byte("a"), data...)

}*/


// 序列化正常交易
func (nor *NorTransaction) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(nor)
	if err != nil {
		logrus.Panic("正常交易编码出错...")
	}

	return buffer.Bytes()
}

// 反序列化正常交易
func NorDeserialize(data []byte) *NorTransaction {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var norTrans NorTransaction

	err := decoder.Decode(&norTrans)
	if err != nil {
		logrus.Panic("微块解码错误...")
	}
	return &norTrans
}

// 获取打包该交易的委员会ID【同样的交易返回固定的ID】
func (norTx *NorTransaction) GetCommitteeID() uint64 {
	// 1. 交易取哈希
	// 2. 哈希生成特定符号
	norSerInfo := norTx.Serialize()
	hash := sha256.Sum256(norSerInfo)

	// 获取哈希值的big.Int
	hashValue := big.Int{}
	hashValue.SetBytes(hash[:])

	// 将委员会数目编码成big.Int
	committeeValue := big.Int{}
	committeeValue.SetInt64(int64(committeeNum))

	// 自动实现取模运算，得到big.Int
	res := big.Int{}
	res.Mod(&hashValue, &committeeValue)

	valueStr := res.String()
	txID, _ := strconv.Atoi(valueStr)

	return uint64(txID) + 1
}

// 序列化分发交易
func (gen *GeneTransaction) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(gen)
	if err != nil {
		logrus.Panic("分发交易编码出错...")
	}

	return buffer.Bytes()
}

// 反序列化分发交易
func GeneDeserialize(data []byte) *GeneTransaction {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var geneTrans GeneTransaction
	err := decoder.Decode(&geneTrans)
	if err != nil {
		logrus.Panic("分发交易解码错误...")
	}
	return &geneTrans
}

// 序列化转让交易
func (transfer *TransferTransaction) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(transfer)
	if err != nil {
		logrus.Panic("转让交易编码出错...")
	}

	return buffer.Bytes()
}

// 反序列化转让交易
func TransferDeserialize(data []byte) *TransferTransaction {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var transferTrans TransferTransaction
	err := decoder.Decode(&transferTrans)
	if err != nil {
		logrus.Panic("转让交易解码错误...")
	}
	return &transferTrans
}
