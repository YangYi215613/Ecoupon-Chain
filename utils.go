package main

import (
	"bytes"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func StartConfigs() {
	// 1 获取系统环境变量
	GetEnv()

	// 2 日志配置
	InitLogger()

	// 3 节点发现
	NodeFind()

	// 4 钱包配置
	// 判断钱包文件是否存在
	_, err := os.Stat("./wallets.dat")

	if err == nil {
		return
	}

	for _, item := range allNodeName {
		NewWallet(item)
	}
}

func GetEnv() {
	networkAddr = os.Getenv("NETWORK_ADDR") + ":8989"
	nodeNumber, _ = strconv.Atoi(os.Getenv("NODE_NUMBER"))
	nodeName = os.Getenv("NODENAME")
	allNodeName = strings.Split(os.Getenv("ALL_NODE_NAME"), ";")
	committeeNum, _ = strconv.Atoi(os.Getenv("COMMITTEE_NUM"))
	isCommitteeLeader, _ = strconv.ParseBool(os.Getenv("IS_COMMITTEELEADER"))
	nodeCommitteeID, _ = strconv.Atoi(os.Getenv("NODE_COMMITTEEID"))
}


// 节点发现
func NodeFind() {
	var hasNumber []int
	var flag bool
	
	rand.Seed(time.Now().UnixNano())
	for i:=0; i<maxConn; i++ {
		var randomNumber int
		for {
			flag = true
			randomNumber = rand.Intn(20) + 1
			for _, item := range hasNumber{
				if item == randomNumber {
					flag = false
					break
				}
			}
			if flag { break }
		}
		FindNodes[i] = "container" + strconv.Itoa(randomNumber) + ":8989"
	}
}


func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num) // BigEndian大端对齐
	if err != nil {
		logrus.Panic(err)
	}
	return buffer.Bytes()
}