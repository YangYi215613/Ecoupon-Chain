package main

import (
	"fmt"
)

func main() {
	fmt.Printf("\n================= Welcome To Ecoupon-Chain =================\n\n")

	// 1 初始参数配置
	StartConfigs()

	// 2 初始化区块链
	InitialBlockChain()

	// 3 获取区块链对象
	NewBlockChain()

	// 4 运行网络服务
	StartServer()

	// 阻塞
	var ch chan struct{}
	<-ch
}