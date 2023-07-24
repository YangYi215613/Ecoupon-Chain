package main

import (
	"net"
	"github.com/sirupsen/logrus"
)

// 防止网络风暴
var IsAcceptMsg []string

// 防止网络风暴
func StartServer() {
	// 创建服务端
	listener, err := net.Listen(protocol, networkAddr)
	logrus.Infof("节点 < %s > 已启动, 地址为 < %s >...\n", nodeName, networkAddr)

	if err != nil {
		logrus.Warn("服务器 listen err: ", err)
	}
	defer listener.Close()

	for {
		// 等待客户端连接请求
		conn, err := listener.Accept()
		if err != nil {
			logrus.Warn("accept err:", err)
			return
		}

		// 协程处理用户请求
		go HandleMulServerConn(conn)
	}
}

// 处理发来的请求，实现多个服务端之间通讯
func HandleMulServerConn(conn net.Conn) {
	// 函数调用完毕，自动关闭conn
	defer conn.Close()

	buf := make([]byte, 100000) // 通讯数据量
	n, err := conn.Read(buf) // 读取用户数据

	if err != nil {
		logrus.Warn("read err:", err)
		return
	}

	// 判断收到的数据是否在自己的缓存中，如果在，直接退出; 如果不在，广播给其他节点
	for _, item := range IsAcceptMsg {
		// 如果在，直接退出
		if item == string(buf[:n]) {
			return
		}
	}

	// 设置缓存最大长度
	if len(IsAcceptMsg) > nodeNumber*2 {
		// 清空该节点缓存
		IsAcceptMsg = IsAcceptMsg[:0]
	}

	IsAcceptMsg = append(IsAcceptMsg, string(buf[:n]))

	// 显式接收数据
	if string(buf[0]) == "m" {
		logrus.Infof("该节点收到一个微块")
	} else {
		logrus.Infof("该节点接收到数据 < %s >\n", string(buf[:n]))
	}

	// 执行解析命令对数据进行解析
	if string(buf[0]) == "d" {
		// 命令行信息
		go bc.ParseCMDData(buf[:n]) // 解析命令行数据
		return
	} else {
		// 交易 或 微块信息
		go bc.ParseData(buf[:n]) // 解析交易、微块数据
	}

	// FindNodes 为节点发现的节点列表
	for _, item := range FindNodes {
		go sendMessage(item, buf[:n]) // 并发发送消息
	}
}

// 发送数据
func sendMessage(item string, buf []byte) {
	clientConn, err := net.Dial(protocol, item)

	if err != nil {
		logrus.Println("clientConn err:", item)
		return
	}

	defer clientConn.Close()

	_, err = clientConn.Write(buf)

	if err != nil {
		logrus.Println("Write data err:", err)
		return
	}
}
